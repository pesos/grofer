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
	page := NewPage(numCores)

	var scrollWidget viz.ScrollableWidget = page.DiskChart
	utitlitySelected := core.None
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
		case core.Help:
			help.Resize(w, h)
			ui.Render(help)

		default:
			page.Grid.SetRect(0, 0, w, h)
			ui.Render(page.Grid)
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
				if utitlitySelected == core.Help {
					scrollWidget.DisableCursor()
					scrollWidget = page.DiskChart
					scrollWidget.EnableCursor()
					utitlitySelected = core.None
				}

			case "p":
				pause()

			case "?":
				scrollWidget.DisableCursor()
				scrollWidget = help.Table
				scrollWidget.EnableCursor()
				utitlitySelected = core.Help

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
				if utitlitySelected != core.Help {
					scrollWidget.DisableCursor()
					scrollWidget = page.switchTableLeft(page.cpuTableVisible)
					scrollWidget.EnableCursor()
				}

			case "<Right>", "l":
				if utitlitySelected != core.Help {
					scrollWidget.DisableCursor()
					scrollWidget = page.switchTableRight(page.cpuTableVisible)
					scrollWidget.EnableCursor()
				}

			// handle actions
			case "t":
				scrollWidget.DisableCursor()
				scrollWidget = page.ToggleCPUWidget()
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

				case "INFO": // Update Info table
					header, rows := data.HostInfo[0], data.HostInfo[1:]
					page.InfoTable.Header = header
					page.InfoTable.Rows = rows

				case "BATTERY": // Update Battery Gauge
					page.BatteryGauge.Title = " Battery % "

					percent := data.BatteryPercent
					page.BatteryGauge.Percent = percent
					switch {
					case percent < 33:
						page.BatteryGauge.BarColor = ui.ColorRed
					case percent < 67:
						page.BatteryGauge.BarColor = ui.ColorYellow
					default:
						page.BatteryGauge.BarColor = ui.ColorGreen
					}

				case "CPU": // Update CPU stats
					if page.cpuTableVisible {
						page.CPUTable.Data = data.CPUStats
					} else {
						for i, percent := range data.CPUStats {
							cpu := fmt.Sprintf("CPU %d", i)
							page.CPUChart.Data[cpu] = append(page.CPUChart.Data[cpu], percent)
							page.CPUChart.Labels[cpu] = fmt.Sprintf("\t%5.2f %%", percent)
						}
					}

				case "MEM": // Update Memory stats
					page.MemoryChart.MaxVal = data.MemStats[0]
					page.MemoryChart.Data = data.MemStats[1:]
					page.MemoryChart.Labels = append(page.MemoryChart.Labels, fmt.Sprintf("Used: %.2fG/%.2fG", data.MemStats[1], data.MemStats[0]))
					page.MemoryChart.Labels = append(page.MemoryChart.Labels, fmt.Sprintf("Available: %.2fG/%.2fG", data.MemStats[2], data.MemStats[0]))
					page.MemoryChart.Labels = append(page.MemoryChart.Labels, fmt.Sprintf("Free: %.2fG/%.2fG", data.MemStats[3], data.MemStats[0]))
					page.MemoryChart.Labels = append(page.MemoryChart.Labels, fmt.Sprintf("Cached: %.2fG/%.2fG", data.MemStats[4], data.MemStats[0]))

				case "DISK": // Update Disk stats
					page.DiskChart.Header = data.DiskStats[0]
					page.DiskChart.Rows = data.DiskStats[1:]

				case "TEMP":
					page.TemperatureTable.Header = data.TempStats[0]
					page.TemperatureTable.Rows = data.TempStats[1:]

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
						if len(page.NetworkChart.Sparklines[0].Data) > 100 {
							page.NetworkChart.Sparklines[0].Data = page.NetworkChart.Sparklines[0].Data[1:]
						}
						page.NetworkChart.Sparklines[0].Data = append(page.NetworkChart.Sparklines[0].Data, recentBytesRecv)
						if len(page.NetworkChart.Sparklines[1].Data) > 100 {
							page.NetworkChart.Sparklines[1].Data = page.NetworkChart.Sparklines[1].Data[1:]
						}
						page.NetworkChart.Sparklines[1].Data = append(page.NetworkChart.Sparklines[1].Data, recentBytesSent)

					}

					totalBytesRecv = curBytesRecv
					totalBytesSent = curBytesSent

					totalData, units := utils.RoundValues(totalBytesRecv, totalBytesSent, true)

					page.NetworkChart.Sparklines[0].Title = fmt.Sprintf(" Total RX: %5.1f %s", totalData[0], units)
					page.NetworkChart.Sparklines[1].Title = fmt.Sprintf(" Total TX: %5.1f %s", totalData[1], units)

				}
				on.Do(updateUI)
			}

		case <-tick: // Update page with new values
			if utitlitySelected != core.Help {
				ui.Render(page.Grid)
			}
		}
	}
}

