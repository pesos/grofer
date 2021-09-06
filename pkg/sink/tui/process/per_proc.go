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
	"log"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/metrics/process"
	"github.com/pesos/grofer/pkg/sink/tui/misc"
	"github.com/pesos/grofer/pkg/utils"
	viz "github.com/pesos/grofer/pkg/utils/visualization"
)

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
	var help *misc.HelpMenu = misc.NewHelpMenu().ForCommand(misc.PerProcCommand)

	// Create new page and select default table
	page := newPerProcPage()
	utilitySelected := ""
	var scrollableWidget viz.ScrollableWidget = page.ChildProcsTable
	scrollableWidget.EnableCursor()

	var statusMap = map[string]string{
		"R": "Running",
		"S": "Sleep",
		"Z": "Zombie",
		"T": "Stop",
		"I": "Idle",
		"W": "Wait",
		"L": "Lock",
	}

	// variables to pause UI render
	runProc := true
	pause := func() {
		runProc = !runProc
	}

	updateUI := func() {

		// Get Terminal Dimensions
		w, h := ui.TerminalDimensions()

		// Adjust Memory Stats Bar graph values
		page.MemStatsChart.BarGap = ((w / 2) - (4 * page.MemStatsChart.BarWidth)) / 4

		// Adjust Page Faults Bar graph values
		page.PageFaultsChart.BarGap = ((w / 4) - (2 * page.PageFaultsChart.BarWidth)) / 2

		// Adjust Context Switches Bar graph values
		page.CTXSwitchesChart.BarGap = ((w / 4) - (2 * page.CTXSwitchesChart.BarWidth)) / 2

		// Adjust Grid dimensions
		page.Grid.SetRect(0, 0, w, h)

		// Clear UI
		ui.Clear()
		switch utilitySelected {
		case "HELP":
			help.Resize(w, h)
			ui.Render(help)

		default:
			ui.Render(page.Grid)
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
				scrollableWidget = page.ChildProcsTable
				scrollableWidget.EnableCursor()
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
			}

			updateUI()
			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			if runProc {
				// update ctx switches
				switches, units := utils.RoundValues(float64(data.NumCtxSwitches.Voluntary), float64(data.NumCtxSwitches.Involuntary), false)

				page.CTXSwitchesChart.Data = switches
				page.CTXSwitchesChart.Title = " CTX Switches" + units

				// update cpu %
				page.CPUChart.Percent = int(data.CPUPercent)

				// update mem %
				page.MemChart.Percent = int(data.MemoryPercent)

				// update proc info
				page.PIDTable.Rows = [][]string{
					{"[Name](fg:green)", data.Name},
					{"[Command](fg:green)", data.Exe},
					{"[Status](fg:green)", statusMap[data.Status] + " (" + data.Status + ")"},
					{"[Background](fg:green)", strconv.FormatBool(data.Background)},
					{"[Foreground](fg:green)", strconv.FormatBool(data.Foreground)},
					{"[Running](fg:green)", strconv.FormatBool(data.IsRunning)},
					{"[Creation Time](fg:green)", utils.GetDateFromUnix(data.CreateTime)},
					{"[Nice value](fg:green)", strconv.Itoa(int(data.Nice))},
					{"[Thread count](fg:green)", strconv.Itoa(int(data.NumThreads))},
					{"[Child process count](fg:green)", strconv.Itoa(len(data.Children))},
					{"[Last Update](fg:green)", time.Now().Format("15:04:05")},
				}
				page.PIDTable.Title = " PID: " + strconv.Itoa(int(data.Proc.Pid)) + " "

				//update memory stats
				memData := []float64{utils.GetInMB(data.MemoryInfo.RSS, 1),
					utils.GetInMB(data.MemoryInfo.Data, 1),
					utils.GetInMB(data.MemoryInfo.Stack, 1),
					utils.GetInMB(data.MemoryInfo.Swap, 1),
				}
				page.MemStatsChart.Data = memData

				//update page faults
				faults, units := utils.RoundValues(float64(data.PageFault.MinorFaults), float64(data.PageFault.MajorFaults), false)

				page.PageFaultsChart.Data = faults
				page.PageFaultsChart.Title = " Page Faults" + units
				page.ChildProcsTable.Rows = getChildProcs(data)

				on.Do(updateUI)
			}

		case <-tick:
			if utilitySelected == "" {
				ui.Render(page.Grid)
			}
		}
	}
}
