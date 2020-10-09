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
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	h "github.com/pesos/grofer/src/display/misc"
	info "github.com/pesos/grofer/src/general"
	"github.com/pesos/grofer/src/utils"
	proc "github.com/shirou/gopsutil/process"
)

var runAllProc = true
var helpVisible = false

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
	ctx context.Context,
	refreshRate uint64) error {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	var on sync.Once
	var help *h.HelpMenu = h.NewHelpMenu()
	h.SelectHelpMenu("proc")

	myPage := NewAllProcsPage()

	updateUI := func() {
		w, h := ui.TerminalDimensions()
		myPage.Grid.SetRect(0, 0, w, h)
		help.Resize(w, h)
		ui.Clear()
		if helpVisible {
			ui.Render(help)
		} else {
			ui.Render(myPage.Grid)
		}
	}

	updateUI() // Render empty UI

	pauseProc := func() {
		runAllProc = !runAllProc
	}

	uiEvents := ui.PollEvents()
	tick := time.Tick(time.Duration(refreshRate) * time.Millisecond)

	previousKey := ""
	selectedStyle := ui.NewStyle(ui.ColorYellow, ui.ColorClear, ui.ModifierBold)
	killingStyle := ui.NewStyle(ui.ColorWhite, ui.ColorMagenta, ui.ModifierBold)
	errorStyle := ui.NewStyle(ui.ColorBlack, ui.ColorRed, ui.ModifierBold)

	// updates process list immediately
	updateProcs := func() {
		if runAllProc {
			procs, err := proc.Processes()
			if err == nil {
				myPage.BodyList.Rows = getData(procs)
			}
		}
	}

	// whether a process is selected for killing (UI controls are paused)
	killSelected := false
	var pidToKill int32

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": //q or Ctrl-C to quit
				return info.ErrCanceledByUser
			case "?":
				helpVisible = !helpVisible
			case "<Resize>":
				updateUI() // updateUI only during resize event
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
				case "<Escape>":
					if killSelected {
						runAllProc = true
						killSelected = false
						myPage.BodyList.SelectedRowStyle = selectedStyle
					}
				case "s": //s to pause
					if !killSelected {
						pauseProc()
					}
				case "j", "<Down>":
					if !killSelected {
						myPage.BodyList.ScrollDown()
					}
				case "k", "<Up>":
					if !killSelected {
						myPage.BodyList.ScrollUp()
					}
				case "<C-d>":
					if !killSelected {
						myPage.BodyList.ScrollHalfPageDown()
					}
				case "<C-u>":
					if !killSelected {
						myPage.BodyList.ScrollHalfPageUp()
					}
				case "<C-f>":
					if !killSelected {
						myPage.BodyList.ScrollPageDown()
					}
				case "<C-b>":
					if !killSelected {
						myPage.BodyList.ScrollPageUp()
					}
				case "g":
					if !killSelected && previousKey == "g" {
						myPage.BodyList.ScrollTop()
					}
				case "<Home>":
					if !killSelected {
						myPage.BodyList.ScrollTop()
					}
				case "G", "<End>":
					if !killSelected {
						myPage.BodyList.ScrollBottom()
					}
				case "K", "<F9>":
					if myPage.BodyList.SelectedRow < len(myPage.BodyList.Rows) {
						row := myPage.BodyList.Rows[myPage.BodyList.SelectedRow]
						// get PID from the data
						pid64, err := strconv.ParseInt(strings.SplitN(row, " ", 2)[0], 10, 32)
						if err != nil {
							return fmt.Errorf("Failed to get PID of process: %v", err)
						}
						pidToKill = int32(pid64)

						if !killSelected {
							runAllProc = false
							killSelected = true
							myPage.BodyList.SelectedRowStyle = killingStyle
						} else {
							// get process and kill it
							procToKill, err := proc.NewProcess(pidToKill)
							myPage.BodyList.SelectedRowStyle = selectedStyle
							if err == nil {
								err = procToKill.Kill()
								if err != nil {
									myPage.BodyList.SelectedRowStyle = errorStyle
								}
							} else {
								myPage.BodyList.SelectedRowStyle = errorStyle
							}
							runAllProc = true
							killSelected = false
							updateProcs()
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
			if runAllProc {
				myPage.BodyList.SelectedRowStyle = selectedStyle
				myPage.BodyList.Rows = getData(data)

				on.Do(updateUI)
			}

		case <-tick: // Update page with new values
			if killSelected {
				exists, _ := proc.PidExists(pidToKill)
				if !exists {
					runAllProc = true
					killSelected = false
					myPage.BodyList.SelectedRowStyle = selectedStyle
					updateProcs()
				}
			} else {
				myPage.BodyList.SelectedRowStyle = selectedStyle
			}
			if !helpVisible {
				ui.Render(myPage.Grid)
			}
		}
	}
}
