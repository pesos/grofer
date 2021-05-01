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
	"strings"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	h "github.com/pesos/grofer/src/display/misc"
	info "github.com/pesos/grofer/src/general"

	"github.com/pesos/grofer/src/container"
	"github.com/pesos/grofer/src/utils"
)

func getContainers(metrics []container.PerContainerMetrics, sizes []int) []string {
	rows := []string{}

	for _, metric := range metrics {
		row := " "
		row += metric.ID + strings.Repeat(" ", sizes[0]-len(metric.ID)) + " \r"

		if len(metric.Image) >= sizes[1] {
			metric.Image = metric.Image[:sizes[1]]
		}
		row += metric.Image + strings.Repeat(" ", sizes[1]-len(metric.Image)) + " \r"

		metric.Name = strings.TrimLeft(metric.Name, "/")
		if len(metric.Name) >= sizes[2] {
			metric.Name = metric.Name[:sizes[2]]
		}
		row += metric.Name + strings.Repeat(" ", sizes[2]-len(metric.Name)) + " \r"

		if len(metric.Status) >= sizes[3] {
			metric.Status = metric.Status[:sizes[3]]
		}
		row += metric.Status + strings.Repeat(" ", sizes[3]-len(metric.Status)) + " \r"

		if len(metric.State) >= sizes[4] {
			metric.State = metric.State[:sizes[4]]
		}
		row += metric.State + strings.Repeat(" ", sizes[4]-len(metric.State)) + " \r"

		cpu := fmt.Sprintf("%.1f%%", metric.Cpu)
		if len(cpu) >= sizes[5] {
			cpu = cpu[:sizes[5]]
		}
		row += cpu + strings.Repeat(" ", sizes[5]-len(cpu)) + " \r"

		mem := fmt.Sprintf("%.1f%%", metric.Mem)
		if len(mem) >= sizes[6] {
			mem = mem[:sizes[6]]
		}
		row += mem + strings.Repeat(" ", sizes[6]-len(mem)) + " \r"

		netVals, units := utils.RoundValues(metric.Net.Rx, metric.Net.Tx, true)
		units = strings.Trim(units, " \n\r")
		net := fmt.Sprintf("%.1f%s/%.1f%s", netVals[0], units, netVals[1], units)
		if len(net) >= sizes[7] {
			net = net[:sizes[7]]
		}
		row += net + strings.Repeat(" ", sizes[7]-len(net)) + " \r"

		blkVals, units := utils.RoundValues(float64(metric.Blk.Read), float64(metric.Blk.Write), true)
		units = strings.Trim(units, " \n\r")
		blk := fmt.Sprintf("%.2f%s/%.2f%s", blkVals[0], units, blkVals[1], units)
		if len(blk) >= sizes[8] {
			blk = blk[:sizes[8]]
		}
		row += blk

		rows = append(rows, row)
	}

	return rows
}

var runProc = true
var helpVisible = false

func OverallVisuals(ctx context.Context, dataChannel chan container.ContainerMetrics, refreshRate uint64) error {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	var on sync.Once

	var help *h.HelpMenu = h.NewHelpMenu()
	h.SelectHelpMenu("cont")

	// Create new page
	myPage := NewOverallContainerPage()

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
	selectedStyle := ui.NewStyle(ui.ColorYellow, ui.ColorClear, ui.ModifierBold)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": //q or Ctrl-C to quit
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
			}

		case data := <-dataChannel:
			myPage.BodyList.SelectedRowStyle = selectedStyle
			if runProc {
				// update cpu %
				myPage.CPUChart.Percent = int(data.TotalCPU)

				// update mem %
				myPage.MemChart.Percent = int(data.TotalMem)

				// update Net RX and TX
				netVals, units := utils.RoundValues(data.TotalNet.Rx, data.TotalNet.Tx, true)
				myPage.NetChart.Data = netVals
				myPage.NetChart.Title = " Net I/O " + units

				//update page faults
				blkVals, units := utils.RoundValues(float64(data.TotalBlk.Read), float64(data.TotalBlk.Write), true)
				myPage.BlkChart.Data = blkVals
				myPage.BlkChart.Title = " Block I/O " + units

				myPage.BodyList.Rows = getContainers(data.PerContainer, myPage.HeadingTable.ColumnWidths)

				on.Do(updateUI)
			}

		case <-tick:
			if !helpVisible {
				ui.Render(myPage.Grid)
			}
		}
	}

}
