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
	h "github.com/pesos/grofer/src/display/misc"
	info "github.com/pesos/grofer/src/general"
	"github.com/pesos/grofer/src/utils"
)

var isCPUSet = false

var run = true
var helpVisible = false

// RenderCharts handles plotting graphs and charts for system stats in general.
func RenderCharts(ctx context.Context,
	dataChannel chan utils.DataStats,
	refreshRate uint64) error {

	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var on sync.Once
	var totalBytesRecv float64
	var totalBytesSent float64
	var help *h.HelpMenu = h.NewHelpMenu()
	h.SelectHelpMenu("main")

	// Get number of cores in machine
	numCores := runtime.NumCPU()
	isCPUSet = true

	// Create new page
	myPage := NewPage(numCores)

	// Pause to pause updating data
	pause := func() {
		run = !run
	}

	updateUI := func() {

		// Get Terminal Dimensions adn clear the UI
		w, h := ui.TerminalDimensions()

		// Calculate Heigth offset
		height := int(h / numCores)
		heightOffset := h - (height * numCores)
		_ = heightOffset

		if isCPUSet {
			myPage.Grid.SetRect(0, 0, w, h)
		}

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
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e := <-uiEvents: // For keyboard events
			switch e.ID {
			case "q", "<C-c>": // q or Ctrl-C to quit
				return info.ErrCanceledByUser

			case "<Resize>":
				updateUI()

			case "?": // s to stop
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
				}
				if numCores > 8 {
					switch e.ID {
					// Use <Up> and <Down> to scroll CPUTable when number of cores is greater than 8,
					// and use <j> and <k> to scroll the DiskChart
					case "<Down>":
						myPage.CPUTable.ScrollDown()
						ui.Render(myPage.Grid)
					case "<Up>":
						myPage.CPUTable.ScrollUp()
						ui.Render(myPage.Grid)

					case "j":
						myPage.DiskChart.ScrollDown()
						ui.Render(myPage.Grid)

					case "k":
						myPage.DiskChart.ScrollUp()
						ui.Render(myPage.Grid)
					}
				} else {
					switch e.ID {
					// Use <Down><Up>/<j><k> to scroll DiskChart when number of CPUs is less than, equal to 8
					case "j", "<Down>":
						myPage.DiskChart.ScrollDown()
						ui.Render(myPage.Grid)

					case "k", "<Up>":
						myPage.DiskChart.ScrollUp()
						ui.Render(myPage.Grid)
					}
				}
			}

		case data := <-dataChannel:
			if run {
				switch data.FieldSet {

				case "CPU": // Update CPU stats
					avgLoad := 0.0
					// Individual line charts for each CPU core when < 8
					for i, x := range data.CpuStats {
						key := fmt.Sprintf("CPU%d", i)
						if len(myPage.CPUGraph.Data[key]) > 100 {
							myPage.CPUGraph.Data[key] = myPage.CPUGraph.Data[key][1:]
						}
						myPage.CPUGraph.Data[key] = append(myPage.CPUGraph.Data[key], x)
						myPage.CPUGraph.Labels[key] = fmt.Sprintf("%3.0f%%", x)
						avgLoad += x
					}
					// Generate an Average Graph for CPUs when number of cores > 8
					if numCores > 8 {
						myPage.CPUTable.Data = data.CpuStats
						avgLoad /= float64(numCores)
						if len(myPage.AvgCPUGraph.Data["Average CPU Load:"]) > 100 {
							myPage.AvgCPUGraph.Data["Average CPU Load:"] = myPage.AvgCPUGraph.Data["Average CPU Load:"][1:]
						}
						myPage.AvgCPUGraph.Data["Average CPU Load:"] = append(myPage.AvgCPUGraph.Data["Average CPU Load:"], avgLoad)
						myPage.AvgCPUGraph.Labels["Average CPU Load:"] = fmt.Sprintf("%3.2f%%", avgLoad)
						// Change LineColor based on percentage
						if avgLoad > 70.0 {
							myPage.AvgCPUGraph.LineColors["Average CPU Load:"] = ui.ColorRed
						} else if avgLoad > 40.0 {
							myPage.AvgCPUGraph.LineColors["Average CPU Load:"] = ui.ColorYellow
						} else {
							myPage.AvgCPUGraph.LineColors["Average CPU Load:"] = ui.ColorGreen
						}
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
					myPage.TemperatureTable.Rows = data.TempStats

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
			if !helpVisible {
				ui.Render(myPage.Grid)
			}
		}
	}
}

func RenderCPUinfo(ctx context.Context,
	dataChannel chan *info.CPULoad,
	refreshRate uint64) error {

	var on sync.Once
	var help *h.HelpMenu = h.NewHelpMenu()
	h.SelectHelpMenu("main")

	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	numCores := runtime.NumCPU()
	myPage := NewCPUPage(numCores)

	pause := func() {
		run = !run
	}

	// Re render UI
	updateUI := func() {
		w, h := ui.TerminalDimensions()
		myPage.Grid.SetRect(0, 0, w, h)
		help.Resize(w, h)
		if helpVisible {
			ui.Clear()
			ui.Render(help)
		} else {
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
				return info.ErrCanceledByUser
			case "<Resize>":
				updateUI()

			case "?": // s to stop
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
				}
				if numCores > 8 {
					switch e.ID {
					case "j", "<Down>":
						myPage.CPUTable.ScrollDown()
						ui.Render(myPage.Grid)
					case "k", "<Up>":
						myPage.CPUTable.ScrollUp()
						ui.Render(myPage.Grid)
					}
				}
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
			if !helpVisible {
				ui.Render(myPage.Grid)
			}
		}
	}
}
