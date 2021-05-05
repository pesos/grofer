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
	h "github.com/pesos/grofer/src/display/misc"
	info "github.com/pesos/grofer/src/general"

	"github.com/pesos/grofer/src/container"
	"github.com/pesos/grofer/src/utils"
)

// ContainerVisuals provides the UI for per container metrics
func ContainerVisuals(ctx context.Context, dataChannel chan container.PerContainerMetrics, refreshRate uint64) error {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	var selectedTable *utils.Table
	var on sync.Once
	var help *h.HelpMenu = h.NewHelpMenu()
	h.SelectHelpMenu("cont_cid")

	// Create new page
	myPage := NewPerContainerPage()

	tableMap := map[string]*utils.Table{
		"0": nil,
		"1": myPage.MountTable,
		"2": myPage.NetworkTable,
		"3": myPage.CPUUsageTable,
		"4": myPage.PortMapTable,
		"5": myPage.ProcTable,
	}

	pause := func() {
		runProc = !runProc
	}

	updateUI := func() {

		// Get Terminal Dimensions adn clear the UI
		w, h := ui.TerminalDimensions()

		// Adjust Blk chart Bar graph values
		myPage.BlkChart.BarGap = ((w / 4) - (2 * myPage.BlkChart.BarWidth)) / 2

		// Adjust Net chart Bar graph values
		myPage.NetChart.BarGap = ((w / 4) - (2 * myPage.NetChart.BarWidth)) / 2

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
				return info.ErrCanceledByUser
			case "<Resize>":
				updateUI()
			case "?":
				helpVisible = !helpVisible
			case "0", "1", "2", "3", "4", "5":
				if !helpVisible {
					if selectedTable != nil {
						selectedTable.ShowCursor = false
					}
					selectedTable = tableMap[e.ID]
				}
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
			} else if selectedTable != nil {
				selectedTable.ShowCursor = true
				switch e.ID {
				case "j", "<Down>":
					selectedTable.ScrollDown()
				case "k", "<Up>":
					selectedTable.ScrollUp()
				case "<C-d>":
					selectedTable.ScrollHalfPageDown()
				case "<C-u>":
					selectedTable.ScrollHalfPageUp()
				case "<C-f>":
					selectedTable.ScrollPageDown()
				case "<C-b>":
					selectedTable.ScrollPageUp()
				case "g":
					if previousKey == "g" {
						selectedTable.ScrollTop()
					}
				case "<Home>":
					selectedTable.ScrollTop()
				case "G", "<End>":
					selectedTable.ScrollBottom()
				}
				ui.Render(myPage.Grid)
				if previousKey == "g" {
					previousKey = ""
				} else {
					previousKey = e.ID
				}
			} else {
				switch e.ID {
				case "?":
					updateUI()
				case "s": //s to pause
					pause()
				}

				ui.Render(myPage.Grid)
				if previousKey == "g" {
					previousKey = ""
				} else {
					previousKey = e.ID
				}
			}

		case data := <-dataChannel:
			// myPage.BodyList.SelectedRowStyle = selectedStyle
			if runProc {
				// update cpu %
				myPage.CPUChart.Percent = int(data.Cpu)

				// update mem %
				myPage.MemChart.Percent = int(data.Mem)

				// update Net RX and TX
				netVals, units := utils.RoundValues(data.Net.Rx, data.Net.Tx, true)
				myPage.NetChart.Data = netVals
				myPage.NetChart.Title = " Net I/O " + units

				// update Block IO
				blkVals, units := utils.RoundValues(float64(data.Blk.Read), float64(data.Blk.Write), true)
				myPage.BlkChart.Data = blkVals
				myPage.BlkChart.Title = " Block I/O " + units

				// update details table
				myPage.DetailsTable.Header = []string{"Name", data.Name}
				myPage.DetailsTable.Rows = [][]string{
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
				myPage.MountTable.Rows = mountData

				// update network settings table
				netData := [][]string{}
				for _, n := range data.NetInfo {
					netData = append(netData, []string{
						n.Name,
						n.Driver,
						n.Ip,
						strconv.FormatBool(n.Ingress),
					})
				}
				myPage.NetworkTable.Rows = netData

				// update per cpu table
				cpuData := [][]string{}
				for i, c := range data.PerCPU {
					cpuData = append(cpuData, []string{
						fmt.Sprintf("CPU %d", i),
						c,
					})
				}
				myPage.CPUUsageTable.Rows = cpuData

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
				myPage.PortMapTable.Rows = portData

				// Update proc table
				procData := [][]string{}
				for _, p := range data.Procs {
					procData = append(procData, []string{
						p.PID,
						p.UID,
						p.CMD,
					})
				}
				myPage.ProcTable.Rows = procData

				on.Do(updateUI)
			}

		case <-tick:
			if !helpVisible {
				ui.Render(myPage.Grid)
			}
		}
	}

}
