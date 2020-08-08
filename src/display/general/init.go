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
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// MainPage contains the ui widgets for the ui rendered by the grofer command
type MainPage struct {
	Grid         *ui.Grid
	MemoryChart  *widgets.BarChart
	DiskChart    *widgets.Table
	NetworkChart *widgets.Plot
	CPUCharts    []*widgets.Gauge
	NetPara      *widgets.Paragraph
}

// NewPage returns a new page initialized from the MainPage struct
func NewPage(numCores int) *MainPage {
	page := &MainPage{
		Grid:         ui.NewGrid(),
		MemoryChart:  widgets.NewBarChart(),
		DiskChart:    widgets.NewTable(),
		NetworkChart: widgets.NewPlot(),
		CPUCharts:    make([]*widgets.Gauge, 0),
		NetPara:      widgets.NewParagraph(),
	}
	page.InitGeneral(numCores)
	return page
}

// InitGeneral initializes all ui elements for the ui rendered by the grofer command
func (page *MainPage) InitGeneral(numCores int) {

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
	page.NetPara.Text = "[Total RX](fg:red): 0\n\n[Total TX](fg:green): 0"
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

	// Get Terminal Dimensions
	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(w/2, 0, w, h)
}
