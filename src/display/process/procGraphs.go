/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package process

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/pesos/grofer/src/process"
	"github.com/pesos/grofer/src/utils"
)

var runProc = true

func getChildProcs(proc *process.Process) []string {
	childProcs := []string{"PID                   Command"}
	for _, proc := range proc.Children {
		exe, err := proc.Exe()
		if err == nil {
			temp := strconv.Itoa(int(proc.Pid))
			for i := 0; i < 22-len(strconv.Itoa(int(proc.Pid))); i++ {
				temp = temp + " "
			}
			temp = temp + "[" + exe + "](fg:green)"
			childProcs = append(childProcs, temp)
		} else {
			childProcs = append(childProcs, "["+strconv.Itoa(int(proc.Pid))+"](fg:yellow)"+"            "+"NA")
		}
	}

	return childProcs
}

// ProcVisuals renders graphs and charts for per-process stats.
func ProcVisuals(endChannel chan os.Signal,
	dataChannel chan *process.Process,
	wg *sync.WaitGroup) {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	// Create new page
	myPage := NewPerProcPage()

	var statusMap map[string]string = map[string]string{
		"R": "Running",
		"S": "Sleep",
		"Z": "Zombie",
		"T": "Stop",
		"I": "Idle",
		"W": "Wait",
		"L": "Lock",
	}

	pause := func() {
		runProc = !runProc
	}

	uiEvents := ui.PollEvents()
	tick := time.Tick(1 * time.Second)

	previousKey := ""
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": //q or Ctrl-C to quit
				endChannel <- os.Kill
				ui.Close()
				wg.Done()
				return
			case "s": //s to pause
				pause()
			case "j", "<Down>":
				myPage.ChildProcsList.ScrollDown()
			case "k", "<Up>":
				myPage.ChildProcsList.ScrollUp()
			case "<C-d>":
				myPage.ChildProcsList.ScrollHalfPageDown()
			case "<C-u>":
				myPage.ChildProcsList.ScrollHalfPageUp()
			case "<C-f>":
				myPage.ChildProcsList.ScrollPageDown()
			case "<C-b>":
				myPage.ChildProcsList.ScrollPageUp()
			case "g":
				if previousKey == "g" {
					myPage.ChildProcsList.ScrollTop()
				}
			case "<Home>":
				myPage.ChildProcsList.ScrollTop()
			case "G", "<End>":
				myPage.ChildProcsList.ScrollBottom()
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			if runProc {
				// update ctx switches
				switches, units := utils.RoundValues(float64(data.NumCtxSwitches.Voluntary), float64(data.NumCtxSwitches.Involuntary))

				myPage.CTXSwitchesChart.Data = switches
				myPage.CTXSwitchesChart.Title = " CTX Switches" + units

				// update cpu %
				myPage.CPUChart.Percent = int(data.CPUPercent)

				// update mem %
				myPage.MemChart.Percent = int(data.MemoryPercent)

				// update proc info
				myPage.PIDTable.Rows = [][]string{
					[]string{"[Name](fg:yellow)", data.Name},
					[]string{"[Command](fg:yellow)", data.Exe},
					[]string{"[Status](fg:yellow)", statusMap[data.Status] + " (" + data.Status + ")"},
					[]string{"[Background](fg:yellow)", strconv.FormatBool(data.Background)},
					[]string{"[Foreground](fg:yellow)", strconv.FormatBool(data.Foreground)},
					[]string{"[Running](fg:yellow)", strconv.FormatBool(data.IsRunning)},
					[]string{"[Creation Time](fg:yellow)", utils.GetDateFromUnix(data.CreateTime)},
					[]string{"[Nice value](fg:yellow)", strconv.Itoa(int(data.Nice))},
					[]string{"[Thread count](fg:yellow)", strconv.Itoa(int(data.NumThreads))},
					[]string{"[Child process count](fg:yellow)", strconv.Itoa(len(data.Children))},
				}
				myPage.PIDTable.Title = " PID: " + strconv.Itoa(int(data.Proc.Pid)) + " "

				//update memory stats
				memData := []float64{utils.GetInMB(data.MemoryInfo.RSS, 1),
					utils.GetInMB(data.MemoryInfo.Data, 1),
					utils.GetInMB(data.MemoryInfo.Stack, 1),
					utils.GetInMB(data.MemoryInfo.Swap, 1),
				}
				myPage.MemStatsChart.Data = memData

				//update page faults
				faults, units := utils.RoundValues(float64(data.PageFault.MinorFaults), float64(data.PageFault.MajorFaults))

				myPage.PageFaultsChart.Data = faults
				myPage.PageFaultsChart.Title = " Page Faults" + units
				myPage.ChildProcsList.Rows = getChildProcs(data)
			}

		case <-tick:
			w, h := ui.TerminalDimensions()
			ui.Clear()
			myPage.Grid.SetRect(0, 0, w, h)
			ui.Render(myPage.Grid)
		}
	}
}
