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
	"sort"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	h "github.com/pesos/grofer/src/display/misc"
	info "github.com/pesos/grofer/src/general"

	"github.com/pesos/grofer/src/container"
	"github.com/pesos/grofer/src/utils"
)

var runProc = true
var helpVisible = false

// OverallVisuals provides the UI for overall container metrics
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
					myPage.DetailsTable.ScrollDown()
				case "k", "<Up>":
					myPage.DetailsTable.ScrollUp()
				case "<C-d>":
					myPage.DetailsTable.ScrollHalfPageDown()
				case "<C-u>":
					myPage.DetailsTable.ScrollHalfPageUp()
				case "<C-f>":
					myPage.DetailsTable.ScrollPageDown()
				case "<C-b>":
					myPage.DetailsTable.ScrollPageUp()
				case "g":
					if previousKey == "g" {
						myPage.DetailsTable.ScrollTop()
					}
				case "<Home>":
					myPage.DetailsTable.ScrollTop()
				case "G", "<End>":
					myPage.DetailsTable.ScrollBottom()
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
				// update cpu %
				myPage.CPUChart.Percent = int(data.TotalCPU)

				// update mem %
				myPage.MemChart.Percent = int(data.TotalMem)

				// update Net RX and TX
				netVals, units := utils.RoundValues(data.TotalNet.Rx, data.TotalNet.Tx, true)
				myPage.NetChart.Data = netVals
				myPage.NetChart.Title = " Net I/O " + units

				// update Block IO
				blkVals, units := utils.RoundValues(float64(data.TotalBlk.Read), float64(data.TotalBlk.Write), true)
				myPage.BlkChart.Data = blkVals
				myPage.BlkChart.Title = " Block I/O " + units

				// Sort container data
				sort.Slice(data.PerContainer, func(i, j int) bool { return data.PerContainer[i].ID > data.PerContainer[j].ID })

				// update container details table
				containerData := [][]string{}
				for _, c := range data.PerContainer {
					netVals, units := utils.RoundValues(c.Net.Rx, c.Net.Tx, true)
					net := fmt.Sprintf("%.1f%s/%.1f%s", netVals[0], units, netVals[1], units)

					blkVals, units := utils.RoundValues(float64(c.Blk.Read), float64(c.Blk.Write), true)
					blk := fmt.Sprintf("%.2f%s/%.2f%s", blkVals[0], units, blkVals[1], units)
					containerData = append(containerData, []string{
						c.ID,
						c.Image,
						c.Name,
						c.Status,
						c.State,
						fmt.Sprintf("%.2f%%", c.Cpu),
						fmt.Sprintf("%.2f%%", c.Mem),
						net,
						blk,
					})
				}

				myPage.DetailsTable.Rows = containerData

				on.Do(updateUI)
			}

		case <-tick:
			if !helpVisible {
				ui.Render(myPage.Grid)
			}
		}
	}

}
