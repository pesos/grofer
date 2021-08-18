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

package general

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/metrics/general"
	"github.com/pesos/grofer/pkg/sink/tui/misc"
	"github.com/pesos/grofer/pkg/utils"
	viz "github.com/pesos/grofer/pkg/utils/visualization"
)

var (
	run = true
)

// RenderCharts handles plotting graphs and charts for system stats in general.
func RenderCharts(ctx context.Context, dataChannel chan general.AggregatedMetrics, refreshRate uint64) error {
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var on sync.Once
	var totalBytesRecv float64
	var totalBytesSent float64
	var help *misc.HelpMenu = misc.NewHelpMenu().ForCommand(misc.RootCommand)

	// Get number of cores in machine
	numCores := runtime.NumCPU()

	// Create new page
	myPage := NewPage(numCores)

	var scrollWidget viz.ScrollableWidget = myPage.DiskChart
	utitlitySelected := ""
	previousKey := ""

	// Pause to pause updating data
	pause := func() {
		run = !run
	}

	updateUI := func() {

		// Get Terminal Dimensions and clear the UI
		w, h := ui.TerminalDimensions()

		ui.Clear()

		switch utitlitySelected {
		case "HELP":
			help.Resize(w, h)
			ui.Render(help)

		default:
			myPage.Grid.SetRect(0, 0, w, h)
			ui.Render(myPage.Grid)
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

		case e := <-uiEvents: // For keyboard events
			switch e.ID {
			case "q", "<C-c>": // q or Ctrl-C to quit
				return core.ErrCanceledByUser

			case "<Resize>":
				updateUI()

			case "<Escape>":
				if utitlitySelected == "HELP" {
					scrollWidget.DisableCursor()
					scrollWidget = myPage.DiskChart
					scrollWidget.EnableCursor()
					utitlitySelected = ""
				}

			case "p":
				pause()

			case "?":
				scrollWidget.DisableCursor()
				scrollWidget = help.Table
				scrollWidget.EnableCursor()
				utitlitySelected = "HELP"

			// handle table navigations
			case "j", "<Down>":
				scrollWidget.ScrollDown()

			case "k", "<Up>":
				scrollWidget.ScrollUp()

			case "<C-d>":
				scrollWidget.ScrollHalfPageDown()

			case "<C-u>":
				scrollWidget.ScrollHalfPageUp()

			case "<C-f>":
				scrollWidget.ScrollPageDown()

			case "<C-b>":
				scrollWidget.ScrollPageUp()

			case "g":
				if previousKey == "g" {
					scrollWidget.ScrollTop()
				}

			case "<Home>":
				scrollWidget.ScrollTop()

			case "G", "<End>":
				scrollWidget.ScrollBottom()

			// handle table switching
			case "<Left>", "h":
				if utitlitySelected != "HELP" {
					scrollWidget.DisableCursor()
					scrollWidget = myPage.SwitchTableLeft(myPage.cpuTableVisible)
					scrollWidget.EnableCursor()
				}

			case "<Right>", "l":
				if utitlitySelected != "HELP" {
					scrollWidget.DisableCursor()
					scrollWidget = myPage.SwitchTableRight(myPage.cpuTableVisible)
					scrollWidget.EnableCursor()
				}

			// handle actions
			case "t":
				scrollWidget.DisableCursor()
				scrollWidget = myPage.ToggleCPUWidget()
				scrollWidget.EnableCursor()
			}

			updateUI()
			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			if run {
				switch data.FieldSet {

				case "CPU": // Update CPU stats
					avgLoad := 0.0
					myPage.CPUGauge.Labels = nil
					// Individual line charts for each CPU core when < 8
					for _, x := range data.CpuStats {
						myPage.CPUGauge.Labels = append(myPage.CPUGauge.Labels, fmt.Sprintf("%.1f%%", x))
						avgLoad += x
					}

					if myPage.cpuTableVisible {
						myPage.CPUTable.Data = data.CpuStats
					} else {
						myPage.CPUGauge.Values = data.CpuStats
					}
					// Generate an Average Graph for CPUs when number of cores > 8
					avgLoad /= float64(numCores)
					if len(myPage.AvgCPUGraph.Data["Average CPU Load:"]) > 100 {
						myPage.AvgCPUGraph.Data["Average CPU Load:"] = myPage.AvgCPUGraph.Data["Average CPU Load:"][1:]
					}

					myPage.AvgCPUGraph.Data["Average CPU Load:"] = append(myPage.AvgCPUGraph.Data["Average CPU Load:"], avgLoad)
					myPage.AvgCPUGraph.Labels["Average CPU Load:"] = fmt.Sprintf("%3.2f%%", avgLoad)
					// Change LineColor based on percentage
					if avgLoad > 66.6 {
						myPage.AvgCPUGraph.LineColors["Average CPU Load:"] = ui.ColorRed
					} else if avgLoad > 33.3 {
						myPage.AvgCPUGraph.LineColors["Average CPU Load:"] = ui.ColorYellow
					} else {
						myPage.AvgCPUGraph.LineColors["Average CPU Load:"] = ui.ColorGreen
					}

				case "MEM": // Update Memory stats
					myPage.MemoryChart.MaxVal = data.MemStats[0]
					myPage.MemoryChart.Data = data.MemStats[1:]
					myPage.MemoryChart.Labels = append(myPage.MemoryChart.Labels, fmt.Sprintf("Used: %.2fG/%.2fG", data.MemStats[1], data.MemStats[0]))
					myPage.MemoryChart.Labels = append(myPage.MemoryChart.Labels, fmt.Sprintf("Available: %.2fG/%.2fG", data.MemStats[2], data.MemStats[0]))
					myPage.MemoryChart.Labels = append(myPage.MemoryChart.Labels, fmt.Sprintf("Free: %.2fG/%.2fG", data.MemStats[3], data.MemStats[0]))
					myPage.MemoryChart.Labels = append(myPage.MemoryChart.Labels, fmt.Sprintf("Cached: %.2fG/%.2fG", data.MemStats[4], data.MemStats[0]))

				case "DISK": // Update Disk stats
					myPage.DiskChart.Header = data.DiskStats[0]
					myPage.DiskChart.Rows = data.DiskStats[1:]

				case "TEMP":
					myPage.TemperatureTable.Header = data.TempStats[0]
					myPage.TemperatureTable.Rows = data.TempStats[1:]

				case "NET": // Update Network stats
					var curBytesRecv, curBytesSent float64

					for _, netInterface := range data.NetStats {
						curBytesRecv += netInterface[1]
						curBytesSent += netInterface[0]
					}

					var recentBytesRecv, recentBytesSent float64

					if totalBytesRecv != 0 {
						recentBytesRecv = curBytesRecv - totalBytesRecv
						recentBytesSent = curBytesSent - totalBytesSent

						if int(recentBytesRecv) < 0 {
							recentBytesRecv = 0
						}
						if int(recentBytesSent) < 0 {
							recentBytesSent = 0
						}
						if len(myPage.NetworkChart.Sparklines[0].Data) > 100 {
							myPage.NetworkChart.Sparklines[0].Data = myPage.NetworkChart.Sparklines[0].Data[1:]
						}
						myPage.NetworkChart.Sparklines[0].Data = append(myPage.NetworkChart.Sparklines[0].Data, recentBytesRecv)
						if len(myPage.NetworkChart.Sparklines[1].Data) > 100 {
							myPage.NetworkChart.Sparklines[1].Data = myPage.NetworkChart.Sparklines[1].Data[1:]
						}
						myPage.NetworkChart.Sparklines[1].Data = append(myPage.NetworkChart.Sparklines[1].Data, recentBytesSent)

					}

					totalBytesRecv = curBytesRecv
					totalBytesSent = curBytesSent

					totalData, units := utils.RoundValues(totalBytesRecv, totalBytesSent, true)

					myPage.NetworkChart.Sparklines[0].Title = fmt.Sprintf(" Total RX: %5.1f %s", totalData[0], units)
					myPage.NetworkChart.Sparklines[1].Title = fmt.Sprintf(" Total TX: %5.1f %s", totalData[1], units)

				}
				on.Do(updateUI)
			}

		case <-tick: // Update page with new values
			if utitlitySelected != "HELP" {
				ui.Render(myPage.Grid)
			}
		}
	}
}

