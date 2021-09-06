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

package container

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/sink/tui/misc"
	"github.com/pesos/grofer/pkg/utils"

	"github.com/pesos/grofer/pkg/metrics/container"
	viz "github.com/pesos/grofer/pkg/utils/visualization"
)

// PerContainerVisuals provides the UI for per container metrics
func PerContainerVisuals(ctx context.Context, dataChannel chan container.PerContainerMetrics, refreshRate uint64) error {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	var on sync.Once
	var help *misc.HelpMenu = misc.NewHelpMenu().ForCommand(misc.PerContainerCommand)

	// Create new page
	page := newPerContainerPage()

	var scrollableWidget viz.ScrollableWidget = page.DetailsTable
	scrollableWidget.EnableCursor()
	tableMap := map[string]*viz.Table{
		"1": page.DetailsTable,
		"2": page.MountTable,
		"3": page.NetworkTable,
		"4": page.CPUUsageTable,
		"5": page.PortMapTable,
		"6": page.ProcTable,
	}

	utilitySelected := ""

	// variables to pause UI rendering
	runProc := true
	pause := func() {
		runProc = !runProc
	}

	updateUI := func() {

		// Get Terminal Dimensions and clear the UI
		w, h := ui.TerminalDimensions()

		// Adjust Blk chart Bar graph values
		page.BlkChart.BarGap = ((w / 4) - (2 * page.BlkChart.BarWidth)) / 2

		// Adjust Net chart Bar graph values
		page.NetChart.BarGap = ((w / 4) - (2 * page.NetChart.BarWidth)) / 2

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

	updateUI() // Initialize empty UI

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
				scrollableWidget = page.DetailsTable
				scrollableWidget.EnableCursor()
				updateUI()

			// handle table selection
			case "1", "2", "3", "4", "5", "6":
				if utilitySelected == "" {
					scrollableWidget.DisableCursor()
					scrollableWidget = tableMap[e.ID]
					scrollableWidget.EnableCursor()
				}

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
			// page.BodyList.SelectedRowStyle = selectedStyle
			if runProc {
				// update cpu %
				page.CPUChart.Percent = int(data.CPU)

				// update mem %
				page.MemChart.Percent = int(data.Mem)

				// update Net RX and TX
				netVals, units := utils.RoundValues(data.Net.Rx, data.Net.Tx, true)
				page.NetChart.Data = netVals
				page.NetChart.Title = " Net I/O " + units

				// update Block IO
				blkVals, units := utils.RoundValues(float64(data.Blk.Read), float64(data.Blk.Write), true)
				page.BlkChart.Data = blkVals
				page.BlkChart.Title = " Block I/O " + units

				// update details table
				page.DetailsTable.Header = []string{"Name", data.Name}
				page.DetailsTable.Rows = [][]string{
					{"Image", data.Image},
					{"ID", data.ID},
					{"Status", data.Status},
					{"State", data.State},
					{"Pid", data.Pid},
				}

				// update mount volumes table
				mountData := [][]string{}
				for _, m := range data.Mounts {
					mountData = append(mountData, []string{
						m.Src,
						m.Dst,
						m.Mode,
					})
				}
				page.MountTable.Rows = mountData

				// update network settings table
				netData := [][]string{}
				for _, n := range data.NetInfo {
					netData = append(netData, []string{
						n.Name,
						n.Driver,
						n.IP,
						strconv.FormatBool(n.Ingress),
					})
				}
				page.NetworkTable.Rows = netData

				// update per cpu table
				cpuData := [][]string{}
				for i, c := range data.PerCPU {
					cpuData = append(cpuData, []string{
						fmt.Sprintf("CPU %d", i),
						c,
					})
				}
				page.CPUUsageTable.Rows = cpuData

				// Update port map table
				portData := [][]string{}
				for _, p := range data.PortMap {
					portData = append(portData, []string{
						p.IP,
						fmt.Sprintf("%d", p.Host),
						fmt.Sprintf("%d", p.Container),
						p.Protocol,
					})
				}
				page.PortMapTable.Rows = portData

				// Update proc table
				procData := [][]string{}
				for _, p := range data.Procs {
					procData = append(procData, []string{
						p.PID,
						p.UID,
						p.CMD,
					})
				}
				page.ProcTable.Rows = procData

				on.Do(updateUI)
			}

		case <-tick:
			if utilitySelected == "" {
				ui.Render(page.Grid)
			}
		}
	}

}