// RenderCPUinfo displays the CPU info page
func RenderCPUinfo(ctx context.Context, dataChannel chan *general.CPULoad, refreshRate uint64) error {
	var on sync.Once
	var help *misc.HelpMenu = misc.NewHelpMenu().ForCommand(misc.RootCommand)

	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	numCores := runtime.NumCPU()
	page := NewCPUPage(numCores)

	var scrollWidget viz.ScrollableWidget = page.CPUTable

	previousKey := ""
	utilitySelected := core.None

	pause := func() {
		run = !run
	}

	// Re render UI
	updateUI := func() {
		w, h := ui.TerminalDimensions()
		page.Grid.SetRect(0, 0, w, h)

		ui.Clear()

		switch utilitySelected {
		case core.Help:
			help.Resize(w, h)
			ui.Render(help)
		default:
			ui.Render(page.Grid)
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
				scrollWidget = page.CPUTable
				scrollWidget.EnableCursor()
				utilitySelected = core.None

			case "p":
				pause()

			case "?":
				scrollWidget.DisableCursor()
				scrollWidget = help.Table
				scrollWidget.EnableCursor()
				utilitySelected = core.Help

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
				page.UsrChart.Percent = data.Usr
				page.NiceChart.Percent = data.Nice
				page.SysChart.Percent = data.Sys
				page.IowaitChart.Percent = data.Iowait
				page.IrqChart.Percent = data.Irq
				page.SoftChart.Percent = data.Soft
				page.StealChart.Percent = data.Steal
				page.IdleChart.Percent = data.Idle

				if numCores > 8 {
					rows := [][]string{}
					for j := 0; j < len(data.CPURates[0]); j++ {
						rows = append(rows, []string{
							data.CPURates[0][j],
							data.CPURates[1][j],
						})
					}

					page.CPUTable.Rows = rows
				} else {
					page.CPUChart.Rows = data.CPURates
				}

				on.Do(func() {
					w, h := ui.TerminalDimensions()
					ui.Clear()
					page.Grid.SetRect(0, 0, w, h)
					ui.Render(page.Grid)
				})
			}

		case <-tick:
			if utilitySelected != core.Help {
				ui.Render(page.Grid)
			}
		}
	}
}

// RenderBatteryinfo displays the Battery info page
func RenderBatteryinfo(ctx context.Context, dataChannel chan general.BatteryData, refreshRate uint64) error {
	var on sync.Once
	var help *misc.HelpMenu = misc.NewHelpMenu().ForCommand(misc.RootCommand)

	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	page := NewBatteryPage()

	utilitySelected := core.None

	// Pause to pause updating data
	pause := func() {
		run = !run
	}

	// Re render UI
	updateUI := func() {
		w, h := ui.TerminalDimensions()
		page.Grid.SetRect(0, 0, 200, 55)

		ui.Clear()

		switch utilitySelected {
		case core.Help:
			help.Resize(w, h)
			ui.Render(help)
		default:
			ui.Render(page.Grid)
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
				utilitySelected = core.None

			case "p":
				pause()

			case "?":
				utilitySelected = core.Help
			}
			updateUI()
		case data := <-dataChannel:
			if run {
				header, rows := data.Battery[0], data.Battery[1:]
				page.Battery.Header = header
				page.Battery.Rows = rows
				on.Do(updateUI)
			}
		case <-tick:
			if utilitySelected != core.Help {
				ui.Render(page.Grid)
			}
		}
	}

}