func RenderCPUinfo(ctx context.Context, dataChannel chan *general.CPULoad, refreshRate uint64) error {
	var on sync.Once
	var help *misc.HelpMenu = misc.NewHelpMenu().ForCommand(misc.RootCommand)

	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	numCores := runtime.NumCPU()
	myPage := NewCPUPage(numCores)

	var scrollWidget viz.ScrollableWidget = myPage.CPUTable

	previousKey := ""
	utilitySelected := ""

	pause := func() {
		run = !run
	}

	// Re render UI
	updateUI := func() {
		w, h := ui.TerminalDimensions()
		myPage.Grid.SetRect(0, 0, w, h)

		ui.Clear()

		switch utilitySelected {
		case "HELP":
			help.Resize(w, h)
			ui.Render(help)
		default:
			ui.Render(myPage.Grid)
		}
	}

	updateUI()

	uiEvents := ui.PollEvents()
	t := time.NewTicker(time.Duration(refreshRate) * time.Millisecond)
	tick := t.C
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e := <-uiEvents: // For keyboard events
			switch e.ID {
			case "q", "<C-c>": // q or Ctrl-C to quit
				return core.ErrCanceledByUser

			case "<Resize>":
				updateUI()

			case "<Escape>":
				scrollWidget.DisableCursor()
				scrollWidget = myPage.CPUTable
				scrollWidget.EnableCursor()
				utilitySelected = ""

			case "p":
				pause()

			case "?":
				scrollWidget.DisableCursor()
				scrollWidget = help.Table
				scrollWidget.EnableCursor()
				utilitySelected = "HELP"

			// handle table navigations
			case "j", "<Down>":
				scrollWidget.ScrollDown()

			case "k", "<Up>":
				scrollWidget.ScrollUp()

			case "<C-d>":
				scrollWidget.ScrollHalfPageDown()

			case "<C-u>":
				scrollWidget.ScrollHalfPageUp()

			case "<C-f>":
				scrollWidget.ScrollPageDown()

			case "<C-b>":
				scrollWidget.ScrollPageUp()

			case "g":
				if previousKey == "g" {
					scrollWidget.ScrollTop()
				}

			case "<Home>":
				scrollWidget.ScrollTop()

			case "G", "<End>":
				scrollWidget.ScrollBottom()

			}

			updateUI()
			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel: // Update chart values
			if run {
				myPage.UsrChart.Percent = data.Usr
				myPage.NiceChart.Percent = data.Nice
				myPage.SysChart.Percent = data.Sys
				myPage.IowaitChart.Percent = data.Iowait
				myPage.IrqChart.Percent = data.Irq
				myPage.SoftChart.Percent = data.Soft
				myPage.StealChart.Percent = data.Steal
				myPage.IdleChart.Percent = data.Idle

				if numCores > 8 {
					rows := [][]string{}
					for j := 0; j < len(data.CPURates[0]); j++ {
						rows = append(rows, []string{
							data.CPURates[0][j],
							data.CPURates[1][j],
						})
					}

					myPage.CPUTable.Rows = rows
				} else {
					myPage.CPUChart.Rows = data.CPURates
				}

				on.Do(func() {
					w, h := ui.TerminalDimensions()
					ui.Clear()
					myPage.Grid.SetRect(0, 0, w, h)
					ui.Render(myPage.Grid)
				})
			}

		case <-tick:
			if utilitySelected != "HELP" {
				ui.Render(myPage.Grid)
			}
		}
	}
}
