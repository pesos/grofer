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
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	ui "github.com/gizak/termui/v3"
	"github.com/pesos/grofer/src/display/misc"
	info "github.com/pesos/grofer/src/general"

	"github.com/pesos/grofer/src/container"
	"github.com/pesos/grofer/src/utils"
)

var runProc = true
var helpVisible = false
var errorVisible = false

var sortIdx = -1
var sortAsc = false

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

var header = []string{
	"ID",
	"Image",
	"Name",
	"Status",
	"State",
	"CPU",
	"Memory",
	"Net I/O",
	"Block I/O",
}

// OverallVisuals provides the UI for overall container metrics
func OverallVisuals(ctx context.Context, cli *client.Client, dataChannel chan container.ContainerMetrics, refreshRate uint64) error {

	if err := ui.Init(); err != nil {
		return err
	}

	defer ui.Close()

	var on sync.Once

	// create widgets for help and error
	var help *misc.HelpMenu = misc.NewHelpMenu()
	misc.SelectHelpMenu("cont")
	var errorBox *misc.ErrorBox = misc.NewErrorBox()

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
		errorBox.Resize(w, h)
		if helpVisible {
			ui.Clear()
			ui.Render(help)
		} else if errorVisible {
			ui.Clear()
			ui.Render(errorBox)
		} else {
			ui.Render(myPage.Grid)
		}
	}

	updateDetails := func(data container.ContainerMetrics) {
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

		if sortIdx != -1 {
			utils.SortData(myPage.DetailsTable.Rows, sortIdx, sortAsc, "CONTAINER")
		}
	}

	updateUI() // Initialize empty UI

	uiEvents := ui.PollEvents()
	t := time.NewTicker(time.Duration(refreshRate) * time.Millisecond)
	tick := t.C

	previousKey := ""

	selectedStyle := ui.ColorCyan
	actionStyle := ui.ColorMagenta

	cid := ""
	actionSelected := ""
	actions := map[string]string{
		"P": "pause",
		"U": "unpause",
		"R": "restart",
		"S": "stop",
		"K": "kill",
		"X": "remove",
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": //q or Ctrl-C to quit
				close(dataChannel)
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
			} else if errorVisible {
				switch e.ID {
				case "<Escape>":
					errorVisible = false
					updateUI()
				case "?":
					helpVisible = true
					updateUI()
				}
			} else {
				if actionSelected == "" {
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

					// Sort Ascending
					case "1", "2", "3", "4", "5", "6", "7":
						myPage.DetailsTable.Header = append([]string{}, header...)
						idx, _ := strconv.Atoi(e.ID)
						sortIdx = idx - 1
						myPage.DetailsTable.Header[sortIdx] = header[sortIdx] + " " + UP_ARROW
						sortAsc = true
						utils.SortData(myPage.DetailsTable.Rows, sortIdx, sortAsc, "CONTAINER")

					// Sort Descending
					case "<F1>", "<F2>", "<F3>", "<F4>", "<F5>", "<F6>", "<F7>":
						myPage.DetailsTable.Header = append([]string{}, header...)
						idx, _ := strconv.Atoi(e.ID[2:3])
						sortIdx = idx - 1
						myPage.DetailsTable.Header[sortIdx] = header[sortIdx] + " " + DOWN_ARROW
						sortAsc = false
						utils.SortData(myPage.DetailsTable.Rows, sortIdx, sortAsc, "CONTAINER")

					// Disable Sort
					case "0":
						myPage.DetailsTable.Header = append([]string{}, header...)
						sortIdx = -1

					// Container Selction
					case "P", "U", "S", "R", "K", "X":
						if myPage.DetailsTable.SelectedRow < len(myPage.DetailsTable.Rows) {
							cid = myPage.DetailsTable.Rows[myPage.DetailsTable.SelectedRow][0]

							runProc = false
							actionSelected = actions[e.ID]
							myPage.DetailsTable.CursorColor = actionStyle
						}
					}
				} else {

					var err error = nil
					actionSet := true

					switch e.ID {
					case "<Escape>":
						if actionSelected != "" {
							runProc = true
							actionSelected = ""
							myPage.DetailsTable.CursorColor = selectedStyle
							actionSet = false
						}

					// Pause Action
					case "P":
						if actionSelected == "pause" {
							err = cli.ContainerPause(ctx, cid)
							if err == nil {
								err = container.ContainerWait(ctx, cli, cid, "paused")
							} else {
								misc.SetErrorString(fmt.Sprintf("Error pausing container with ID: %s", cid))
							}
						}

					// Unpause Action
					case "U":
						if actionSelected == "unpause" {
							err = cli.ContainerUnpause(ctx, cid)
							if err == nil {
								err = container.ContainerWait(ctx, cli, cid, "running")
							} else {
								misc.SetErrorString(fmt.Sprintf("Error un-pausing container with ID: %s", cid))
							}
						}

					// Restart Action
					case "R":
						if actionSelected == "restart" {
							err = cli.ContainerRestart(ctx, cid, nil)
							if err == nil {
								err = container.ContainerWait(ctx, cli, cid, "running")
							} else {
								misc.SetErrorString(fmt.Sprintf("Error restarting container with ID: %s", cid))
							}
						}

					// Stop Action
					case "S":
						if actionSelected == "stop" {
							err = cli.ContainerStop(ctx, cid, nil)
							if err == nil {
								err = container.ContainerWait(ctx, cli, cid, "exited")
							} else {
								misc.SetErrorString(fmt.Sprintf("Error stopping container with ID: %s", cid))
							}
						}

					// Kill action
					case "K":
						if actionSelected == "kill" {
							err = cli.ContainerKill(ctx, cid, "")
							if err == nil {
								err = container.ContainerWait(ctx, cli, cid, "exited")
							} else {
								misc.SetErrorString(fmt.Sprintf("Error killing container with ID: %s", cid))
							}
						}

					// Remove action
					case "X":
						if actionSelected == "remove" {
							err = cli.ContainerRemove(ctx, cid, types.ContainerRemoveOptions{
								RemoveVolumes: true,
								Force:         true,
							})
							if err == nil {
								container.ContainerWait(ctx, cli, cid, "removed")
							} else {
								misc.SetErrorString(fmt.Sprintf("Error removing container with ID: %s", cid))
							}
						}
					}

					if actionSet {
						<-dataChannel
						data, _ := container.GetOverallMetrics(ctx, cli, true)
						updateDetails(data)
						updateUI()

						if err != nil {
							errorVisible = true
							updateUI()
						} else {
							errorVisible = false
						}
					}

					myPage.DetailsTable.CursorColor = selectedStyle

					runProc = true
					actionSelected = ""
				}

				if !errorVisible {
					ui.Render(myPage.Grid)
				}
				if previousKey == "g" {
					previousKey = ""
				} else {
					previousKey = e.ID
				}
			}

		case data := <-dataChannel:
			if runProc {
				updateDetails(data)
				on.Do(updateUI)
			}

		case <-tick:
			if !helpVisible && !errorVisible {
				ui.Render(myPage.Grid)
			}
		}
	}

}
