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
	"syscall"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/sink/tui/misc"
	"github.com/pesos/grofer/pkg/utils"
	proc "github.com/shirou/gopsutil/process"
)

var runAllProc = true
var helpVisible = false
var sendSignal = false
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

func AllProcVisuals(ctx context.Context, dataChannel chan []*proc.Process, refreshRate uint64) error {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	var on sync.Once
	var signals *misc.SignalTable = misc.NewSignalTable()
	var help *misc.HelpMenu = misc.NewHelpMenu().ForCommand(misc.ProcCommand)

	myPage := NewAllProcsPage()

	updateUI := func() {
		w, h := ui.TerminalDimensions()
		myPage.Grid.SetRect(0, 0, w, h)
		help.Resize(w, h)
		if helpVisible {
			ui.Clear()
			ui.Render(help)
		} else {
			if sendSignal {
				signals.SetRect(0, 0, w/6, h)
				myPage.Grid.SetRect(w/6, 0, w, h)
				ui.Render(signals)
			}
			ui.Render(myPage.Grid)
		}
	}

	updateUI() // Render empty UI

	pauseProc := func() {
		runAllProc = !runAllProc
	}

	uiEvents := ui.PollEvents()
	t := time.NewTicker(time.Duration(refreshRate) * time.Millisecond)
	tick := t.C

	previousKey := ""
	selectedStyle := myPage.ProcTable.CursorColor
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
	var handledPreviousKey bool

	for {
		handledPreviousKey = false
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": //q or Ctrl-C to quit
				return core.ErrCanceledByUser
			case "?":
				helpVisible = !helpVisible
				updateUI()
			case "<Resize>":
				updateUI() // updateUI only during resize event
			}
			if helpVisible { // keybindings for help menu
				switch e.ID {
				case "<Escape>":
					helpVisible = false
					updateUI()
				case "j", "<Down>":
					help.ScrollDown()
					ui.Render(help)
				case "k", "<Up>":
					help.ScrollUp()
					ui.Render(help)
				}
			} else {
				if !killSelected { // keybindings for main proc table
					switch e.ID {
					case "s": //s to pause
						pauseProc()
					case "j", "<Down>":
						myPage.ProcTable.ScrollDown()
					case "k", "<Up>":
						myPage.ProcTable.ScrollUp()
					case "<C-d>":
						myPage.ProcTable.ScrollHalfPageDown()
					case "<C-u>":
						myPage.ProcTable.ScrollHalfPageUp()
					case "<C-f>":
						myPage.ProcTable.ScrollPageDown()
					case "<C-b>":
						myPage.ProcTable.ScrollPageUp()
					case "g":
						if previousKey == "g" {
							myPage.ProcTable.ScrollTop()
							handledPreviousKey = true
						}
					case "<Home>":
						myPage.ProcTable.ScrollTop()
					case "G", "<End>":
						myPage.ProcTable.ScrollBottom()
					case "K", "<F9>":
						if myPage.ProcTable.SelectedRow < len(myPage.ProcTable.Rows) {
							// get PID from the data
							row := myPage.ProcTable.Rows[myPage.ProcTable.SelectedRow]
							pid, err := strconv.Atoi(row[0])
							if err != nil {
								return fmt.Errorf("failed to get PID of process: %v", err)
							}

							// Set pid to kill
							pidToKill = int32(pid)
							runAllProc = false
							killSelected = true
							myPage.ProcTable.CursorColor = killingStyle
							// open the signal selector
							sendSignal = true
							updateUI()
						}
					// Sort Ascending
					case "1", "2", "3", "4", "5", "6", "7", "8":
						myPage.ProcTable.Header = append([]string{}, header...)
						idx, _ := strconv.Atoi(e.ID)
						sortIdx = idx - 1
						myPage.ProcTable.Header[sortIdx] = header[sortIdx] + " " + UP_ARROW
						sortAsc = true
						utils.SortData(myPage.ProcTable.Rows, sortIdx, sortAsc, "PROCS")

					// Sort Descending
					case "<F1>", "<F2>", "<F3>", "<F4>", "<F5>", "<F6>", "<F7>", "<F8>":
						myPage.ProcTable.Header = append([]string{}, header...)
						idx, _ := strconv.Atoi(e.ID[2:3])
						sortIdx = idx - 1
						myPage.ProcTable.Header[sortIdx] = header[sortIdx] + " " + DOWN_ARROW
						sortAsc = false
						utils.SortData(myPage.ProcTable.Rows, sortIdx, sortAsc, "PROCS")

					// Disable Sort
					case "0":
						myPage.ProcTable.Header = append([]string{}, header...)
						sortIdx = -1
					}
				} else { // keybindings for signal menu
					switch e.ID {
					case "<Escape>":
						if killSelected {
							runAllProc = true
							killSelected = false
							myPage.ProcTable.CursorColor = selectedStyle
						}
						sendSignal = false
						updateUI()
					case "K", "<F9>":
						// get process and kill it
						procToKill, err := proc.NewProcess(pidToKill)
						myPage.ProcTable.CursorColor = selectedStyle
						if err == nil {
							err = procToKill.SendSignal(syscall.SIGTERM)
							if err != nil {
								myPage.ProcTable.CursorColor = errorStyle
							}
						} else {
							myPage.ProcTable.CursorColor = errorStyle
						}
						runAllProc = true
						killSelected = false
						updateProcs()
						sendSignal = false
						updateUI()
					case "j", "<Down>":
						signals.Table.ScrollDown()
						ui.Render(signals)
					case "k", "<Up>":
						signals.Table.ScrollUp()
						ui.Render(signals)
					case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
						/*
						* The signal selector can be navigated by entering the number beside the
						* desired signal. Double digit numbers are handled by checking the previous
						* key and, if it is among 1,2 and 3, navigate to the corresponding double
						* digit number (as there are currently 31 supported signals).
						* For example, pressing 25 would first navigate to signal 2, then to signal 25
						 */
						scrollIdx, _ := strconv.Atoi(e.ID)
						if _, checkPrev := map[string]bool{"1": true, "2": true, "3": true}[previousKey]; checkPrev {
							prevIdx, _ := strconv.Atoi(previousKey)
							scrollIdx = 10*prevIdx + scrollIdx
							handledPreviousKey = true
						}
						signals.Table.ScrollToIndex(scrollIdx - 1) // account for 0-indexing
						ui.Render(signals)
					case "<Enter>":
						signalToSend := signals.SelectedSignal()
						procToKill, err := proc.NewProcess(pidToKill)
						myPage.ProcTable.CursorColor = selectedStyle
						if err == nil {
							err = procToKill.SendSignal(signalToSend)
							if err != nil {
								myPage.ProcTable.CursorColor = errorStyle
							}
						} else {
							myPage.ProcTable.CursorColor = errorStyle
						}
						runAllProc = true
						killSelected = false
						updateProcs()
						sendSignal = false
						updateUI()
					}
				}

				ui.Render(myPage.Grid)
				if handledPreviousKey {
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
					utils.SortData(myPage.ProcTable.Rows, sortIdx, sortAsc, "PROCS")
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
				if sendSignal {
					ui.Render(signals)
				}
				ui.Render(myPage.Grid)
			}
		}
	}
}
