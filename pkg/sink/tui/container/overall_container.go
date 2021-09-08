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
	"github.com/pesos/grofer/pkg/core"

	containerMetrics "github.com/pesos/grofer/pkg/metrics/container"
	"github.com/pesos/grofer/pkg/sink/tui/misc"
	"github.com/pesos/grofer/pkg/utils"
	viz "github.com/pesos/grofer/pkg/utils/visualization"
)

// OverallVisuals provides the UI for overall container metrics
func OverallVisuals(ctx context.Context, cli *client.Client, all bool, dataChannel chan containerMetrics.OverallMetrics, refreshRate uint64) error {
	if err := ui.Init(); err != nil {
		return err
	}

	defer ui.Close()

	var on sync.Once

	// create widgets for help, actions and error
	var help *misc.HelpMenu = misc.NewHelpMenu().ForCommand(misc.ContainerCommand)
	var errorBox *misc.ErrorBox = misc.NewErrorBox()
	var actions *misc.ActionTable = misc.NewActionTable()

	// Create new page and select table
	page := newOverallContainerPage()
	var scrollableWidget viz.ScrollableWidget = page.DetailsTable
	scrollableWidget.EnableCursor()
	utilitySelected := core.None

	// variables for sorting
	sortIdx := -1
	sortAsc := false
	header := []string{
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

	// variables to pause UI rendering
	runProc := true
	pause := func() {
		runProc = !runProc
	}

	previousKey := ""

	selectedStyle := ui.ColorCyan
	actionStyle := ui.ColorMagenta

	cid := ""

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
		case core.Help:
			help.Resize(w, h)
			ui.Render(help)

		case core.Error:
			errorBox.Resize(w, h)
			ui.Render(errorBox)

		case core.Action:
			page.DetailsTable.CursorColor = actionStyle
			actions.SetRect(0, 0, w/6, h)
			page.Grid.SetRect(w/6, 0, w, h)
			ui.Render(actions)
			ui.Render(page.Grid)

		default:
			page.DetailsTable.CursorColor = selectedStyle
			ui.Render(page.Grid)
		}
	}

	updateDetails := func(data containerMetrics.OverallMetrics) {
		// update cpu %
		page.CPUChart.Percent = int(data.TotalCPU)

		// update mem %
		page.MemChart.Percent = int(data.TotalMem)

		// update Net RX and TX
		netVals, units := utils.RoundValues(data.TotalNet.Rx, data.TotalNet.Tx, true)
		page.NetChart.Data = netVals
		page.NetChart.Title = " Net I/O " + units

		// update Block IO
		blkVals, units := utils.RoundValues(float64(data.TotalBlk.Read), float64(data.TotalBlk.Write), true)
		page.BlkChart.Data = blkVals
		page.BlkChart.Title = " Block I/O " + units

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
				fmt.Sprintf("%.2f%%", c.CPU),
				fmt.Sprintf("%.2f%%", c.Mem),
				net,
				blk,
			})
		}

		page.DetailsTable.Rows = containerData

		if sortIdx != -1 {
			utils.SortData(page.DetailsTable.Rows, sortIdx, sortAsc, "CONTAINER")
		}
	}

	updateUI() // Initialize empty UI

	uiEvents := ui.PollEvents()
	t := time.NewTicker(time.Duration(refreshRate) * time.Millisecond)
	tick := t.C

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-uiEvents:

			switch e.ID {
			case "q", "<C-c>":
				return core.ErrCanceledByUser

			case "<Resize>":
				updateUI()

			case "<Escape>":
				utilitySelected = core.None
				scrollableWidget.DisableCursor()
				scrollableWidget = page.DetailsTable
				scrollableWidget.EnableCursor()
				updateUI()

			case "?":
				scrollableWidget.DisableCursor()
				scrollableWidget = help.Table
				scrollableWidget.EnableCursor()
				utilitySelected = core.Help
				updateUI()

			case "p":
				pause()

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

			// Container Action Selction
			case "<Enter>":
				if utilitySelected == core.None {
					if page.DetailsTable.SelectedRow < len(page.DetailsTable.Rows) {
						// get CID from the data
						cid = page.DetailsTable.Rows[page.DetailsTable.SelectedRow][0]

						runProc = false

						// open the signal selector
						utilitySelected = core.Action
						scrollableWidget = actions.Table
						scrollableWidget.EnableCursor()
					}
				} else if utilitySelected == core.Action {
					var err error

					actionSelected := actions.SelectedAction()

					switch actionSelected {
					// Pause Action
					case "PAUSE":
						err = cli.ContainerPause(ctx, cid)
						if err == nil {
							err = containerMetrics.Wait(ctx, cli, cid, "paused")
						} else {
							errorBox.SetErrorString(fmt.Sprintf("Error pausing container with ID: %s", cid), err)
						}

					// Unpause Action
					case "UNPAUSE":
						err = cli.ContainerUnpause(ctx, cid)
						if err == nil {
							err = containerMetrics.Wait(ctx, cli, cid, "running")
						} else {
							errorBox.SetErrorString(fmt.Sprintf("Error un-pausing container with ID: %s", cid), err)
						}

					// Restart Action
					case "RESTART":
						err = cli.ContainerRestart(ctx, cid, nil)
						if err == nil {
							err = containerMetrics.Wait(ctx, cli, cid, "running")
						} else {
							errorBox.SetErrorString(fmt.Sprintf("Error restarting container with ID: %s", cid), err)
						}

					// Stop action
					case "STOP":
						err = cli.ContainerStop(ctx, cid, nil)
						if err == nil {
							err = containerMetrics.Wait(ctx, cli, cid, "exited")
						} else {
							errorBox.SetErrorString(fmt.Sprintf("Error stopping container with ID: %s", cid), err)
						}

					// Kill action
					case "KILL":
						err = cli.ContainerKill(ctx, cid, "")
						if err == nil {
							err = containerMetrics.Wait(ctx, cli, cid, "exited")
						} else {
							errorBox.SetErrorString(fmt.Sprintf("Error killing container with ID: %s", cid), err)
						}

					// Remove action
					case "REMOVE":
						err = cli.ContainerRemove(ctx, cid, types.ContainerRemoveOptions{
							RemoveVolumes: true,
							Force:         true,
						})
						if err == nil {
							err = containerMetrics.Wait(ctx, cli, cid, "removed")
						} else {
							errorBox.SetErrorString(fmt.Sprintf("Error removing container with ID: %s", cid), err)
						}
					}

					// Flush out stale data
					<-dataChannel
					data, _ := containerMetrics.GetOverallMetrics(ctx, cli, all)
					updateDetails(data)

					// Display error box if action failed/timed out
					if err != nil {
						errorBox.SetErrorString("Timeout Error", err)
						utilitySelected = core.Error
						scrollableWidget.DisableCursor()
						scrollableWidget = errorBox.Table
					} else {
						utilitySelected = core.None
						scrollableWidget.DisableCursor()
						scrollableWidget = page.DetailsTable
						scrollableWidget.EnableCursor()
					}

					updateUI()

					runProc = true
				}

			// Handle sorting

			// Sort Ascending
			case "1", "2", "3", "4", "5", "6", "7":
				if utilitySelected == core.None {
					page.DetailsTable.Header = append([]string{}, header...)
					idx, _ := strconv.Atoi(e.ID)
					sortIdx = idx - 1
					page.DetailsTable.Header[sortIdx] = header[sortIdx] + " " + viz.UpArrow
					sortAsc = true
					utils.SortData(page.DetailsTable.Rows, sortIdx, sortAsc, "CONTAINER")
				}

			// Sort Descending
			case "<F1>", "<F2>", "<F3>", "<F4>", "<F5>", "<F6>", "<F7>":
				if utilitySelected == core.None {
					page.DetailsTable.Header = append([]string{}, header...)
					idx, _ := strconv.Atoi(e.ID[2:3])
					sortIdx = idx - 1
					page.DetailsTable.Header[sortIdx] = header[sortIdx] + " " + viz.DownArrow
					sortAsc = false
					utils.SortData(page.DetailsTable.Rows, sortIdx, sortAsc, "CONTAINER")
				}

			// Disable Sort
			case "0":
				page.DetailsTable.Header = append([]string{}, header...)
				sortIdx = -1
			}

			updateUI()
			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			if runProc {
				updateDetails(data)
				on.Do(updateUI)
			}

		case <-tick:
			if utilitySelected == core.None {
				ui.Render(page.Grid)
			}
		}
	}

}
