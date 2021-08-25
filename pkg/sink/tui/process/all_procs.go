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
	viz "github.com/pesos/grofer/pkg/utils/visualization"
	proc "github.com/shirou/gopsutil/process"
)

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

	page := newAllProcPage()
	utilitySelected := ""
	var scrollableWidget viz.ScrollableWidget = page.ProcTable

	sortIdx := -1
	sortAsc := false
	header := []string{
		"PID",
		"Command",
		"CPU",
		"Memory",
		"Status",
		"Foreground",
		"Creation Time",
		"Thread Count",
	}

	updateUI := func() {
		// Adjust grid dimesnions
		w, h := ui.TerminalDimensions()
		page.Grid.SetRect(0, 0, w, h)

		// Clear UI
		ui.Clear()

		switch utilitySelected {
		case "HELP":
			help.Resize(w, h)
			ui.Render(help)

		case "KILL":
			signals.SetRect(0, 0, w/6, h)
			page.Grid.SetRect(w/6, 0, w, h)
			ui.Render(signals)
			ui.Render(page.Grid)

		default:
			ui.Render(page.Grid)
		}
	}

	updateUI() // Render empty UI

	// variables to pause UI render
	runAllProc := true
	pause := func() {
		runAllProc = !runAllProc
	}

	uiEvents := ui.PollEvents()
	t := time.NewTicker(time.Duration(refreshRate) * time.Millisecond)
	tick := t.C

	previousKey := ""
	selectedStyle := page.ProcTable.CursorColor
	killingStyle := ui.ColorMagenta
	errorStyle := ui.ColorRed

	// updates process list immediately
	updateProcs := func() {
		if runAllProc {
			procs, err := proc.Processes()
			if err == nil {
				page.ProcTable.Rows = getData(procs)
			}
		}
	}

	// whether a process is selected for killing (UI controls are paused)
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

			case "<Resize>":
				updateUI()

			case "?":
				scrollableWidget.DisableCursor()
				scrollableWidget = help.Table
				scrollableWidget.EnableCursor()
				utilitySelected = "HELP"
				updateUI()
			case "p":
				pause()

			case "<Escape>":
				utilitySelected = ""
				scrollableWidget.DisableCursor()
				scrollableWidget = page.ProcTable
				scrollableWidget.EnableCursor()
				page.ProcTable.CursorColor = selectedStyle
				runAllProc = true
				updateUI()

			// handle table navigations
			case "j", "<Down>":
				scrollableWidget.ScrollDown()

			case "k", "<Up>":
				scrollableWidget.ScrollUp()

			case "<C-d>":
				scrollableWidget.ScrollHalfPageDown()

			case "<C-u>":
				scrollableWidget.ScrollHalfPageUp()

			case "<C-f>":
				scrollableWidget.ScrollPageDown()

			case "<C-b>":
				scrollableWidget.ScrollPageUp()

			case "g":
				if previousKey == "g" {
					scrollableWidget.ScrollTop()
				}

			case "<Home>":
				scrollableWidget.ScrollTop()

			case "G", "<End>":
				scrollableWidget.ScrollBottom()

			// handle actions
			case "K", "<F9>":
				if utilitySelected == "" {
					if page.ProcTable.SelectedRow < len(page.ProcTable.Rows) {
						// get PID from the data
						row := page.ProcTable.Rows[page.ProcTable.SelectedRow]
						pid, err := strconv.Atoi(row[0])
						if err != nil {
							return fmt.Errorf("failed to get PID of process: %v", err)
						}

						// Set pid to kill
						pidToKill = int32(pid)
						runAllProc = false
						page.ProcTable.CursorColor = killingStyle

						// open the signal selector
						utilitySelected = "KILL"
						scrollableWidget = signals.Table
					}
				} else if utilitySelected == "KILL" {
					// get process and kill it
					procToKill, err := proc.NewProcess(pidToKill)
					page.ProcTable.CursorColor = selectedStyle
					if err == nil {
						err = procToKill.SendSignal(syscall.SIGTERM)
						if err != nil {
							page.ProcTable.CursorColor = errorStyle
						}
					} else {
						page.ProcTable.CursorColor = errorStyle
					}
					scrollableWidget = page.ProcTable
					scrollableWidget.EnableCursor()
					runAllProc = true
					updateProcs()
				}

			case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				/*
				* The signal selector can be navigated by entering the number beside the
				* desired signal. Double digit numbers are handled by checking the previous
				* key and, if it is among 1,2 and 3, navigate to the corresponding double
				* digit number (as there are currently 31 supported signals).
				* For example, pressing 25 would first navigate to signal 2, then to signal 25
				 */
				if utilitySelected == "KILL" {
					scrollIdx, _ := strconv.Atoi(e.ID)
					if _, checkPrev := map[string]bool{"1": true, "2": true, "3": true}[previousKey]; checkPrev {
						prevIdx, _ := strconv.Atoi(previousKey)
						scrollIdx = 10*prevIdx + scrollIdx
						handledPreviousKey = true
					}
					signals.Table.ScrollToIndex(scrollIdx - 1) // account for 0-indexing
					ui.Render(signals)
				} else {
					switch e.ID {
					// Sort Ascending
					case "1", "2", "3", "4", "5", "6", "7", "8":
						page.ProcTable.Header = append([]string{}, header...)
						idx, _ := strconv.Atoi(e.ID)
						sortIdx = idx - 1
						page.ProcTable.Header[sortIdx] = header[sortIdx] + " " + UP_ARROW
						sortAsc = true
						utils.SortData(page.ProcTable.Rows, sortIdx, sortAsc, "PROCS")

					// Disable Sort
					case "0":
						page.ProcTable.Header = append([]string{}, header...)
						sortIdx = -1
					}
				}

			// Sort Descending
			case "<F1>", "<F2>", "<F3>", "<F4>", "<F5>", "<F6>", "<F7>", "<F8>":
				page.ProcTable.Header = append([]string{}, header...)
				idx, _ := strconv.Atoi(e.ID[2:3])
				sortIdx = idx - 1
				page.ProcTable.Header[sortIdx] = header[sortIdx] + " " + DOWN_ARROW
				sortAsc = false
				utils.SortData(page.ProcTable.Rows, sortIdx, sortAsc, "PROCS")

			case "<Enter>":
				if utilitySelected == "KILL" {
					signalToSend := signals.SelectedSignal()
					procToKill, err := proc.NewProcess(pidToKill)
					page.ProcTable.CursorColor = selectedStyle
					if err == nil {
						err = procToKill.SendSignal(signalToSend)
						if err != nil {
							page.ProcTable.CursorColor = errorStyle
						}
					} else {
						page.ProcTable.CursorColor = errorStyle
					}

					runAllProc = true
					utilitySelected = ""
					updateProcs()
				}
				scrollableWidget = page.ProcTable
				scrollableWidget.EnableCursor()
			}

			updateUI()
			if handledPreviousKey {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			if runAllProc {
				page.ProcTable.CursorColor = selectedStyle
				procData := getData(data)
				page.ProcTable.Rows = procData
				if sortIdx != -1 {
					utils.SortData(page.ProcTable.Rows, sortIdx, sortAsc, "PROCS")
				}
				on.Do(updateUI)
			}

		case <-tick: // Update page with new values
			if utilitySelected == "KILL" {
				exists, _ := proc.PidExists(pidToKill)
				if !exists {
					runAllProc = true
					utilitySelected = ""
					page.ProcTable.CursorColor = selectedStyle
					updateProcs()
				}
			} else {
				page.ProcTable.CursorColor = selectedStyle
			}

			if utilitySelected != "HELP" {
				if utilitySelected == "KILL" {
					ui.Render(signals)
				}
				ui.Render(page.Grid)
			}
		}
	}
}
