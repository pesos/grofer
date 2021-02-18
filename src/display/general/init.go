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
	NetPara      *widgets.Paragraph
	CPUCharts    []*widgets.Gauge
}

type CPUPage struct {
	Grid        *ui.Grid
	UsrChart    *widgets.Gauge
	NiceChart   *widgets.Gauge
	SysChart    *widgets.Gauge
	IowaitChart *widgets.Gauge
	IrqChart    *widgets.Gauge
	SoftChart   *widgets.Gauge
	IdleChart   *widgets.Gauge
	StealChart  *widgets.Gauge
	CPUChart    *widgets.Table
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

func NewCPUPage(numCores int) *CPUPage {
	page := &CPUPage{
		Grid:        ui.NewGrid(),
		UsrChart:    widgets.NewGauge(),
		NiceChart:   widgets.NewGauge(),
		SysChart:    widgets.NewGauge(),
		IowaitChart: widgets.NewGauge(),
		IrqChart:    widgets.NewGauge(),
		SoftChart:   widgets.NewGauge(),
		IdleChart:   widgets.NewGauge(),
		StealChart:  widgets.NewGauge(),
		CPUChart:    widgets.NewTable(),
	}
	page.InitCPU(numCores)
	return page
}

// InitGeneral initializes all ui elements for the ui rendered by the grofer command
func (page *MainPage) InitGeneral(numCores int) {

	// Initialize Bar Graph for Memory Chart
	page.MemoryChart.Title = " Memory (RAM) "
	page.MemoryChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.MemoryChart.Labels = []string{"Total", "Available", "Used", "Free"}
	page.MemoryChart.BarWidth = 8
	page.MemoryChart.BarGap = 9
	page.MemoryChart.BarColors = []ui.Color{ui.ColorCyan, ui.ColorGreen}
	page.MemoryChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorClear)}
	page.MemoryChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
	page.MemoryChart.BorderStyle.Fg = ui.ColorCyan

	// Initialize Table for Disk Chart
	page.DiskChart.Title = " Disk "
	page.DiskChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.DiskChart.TextStyle = ui.NewStyle(ui.ColorClear)
	page.DiskChart.TextAlignment = ui.AlignLeft
	page.DiskChart.RowSeparator = false
	page.DiskChart.ColumnWidths = []int{9, 9, 9, 9, 9, 11}
	page.DiskChart.BorderStyle.Fg = ui.ColorCyan

	// Initialize Plot for Network Chart
	page.NetworkChart.Title = " Network data(in mB) "
	page.NetworkChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.NetworkChart.HorizontalScale = 1
	page.NetworkChart.AxesColor = ui.ColorCyan
	page.NetworkChart.LineColors[0] = ui.ColorRed
	page.NetworkChart.LineColors[1] = ui.ColorGreen
	page.NetworkChart.DrawDirection = widgets.DrawLeft
	page.NetworkChart.BorderStyle.Fg = ui.ColorCyan
	page.NetworkChart.DataLabels = []string{"ip kB", "op kB"} //refer issue #214 for details

	// Initialize paragraph for NetPara
	page.NetPara.Text = "[Total RX](fg:red): 0\n\n[Total TX](fg:green): 0"
	page.NetPara.TextStyle = ui.NewStyle(ui.ColorClear)
	page.NetPara.Border = true
	page.NetPara.BorderStyle.Fg = ui.ColorCyan
	page.NetPara.Title = " RX/TX "
	page.NetPara.TitleStyle = ui.NewStyle(ui.ColorClear)

	// Initialize Gauges for each CPU Core usage
	for i := 0; i < numCores; i++ {
		tempGauge := widgets.NewGauge()
		tempGauge.Title = " CPU " + strconv.Itoa(i) + " "
		tempGauge.Percent = 0
		tempGauge.BarColor = ui.ColorBlue
		tempGauge.BorderStyle.Fg = ui.ColorCyan
		tempGauge.TitleStyle.Fg = ui.ColorClear
		tempGauge.LabelStyle.Fg = ui.ColorClear
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

func (page *CPUPage) InitCPU(numCores int) {
	page.UsrChart.Title = " Usr "
	page.UsrChart.Percent = 0
	page.UsrChart.BarColor = ui.ColorBlue
	page.UsrChart.BorderStyle.Fg = ui.ColorCyan
	page.UsrChart.TitleStyle.Fg = ui.ColorClear
	page.UsrChart.LabelStyle.Fg = ui.ColorClear

	page.NiceChart.Title = " Nice "
	page.NiceChart.Percent = 0
	page.NiceChart.BarColor = ui.ColorBlue
	page.NiceChart.BorderStyle.Fg = ui.ColorCyan
	page.NiceChart.TitleStyle.Fg = ui.ColorClear
	page.NiceChart.LabelStyle.Fg = ui.ColorClear

	page.SysChart.Title = " Sys "
	page.SysChart.Percent = 0
	page.SysChart.BarColor = ui.ColorBlue
	page.SysChart.BorderStyle.Fg = ui.ColorCyan
	page.SysChart.TitleStyle.Fg = ui.ColorClear
	page.SysChart.LabelStyle.Fg = ui.ColorClear

	page.IowaitChart.Title = " Iowait "
	page.IowaitChart.Percent = 0
	page.IowaitChart.BarColor = ui.ColorBlue
	page.IowaitChart.BorderStyle.Fg = ui.ColorCyan
	page.IowaitChart.TitleStyle.Fg = ui.ColorClear
	page.IowaitChart.LabelStyle.Fg = ui.ColorClear

	page.IrqChart.Title = " Irq "
	page.IrqChart.Percent = 0
	page.IrqChart.BarColor = ui.ColorBlue
	page.IrqChart.BorderStyle.Fg = ui.ColorCyan
	page.IrqChart.TitleStyle.Fg = ui.ColorClear
	page.IrqChart.LabelStyle.Fg = ui.ColorClear

	page.SoftChart.Title = " Soft "
	page.SoftChart.Percent = 0
	page.SoftChart.BarColor = ui.ColorBlue
	page.SoftChart.BorderStyle.Fg = ui.ColorCyan
	page.SoftChart.TitleStyle.Fg = ui.ColorClear
	page.SoftChart.LabelStyle.Fg = ui.ColorClear

	page.IdleChart.Title = " Idle "
	page.IdleChart.Percent = 0
	page.IdleChart.BarColor = ui.ColorBlue
	page.IdleChart.BorderStyle.Fg = ui.ColorCyan
	page.IdleChart.TitleStyle.Fg = ui.ColorClear
	page.IdleChart.LabelStyle.Fg = ui.ColorClear

	page.StealChart.Title = " Steal "
	page.StealChart.Percent = 0
	page.StealChart.BarColor = ui.ColorBlue
	page.StealChart.BorderStyle.Fg = ui.ColorCyan
	page.StealChart.TitleStyle.Fg = ui.ColorClear
	page.StealChart.LabelStyle.Fg = ui.ColorClear

	page.CPUChart.Title = " CPU "
	page.CPUChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.CPUChart.BorderStyle = ui.NewStyle(ui.ColorCyan)
	page.CPUChart.TextStyle = ui.NewStyle(ui.ColorClear)
	page.CPUChart.TextAlignment = ui.AlignCenter
	page.CPUChart.RowSeparator = true

	columnWidths := []int{}
	for i := 0; i < numCores; i++ {
		columnWidths = append(columnWidths, 9)
	}

	page.CPUChart.ColumnWidths = columnWidths

	page.Grid.Set(
		ui.NewRow(0.17,
			ui.NewCol(0.5, page.UsrChart),
			ui.NewCol(0.5, page.NiceChart),
		),
		ui.NewRow(0.17,
			ui.NewCol(0.5, page.SysChart),
			ui.NewCol(0.5, page.IowaitChart),
		),
		ui.NewRow(0.17,
			ui.NewCol(0.5, page.IrqChart),
			ui.NewCol(0.5, page.SoftChart),
		),
		ui.NewRow(0.17,
			ui.NewCol(0.5, page.IdleChart),
			ui.NewCol(0.5, page.StealChart),
		),
		ui.NewRow(0.30, page.CPUChart),
	)

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)

}
