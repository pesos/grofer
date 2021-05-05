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
	"strconv"
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

		// Adjust Memory Bar graph values
		myPage.MemoryChart.BarGap = ((w / 2) - (4 * myPage.MemoryChart.BarWidth)) / 4

		// Adjust CPU Gauge dimensions
		if isCPUSet {
			for i := 0; i < numCores; i++ {
				myPage.CPUCharts[i].SetRect(0, i*height, w/2, (i+1)*height)
			}
		}

		// Adjust Grid dimensions
		myPage.Grid.SetRect(w/2, 0, w, h-heightOffset)

		help.Resize(w, h)

		if helpVisible {
			ui.Clear()
			ui.Render(help)
		} else {
			ui.Render(myPage.Grid)
			for i := 0; i < numCores; i++ {
				ui.Render(myPage.CPUCharts[i])
			}
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
			}

		case data := <-dataChannel:
			if run {
				switch data.FieldSet {

				case "CPU": // Update CPU stats
					for index, rate := range data.CpuStats {
						myPage.CPUCharts[index].Title = " CPU " + strconv.Itoa(index) + " "
						myPage.CPUCharts[index].Percent = int(rate)
					}

				case "MEM": // Update Memory stats
					myPage.MemoryChart.Data = data.MemStats

				case "DISK": // Update Disk stats
					myPage.DiskChart.Rows = data.DiskStats

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

						myPage.NetworkChart.Data["RX"] = append(myPage.NetworkChart.Data["RX"], recentBytesRecv)
						myPage.NetworkChart.Data["TX"] = append(myPage.NetworkChart.Data["TX"], recentBytesSent)
					}

					totalBytesRecv = curBytesRecv
					totalBytesSent = curBytesSent

					totalData, units := utils.RoundValues(totalBytesRecv, totalBytesSent, true)

					myPage.NetworkChart.Labels["RX"] = fmt.Sprintf("Total: %5.1f %s\n", totalData[0], units)
					myPage.NetworkChart.Labels["TX"] = fmt.Sprintf("Total: %5.1f %s\n", totalData[1], units)

				}
				on.Do(updateUI)
			}

		case <-tick: // Update page with new values
			if !helpVisible {
				ui.Render(myPage.Grid)
				for i := 0; i < numCores; i++ {
					ui.Render(myPage.CPUCharts[i])
				}
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

				myPage.CPUChart.Rows = data.CPURates

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
