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
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/pesos/grofer/src/utils"
	proc "github.com/shirou/gopsutil/process"
)

var runAllProc = true

func getData(procs []*proc.Process) []string {
	var data []string
	for _, info := range procs {
		exe, err := info.Exe()
		if err == nil {
			temp := strconv.Itoa(int(info.Pid))

			for i := 0; i < 12-len(strconv.Itoa(int(info.Pid))); i++ {
				temp = temp + " "
			}

			commands := strings.Split(exe, "/")
			command := commands[len(commands)-1]

			if len(command) > 40 {
				command = command[:40]
			} else {
				temp = temp + "[" + command + "](fg:green)" + strings.Repeat(" ", 41-len(command))
			}

			tempCPU, err := info.CPUPercent()
			cpuPercent := ""
			if err == nil {
				cpuPercent = fmt.Sprintf("%.2f%s", tempCPU, "%")
				temp = temp + cpuPercent
			}
			temp = temp + strings.Repeat(" ", 11-len(cpuPercent))

			tempMem, err := info.MemoryPercent()
			memPercent := ""
			if err == nil {
				memPercent = fmt.Sprintf("%.2f%s", tempMem, "%")
				temp = temp + memPercent
			}
			temp = temp + strings.Repeat(" ", 11-len(memPercent))

			status, err := info.Status()
			if err == nil {
				temp = temp + status
			}
			temp = temp + strings.Repeat(" ", 9-len(status))

			fg, err := info.Foreground()
			if err == nil {
				if fg {
					temp = temp + "True"
					temp = temp + strings.Repeat(" ", 9)
				} else {
					temp = temp + "False"
					temp = temp + strings.Repeat(" ", 8)
				}
			}

			ctime, err := info.CreateTime()
			createTime := ""
			if err == nil {
				createTime := utils.GetDateFromUnix(ctime)
				temp = temp + createTime
			}
			temp = temp + strings.Repeat(" ", 9-len(createTime))

			threads, err := info.NumThreads()
			if err == nil {
				threadCount := strconv.FormatInt(int64(threads), 10)
				temp = temp + threadCount
			}

			data = append(data, temp)
		}
	}
	return data
}

func AllProcVisuals(dataChannel chan []*proc.Process,
	endChannel chan os.Signal,
	refreshRate uint64,
	wg *sync.WaitGroup) {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	var on sync.Once

	myPage := NewAllProcsPage()

	updateUI := func() {
		w, h := ui.TerminalDimensions()
		ui.Clear()
		myPage.Grid.SetRect(0, 0, w, h)
		ui.Render(myPage.Grid)
	}

	updateUI() // Render empty UI

	pause := func() {
		runAllProc = !runAllProc
	}

	uiEvents := ui.PollEvents()
	tick := time.Tick(time.Duration(refreshRate) * time.Millisecond)

	previousKey := ""
	selectedStyle := ui.NewStyle(ui.ColorYellow, ui.ColorClear, ui.ModifierBold)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": //q or Ctrl-C to quit
				endChannel <- os.Kill
				ui.Close()
				wg.Done()
				return
			case "<Resize>":
				updateUI() // updateUI only during resize event
			case "s": //s to pause
				pause()
			case "j", "<Down>":
				myPage.BodyList.ScrollDown()
			case "k", "<Up>":
				myPage.BodyList.ScrollUp()
			case "<C-d>":
				myPage.BodyList.ScrollHalfPageDown()
			case "<C-u>":
				myPage.BodyList.ScrollHalfPageUp()
			case "<C-f>":
				myPage.BodyList.ScrollPageDown()
			case "<C-b>":
				myPage.BodyList.ScrollPageUp()
			case "g":
				if previousKey == "g" {
					myPage.BodyList.ScrollTop()
				}
			case "<Home>":
				myPage.BodyList.ScrollTop()
			case "G", "<End>":
				myPage.BodyList.ScrollBottom()
			}

			ui.Render(myPage.Grid)
			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			myPage.BodyList.SelectedRowStyle = selectedStyle
			if runAllProc {
				myPage.BodyList.Rows = getData(data)

				on.Do(func() {
					w, h := ui.TerminalDimensions()
					ui.Clear()
					myPage.Grid.SetRect(0, 0, w, h)
					ui.Render(myPage.Grid)
				})
			}
		case <-tick: // Update page with new values
			ui.Render(myPage.Grid)
		}
	}
}
