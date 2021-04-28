package container

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	info "github.com/pesos/grofer/src/general"

	"github.com/pesos/grofer/src/container"
	"github.com/pesos/grofer/src/utils"
)

func getContainers(metrics []container.PerContainerMetrics) []string {
	rows := []string{}

	for _, metric := range metrics {
		row := fmt.Sprintf("%s %s %s %s %s %.2f %% %.2f %% %.2f/%.2f %d/%d",
			metric.ContainerID,
			metric.Image,
			metric.Name,
			metric.Status,
			metric.State,
			metric.Cpu,
			metric.Mem,
			metric.Net.Rx,
			metric.Net.Tx,
			metric.Blk.Read,
			metric.Blk.Write,
		)

		rows = append(rows, row)
	}

	return rows
}

var runProc = true
var helpVisible = false

func OverallVisuals(ctx context.Context, dataChannel chan container.ContainerMetrics, refreshRate uint64) error {

	defer ui.Close()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	var on sync.Once

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

		// TODO: HELP SECTION
		// help.Resize(w, h)
		if helpVisible {
			// 	ui.Clear()
			// 	ui.Render(help)
		} else {
			ui.Render(myPage.Grid)
		}
	}

	updateUI() // Initialize empty UI

	uiEvents := ui.PollEvents()
	tick := time.Tick(time.Duration(refreshRate) * time.Millisecond)

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
				// switch e.ID {
				// case "?":
				// 	updateUI()
				// case "<Escape>":
				// 	helpVisible = false
				// 	updateUI()
				// case "j", "<Down>":
				// 	help.List.ScrollDown()
				// 	ui.Render(help)
				// case "k", "<Up>":
				// 	help.List.ScrollUp()
				// 	ui.Render(help)
				// }
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

				// update proc info
				myPage.HeadingTable.Rows = [][]string{{
					" ID",
					" Image",
					" Name",
					" Status",
					" State",
					" CPU",
					" Memory",
					" Net I/O",
					" Block I/O ",
				}}

				myPage.BodyList.Rows = getContainers(data.PerContainer)

				on.Do(updateUI)
			}

		case <-tick:
			if !helpVisible {
				ui.Render(myPage.Grid)
			}
		}
	}

}
