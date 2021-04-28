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

func getContainers(metrics []container.PerContainerMetrics) []string {
	rows := []string{}

	for _, metric := range metrics {
		row := metric.ContainerID + strings.Repeat(" ", 15-len(metric.ContainerID)+2)

		if len(metric.Image) > 15 {
			metric.Image = metric.Image[:15]
		}
		row += metric.Image + strings.Repeat(" ", 16-len(metric.Image)+1)

		metric.Name = strings.ReplaceAll(metric.Name, "/", "")
		if len(metric.Name) > 19 {
			metric.Name = metric.Name[:19]
		}
		row += metric.Name + strings.Repeat(" ", 20-len(metric.Name)+1)

		if len(metric.Status) > 14 {
			metric.Status = metric.Status[:14]
		}
		row += metric.Status + strings.Repeat(" ", 15-len(metric.Status)+1)

		if len(metric.State) > 14 {
			metric.State = metric.State[:14]
		}
		row += metric.State + strings.Repeat(" ", 15-len(metric.State)+1)

		cpu := fmt.Sprintf("%.1f%%", metric.Cpu)
		if len(cpu) > 9 {
			cpu = cpu[:9]
		}
		row += cpu + strings.Repeat(" ", 10-len(cpu)+1)

		mem := fmt.Sprintf("%.1f%%", metric.Mem)
		if len(mem) > 9 {
			mem = mem[:9]
		}
		row += mem + strings.Repeat(" ", 10-len(mem)+1)

		netVals, units := utils.RoundValues(metric.Net.Rx, metric.Net.Tx, true)
		units = strings.Trim(units, " \n\r")
		net := fmt.Sprintf("%.1f%s/%.1f%s", netVals[0], units, netVals[1], units)
		if len(net) > 16 {
			net = net[:16]
		}
		row += net + strings.Repeat(" ", 17-len(net)+1)

		blkVals, units := utils.RoundValues(float64(metric.Blk.Read), float64(metric.Blk.Write), true)
		units = strings.Trim(units, " \n\r")
		blk := fmt.Sprintf("%.2f%s/%.2f%s", blkVals[0], units, blkVals[1], units)
		if len(blk) > 16 {
			blk = blk[:16]
		}
		row += blk + strings.Repeat(" ", 17-len(blk))

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
	h.SelectHelpMenu("proc")

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
