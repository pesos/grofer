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
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
)

var isCPUSet = false

var run = true

// RenderCharts handles plotting graphs and charts for system stats in general.
func RenderCharts(endChannel chan os.Signal,
	memChannel chan []float64,
	cpuChannel chan []float64,
	diskChannel chan [][]string,
	netChannel chan map[string][]float64,
	refreshRate int32,
	wg *sync.WaitGroup) {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var totalBytesRecv float64
	var totalBytesSent float64

	// Get number of cores in machine
	numCores := runtime.NumCPU()
	isCPUSet = true

	// Create new page
	myPage := NewPage(numCores)

	// Initialize slices for Network Data
	ipData := make([]float64, 40)
	opData := make([]float64, 40)

	// Pause to pause updating data
	pause := func() {
		run = !run
	}

	uiEvents := ui.PollEvents()
	tick := time.Tick(time.Duration(refreshRate) * time.Millisecond)
	for {
		select {
		case e := <-uiEvents: // For keyboard events
			switch e.ID {
			case "q", "<C-c>": // q or Ctrl-C to quit
				endChannel <- os.Kill
				wg.Done()
				return

			case "s": // s to stop
				pause()
			}

		case data := <-memChannel: // Update memory values
			if run {
				myPage.MemoryChart.Data = data
			}

		case data := <-diskChannel: // Update disk values
			if run {
				myPage.DiskChart.Rows = data
			}

		case data := <-netChannel: // Update network stats & render braille plots
			if run {

				var curBytesRecv, curBytesSent float64

				for _, netInterface := range data {
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

					ipData = ipData[1:]
					opData = opData[1:]

					ipData = append(ipData, recentBytesRecv)
					opData = append(opData, recentBytesSent)
				}

				totalBytesRecv = curBytesRecv
				totalBytesSent = curBytesSent

				titles := make([]string, 2)

				for i := 0; i < 2; i++ {
					if i == 0 {
						titles[i] = fmt.Sprintf("[Total RX](fg:red): %5.1f %s\n", totalBytesRecv/1024, "mB")
					} else {
						titles[i] = fmt.Sprintf("\n[Total TX](fg:green): %5.1f %s", totalBytesSent/1024, "mB")
					}

				}

				myPage.NetPara.Text = titles[0] + titles[1]

				temp := [][]float64{}
				temp = append(temp, ipData)
				temp = append(temp, opData)
				myPage.NetworkChart.Data = temp
			}

		case cpu_data := <-cpuChannel: // Update Gauge map with newer values
			if run {
				for index, rate := range cpu_data {
					myPage.CPUCharts[index].Title = " CPU " + strconv.Itoa(index) + " "
					myPage.CPUCharts[index].Percent = int(rate)
				}
			}

		case <-tick: // Update page with new values
			w, h := ui.TerminalDimensions()
			ui.Clear()

			height := int(h / numCores)
			heightOffset := h - (height * numCores)

			// Adjust Grid dimensions
			myPage.Grid.SetRect(w/2, 0, w, h-heightOffset)

			// Adjust Memory Bar graph values
			myPage.MemoryChart.BarGap = ((w / 2) - (4 * myPage.MemoryChart.BarWidth)) / 4

			// Adjust CPU Gauge dimensions
			if isCPUSet {
				for i := 0; i < numCores; i++ {
					myPage.CPUCharts[i].SetRect(0, i*height, w/2, (i+1)*height)
					ui.Render(myPage.CPUCharts[i])
				}
			}

			ui.Render(myPage.Grid)

		}
	}
}
