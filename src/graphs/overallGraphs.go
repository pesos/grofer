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

package graphs

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var isCPUSet = false

var run = true

type mainPage struct {
	Grid         *ui.Grid
	MemoryChart  *widgets.BarChart
	DiskChart    *widgets.Table
	NetworkChart *widgets.Plot
	CPUCharts    []*widgets.Gauge
	NetPara      *widgets.Paragraph
}

func newPage(numCores int) *mainPage {
	page := &mainPage{
		Grid:         ui.NewGrid(),
		MemoryChart:  widgets.NewBarChart(),
		DiskChart:    widgets.NewTable(),
		NetworkChart: widgets.NewPlot(),
		CPUCharts:    make([]*widgets.Gauge, 0),
		NetPara:      widgets.NewParagraph(),
	}
	page.init(numCores)
	return page
}

func (page *mainPage) init(numCores int) {

	// Initialize Bar Graph for Memory Chart
	page.MemoryChart.Title = " Memory (RAM) "
	page.MemoryChart.Labels = []string{"Total", "Available", "Used", "Free"}
	page.MemoryChart.BarWidth = 8
	page.MemoryChart.BarGap = 9
	page.MemoryChart.BarColors = []ui.Color{ui.ColorCyan, ui.ColorGreen}
	page.MemoryChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	page.MemoryChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
	page.MemoryChart.BorderStyle.Fg = ui.ColorCyan

	// Initialize Table for Disk Chart
	page.DiskChart.Title = " Disk "
	page.DiskChart.TextStyle = ui.NewStyle(ui.ColorWhite)
	page.DiskChart.TextAlignment = ui.AlignLeft
	page.DiskChart.RowSeparator = false
	page.DiskChart.ColumnWidths = []int{9, 9, 9, 9, 9, 11}
	page.DiskChart.BorderStyle.Fg = ui.ColorCyan

	// Initialize Plot for Network Chart
	page.NetworkChart.Title = " Network data(in mB) "
	page.NetworkChart.HorizontalScale = 1
	page.NetworkChart.AxesColor = ui.ColorCyan
	page.NetworkChart.LineColors[0] = ui.ColorRed
	page.NetworkChart.LineColors[1] = ui.ColorGreen
	page.NetworkChart.DrawDirection = widgets.DrawLeft
	page.NetworkChart.BorderStyle.Fg = ui.ColorCyan
	page.NetworkChart.DataLabels = []string{"ip kB", "op kB"} //refer issue #214 for details

	// Initialize paragraph for NetPara
	page.NetPara.Text = "[Received(kB)](fg:cyan)\n\n[Sent(kB)](fg:red)"
	page.NetPara.Border = true
	page.NetPara.BorderStyle.Fg = ui.ColorCyan
	page.NetPara.Title = " RX/TX "

	// Initialize Gauges for each CPU Core usage
	for i := 0; i < numCores; i++ {
		tempGauge := widgets.NewGauge()
		tempGauge.Title = " CPU " + strconv.Itoa(i) + " "
		tempGauge.Percent = 0
		tempGauge.BarColor = ui.ColorBlue
		tempGauge.BorderStyle.Fg = ui.ColorCyan
		tempGauge.TitleStyle.Fg = ui.ColorWhite
		page.CPUCharts = append(page.CPUCharts, tempGauge)
	}

	// Initialize Grid layout
	page.Grid.Set(
		ui.NewRow(0.34, page.MemoryChart),
		ui.NewRow(0.34,
			ui.NewCol(0.25, page.NetPara),
			ui.NewCol(0.75, page.NetworkChart),
		),
		ui.NewRow(0.34, page.DiskChart),
	)

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w/2, h)
}

// RenderCharts handles plotting graphs and charts for system stats in general.
func RenderCharts(endChannel chan os.Signal, memChannel chan []float64, cpuChannel chan []float64, diskChannel chan [][]string, netChannel chan map[string][]float64, wg *sync.WaitGroup) {

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
	myPage := newPage(numCores)

	// Initialize slices for Network Data
	ipData := make([]float64, 40)
	opData := make([]float64, 40)

	// Pause to pause updating data
	pause := func() {
		run = !run
	}

	uiEvents := ui.PollEvents()
	tick := time.Tick(100 * time.Millisecond)
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
						titles[i] = fmt.Sprintf("[Total RX](fg:red): %5.1f %s\n\n", totalBytesRecv/1024, "mB")
					} else {
						titles[i] = fmt.Sprintf("[Total TX](fg:green): %5.1f %s", totalBytesSent/1024, "mB")
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
			myPage.Grid.SetRect(0, 0, w/2, h)

			height := int(h / (numCores - 1))

			// heightOffset := h - (height*numCores - 1) // There's some extra empty space left

			if isCPUSet {
				for i := 0; i < numCores; i++ {
					myPage.CPUCharts[i].SetRect(w/2, i*height, w, (i+1)*height)
					ui.Render(myPage.CPUCharts[i])
				}
			}

			ui.Render(myPage.Grid)

		}
	}
}
