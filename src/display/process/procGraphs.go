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
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	h "github.com/pesos/grofer/src/display/misc"
	info "github.com/pesos/grofer/src/general"
	"github.com/pesos/grofer/src/process"
	"github.com/pesos/grofer/src/utils"
)

var runProc = true

func getChildProcs(proc *process.Process) [][]string {
	childProcs := [][]string{}
	for _, proc := range proc.Children {
		pid := strconv.Itoa(int(proc.Pid))
		exe, err := proc.Exe()
		cmd := "NA"
		if err == nil {
			cmd = exe
		}
		childProcs = append(childProcs, []string{pid, cmd})
	}
	return childProcs
}

// ProcVisuals renders graphs and charts for per-process stats.
func ProcVisuals(ctx context.Context,
	dataChannel chan *process.Process,
	refreshRate uint64) error {

	defer ui.Close()
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	var on sync.Once
	var help *h.HelpMenu = h.NewHelpMenu()
	h.SelectHelpMenu("proc_pid")

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

	updateUI := func() {

		// Get Terminal Dimensions adn clear the UI
		w, h := ui.TerminalDimensions()

		// Adjust Memory Stats Bar graph values
		myPage.MemStatsChart.BarGap = ((w / 2) - (4 * myPage.MemStatsChart.BarWidth)) / 4

		// Adjust Page Faults Bar graph values
		myPage.PageFaultsChart.BarGap = ((w / 4) - (2 * myPage.PageFaultsChart.BarWidth)) / 2

		// Adjust Context Switches Bar graph values
		myPage.CTXSwitchesChart.BarGap = ((w / 4) - (2 * myPage.CTXSwitchesChart.BarWidth)) / 2

		// Adjust Grid dimensions
		myPage.Grid.SetRect(0, 0, w, h)
		help.Resize(w, h)
		if helpVisible {
			ui.Clear()
			ui.Render(help)
		} else {
			ui.Render(myPage.Grid)
		}
	}

	uiEvents := ui.PollEvents()
	t := time.NewTicker(time.Duration(refreshRate) * time.Millisecond)
	tick := t.C

	previousKey := ""

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				//q or Ctrl-C to quit
				return info.ErrCanceledByUser
			case "<Resize>":
				updateUI()
			case "?":
				helpVisible = !helpVisible
			}
			if helpVisible {
				switch e.ID {
				case "?":
					updateUI()
				case "<Escape>":
					helpVisible = false
					updateUI()
				case "j", "<Down>":
					help.List.ScrollDown()
					ui.Render(help)
				case "k", "<Up>":
					help.List.ScrollUp()
					ui.Render(help)
				}
			} else {
				switch e.ID {
				case "?":
					updateUI()
				case "s": //s to pause
					pause()
				case "j", "<Down>":
					myPage.ChildProcsTable.ScrollDown()
				case "k", "<Up>":
					myPage.ChildProcsTable.ScrollUp()
				case "<C-d>":
					myPage.ChildProcsTable.ScrollHalfPageDown()
				case "<C-u>":
					myPage.ChildProcsTable.ScrollHalfPageUp()
				case "<C-f>":
					myPage.ChildProcsTable.ScrollPageDown()
				case "<C-b>":
					myPage.ChildProcsTable.ScrollPageUp()
				case "g":
					if previousKey == "g" {
						myPage.ChildProcsTable.ScrollTop()
					}
				case "<Home>":
					myPage.ChildProcsTable.ScrollTop()
				case "G", "<End>":
					myPage.ChildProcsTable.ScrollBottom()
				case "<Enter>":
					if myPage.ChildProcsTable.SelectedRow != 0 {
						row := myPage.ChildProcsTable.Rows[myPage.ChildProcsTable.SelectedRow]
						// get PID from the data
						pid, err := strconv.ParseInt(strings.SplitN(row[0], " ", 2)[0], 10, 32)
						if err != nil {
							return fmt.Errorf("Failed to get PID of process: %v", err)
						}
						eg, ctx := errgroup.WithContext(context.Background())
						proc, _ := process.NewProcess(int32(pid))
						dataChannel := make(chan *process.Process, 1)
						eg.Go(func() error {
							return process.Serve(proc, dataChannel, ctx, int64(4*refreshRate/5))
						})
						eg.Go(func() error {
							return ProcVisuals(ctx, dataChannel, refreshRate)
						})
						if err := eg.Wait(); err != nil {
							if err != info.ErrCanceledByUser {
								fmt.Printf("Error: %v\n", err)
							}
						}

					}
				}

				ui.Render(myPage.Grid)
				if previousKey == "g" {
					previousKey = ""
				} else {
					previousKey = e.ID
				}
			}

		case data := <-dataChannel:
			if runProc {
				// update ctx switches
				switches, units := utils.RoundValues(float64(data.NumCtxSwitches.Voluntary), float64(data.NumCtxSwitches.Involuntary), false)

				myPage.CTXSwitchesChart.Data = switches
				myPage.CTXSwitchesChart.Title = " CTX Switches" + units

				// update cpu %
				myPage.CPUChart.Percent = int(data.CPUPercent)

				// update mem %
				myPage.MemChart.Percent = int(data.MemoryPercent)

				// update proc info
				myPage.PIDTable.Rows = [][]string{
					{"[Name](fg:yellow)", data.Name},
					{"[Command](fg:yellow)", data.Exe},
					{"[Status](fg:yellow)", statusMap[data.Status] + " (" + data.Status + ")"},
					{"[Background](fg:yellow)", strconv.FormatBool(data.Background)},
					{"[Foreground](fg:yellow)", strconv.FormatBool(data.Foreground)},
					{"[Running](fg:yellow)", strconv.FormatBool(data.IsRunning)},
					{"[Creation Time](fg:yellow)", utils.GetDateFromUnix(data.CreateTime)},
					{"[Nice value](fg:yellow)", strconv.Itoa(int(data.Nice))},
					{"[Thread count](fg:yellow)", strconv.Itoa(int(data.NumThreads))},
					{"[Child process count](fg:yellow)", strconv.Itoa(len(data.Children))},
					{"[Last Update](fg:yellow)", time.Now().Format("15:04:05")},
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
				faults, units := utils.RoundValues(float64(data.PageFault.MinorFaults), float64(data.PageFault.MajorFaults), false)

				myPage.PageFaultsChart.Data = faults
				myPage.PageFaultsChart.Title = " Page Faults" + units
				myPage.ChildProcsTable.Rows = getChildProcs(data)

				on.Do(updateUI)
			}

		case <-tick:
			if !helpVisible {
				ui.Render(myPage.Grid)
			}
		}
	}
}
