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
	"sort"
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
var sortIdx = -1
var sortAsc = false
var header = []string{
	"PID",
	"Command",
	"CPU",
	"Memory",
	"Status",
	"Foreground",
	"Creation Time",
	"Thread Count",
}

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

func getData(procs []*proc.Process) [][]string {
	procData := [][]string{}
	for _, p := range procs {
		// Get command
		cmd := ""
		exe, err := p.Exe()
		if err == nil {
			cmds := strings.Split(exe, "/")
			cmd = cmds[len(cmds)-1]

			// Get CPU
			cpu := ""
			cpuPercent, err := p.CPUPercent()
			if err == nil {
				cpu = fmt.Sprintf("%.2f%%", cpuPercent)
			}

			// Get Mem
			mem := ""
			memPercent, err := p.MemoryPercent()
			if err == nil {
				mem = fmt.Sprintf("%.2f%%", memPercent)
			}

			// Get Status
			status, _ := p.Status()

			// Get Foreground
			fg, _ := p.Foreground()

			// Get Creation time
			t, err := p.CreateTime()
			ctime := ""
			if err == nil {
				ctime = utils.GetDateFromUnix(t)
			}

			// Get Thread Count
			tc, _ := p.NumThreads()

			// Aggregate row
			r := []string{
				fmt.Sprintf("%d", p.Pid),
				cmd,
				cpu,
				mem,
				status,
				fmt.Sprintf("%t", fg),
				ctime,
				fmt.Sprintf("%d", tc),
			}
			procData = append(procData, r)
		}
	}

	return procData
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
		if helpVisible {
			ui.Clear()
			ui.Render(help)
		} else {
			ui.Render(myPage.Grid)
		}
	}

	updateUI() // Render empty UI

	pauseProc := func() {
		runAllProc = !runAllProc
	}

	intSort := func(i, j int) bool {
		x, _ := strconv.Atoi(myPage.ProcTable.Rows[i][sortIdx])
		y, _ := strconv.Atoi(myPage.ProcTable.Rows[j][sortIdx])
		if sortAsc {
			return x < y
		}
		return x > y
	}

	strSort := func(i, j int) bool {
		if sortAsc {
			return myPage.ProcTable.Rows[i][sortIdx] < myPage.ProcTable.Rows[j][sortIdx]
		}
		return myPage.ProcTable.Rows[i][sortIdx] > myPage.ProcTable.Rows[j][sortIdx]
	}

	floatSort := func(i, j int) bool {
		x1 := myPage.ProcTable.Rows[i][sortIdx]
		y1 := myPage.ProcTable.Rows[j][sortIdx]
		x, _ := strconv.ParseFloat(x1[:len(x1)-1], 32)
		y, _ := strconv.ParseFloat(y1[:len(y1)-1], 32)
		if sortAsc {
			return x < y
		}
		return x > y
	}

	sortFuncs := map[int]func(i, j int) bool{
		0: intSort,
		1: strSort,
		2: floatSort,
		3: floatSort,
		4: strSort,
		5: strSort,
		6: strSort,
		7: intSort,
	}

	uiEvents := ui.PollEvents()
	t := time.NewTicker(time.Duration(refreshRate) * time.Millisecond)
	tick := t.C

	previousKey := ""
	selectedStyle := ui.ColorCyan
	killingStyle := ui.ColorMagenta
	errorStyle := ui.ColorRed

	// updates process list immediately
	updateProcs := func() {
		if runAllProc {
			procs, err := proc.Processes()
			if err == nil {
				myPage.ProcTable.Rows = getData(procs)
			}
		}
	}

	// whether a process is selected for killing (UI controls are paused)
	killSelected := false
	var pidToKill int32

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

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
						myPage.ProcTable.CursorColor = selectedStyle
					}
				case "s": //s to pause
					if !killSelected {
						pauseProc()
					}
				case "j", "<Down>":
					if !killSelected {
						myPage.ProcTable.ScrollDown()
					}
				case "k", "<Up>":
					if !killSelected {
						myPage.ProcTable.ScrollUp()
					}
				case "<C-d>":
					if !killSelected {
						myPage.ProcTable.ScrollHalfPageDown()
					}
				case "<C-u>":
					if !killSelected {
						myPage.ProcTable.ScrollHalfPageUp()
					}
				case "<C-f>":
					if !killSelected {
						myPage.ProcTable.ScrollPageDown()
					}
				case "<C-b>":
					if !killSelected {
						myPage.ProcTable.ScrollPageUp()
					}
				case "g":
					if !killSelected && previousKey == "g" {
						myPage.ProcTable.ScrollTop()
					}
				case "<Home>":
					if !killSelected {
						myPage.ProcTable.ScrollTop()
					}
				case "G", "<End>":
					if !killSelected {
						myPage.ProcTable.ScrollBottom()
					}
				case "K", "<F9>":
					if myPage.ProcTable.SelectedRow < len(myPage.ProcTable.Rows) {
						row := myPage.ProcTable.Rows[myPage.ProcTable.SelectedRow]
						// get PID from the data
						pid, err := strconv.Atoi(row[0])
						if err != nil {
							return fmt.Errorf("failed to get PID of process: %v", err)
						}
						pidToKill = int32(pid)

						if !killSelected {
							runAllProc = false
							killSelected = true
							myPage.ProcTable.CursorColor = killingStyle
						} else {
							// get process and kill it
							procToKill, err := proc.NewProcess(pidToKill)
							myPage.ProcTable.CursorColor = selectedStyle
							if err == nil {
								err = procToKill.Kill()
								if err != nil {
									myPage.ProcTable.CursorColor = errorStyle
								}
							} else {
								myPage.ProcTable.CursorColor = errorStyle
							}
							runAllProc = true
							killSelected = false
							updateProcs()
						}
					}
					// Sort Ascending
				case "1", "2", "3", "4", "5", "6", "7", "8":
					myPage.ProcTable.Header = append([]string{}, header...)
					idx, _ := strconv.Atoi(e.ID)
					sortIdx = idx - 1
					myPage.ProcTable.Header[sortIdx] = header[sortIdx] + " " + UP_ARROW
					sortAsc = true
					sort.Slice(myPage.ProcTable.Rows, sortFuncs[sortIdx])

				// Sort Descending
				case "<F1>", "<F2>", "<F3>", "<F4>", "<F5>", "<F6>", "<F7>", "<F8>":
					myPage.ProcTable.Header = append([]string{}, header...)
					idx, _ := strconv.Atoi(e.ID[2:3])
					sortIdx = idx - 1
					myPage.ProcTable.Header[sortIdx] = header[sortIdx] + " " + DOWN_ARROW
					sortAsc = false
					sort.Slice(myPage.ProcTable.Rows, sortFuncs[sortIdx])

				// Disable Sort
				case "0":
					myPage.ProcTable.Header = append([]string{}, header...)
					sortIdx = -1
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
				myPage.ProcTable.CursorColor = selectedStyle
				procData := getData(data)
				myPage.ProcTable.Rows = procData
				if sortIdx != -1 {
					sort.Slice(myPage.ProcTable.Rows, sortFuncs[sortIdx])
				}
				on.Do(updateUI)
			}

		case <-tick: // Update page with new values
			if killSelected {
				exists, _ := proc.PidExists(pidToKill)
				if !exists {
					runAllProc = true
					killSelected = false
					myPage.ProcTable.CursorColor = selectedStyle
					updateProcs()
				}
			} else {
				myPage.ProcTable.CursorColor = selectedStyle
			}
			if !helpVisible {
				ui.Render(myPage.Grid)
			}
		}
	}
}
