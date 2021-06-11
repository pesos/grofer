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

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pesos/grofer/src/utils"
)

// MainPage contains the ui widgets for the ui rendered by the grofer command
type MainPage struct {
	Grid             *ui.Grid
	MemoryChart      *utils.HorizontalBarChart
	DiskChart        *utils.Table
	NetworkChart     *utils.SparklineGroup
	CPUGraph         *utils.LineGraph
	CPUTable         *utils.CpuTableChart
	AvgCPUGraph      *utils.LineGraph
	TemperatureTable *utils.Table
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
	CPUTable    *utils.Table
}

// NewPage returns a new page initialized from the MainPage struct
func NewPage(numCores int) *MainPage {
	rxSparkLine := utils.NewSparkline()
	rxSparkLine.Data = []float64{}
	txSparkLine := utils.NewSparkline()
	txSparkLine.Data = []float64{}

	memAvailableSparkLine := widgets.NewSparkline()
	memAvailableSparkLine.Data = []float64{}
	memUsedSparkLine := widgets.NewSparkline()
	memUsedSparkLine.Data = []float64{}
	memFreeSparkLine := widgets.NewSparkline()
	memFreeSparkLine.Data = []float64{}
	memCachedSparkLine := widgets.NewSparkline()
	memCachedSparkLine.Data = []float64{}

	page := &MainPage{
		Grid:             ui.NewGrid(),
		MemoryChart:      utils.NewHorizontalBarChart(),
		DiskChart:        utils.NewTable(),
		NetworkChart:     utils.NewSparklineGroup(rxSparkLine, txSparkLine),
		CPUGraph:         utils.NewLineGraph(),
		CPUTable:         utils.NewCpuTableChart(),
		AvgCPUGraph:      utils.NewLineGraph(),
		TemperatureTable: utils.NewTable(),
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
		CPUTable:    utils.NewTable(),
	}
	page.InitCPU(numCores)
	return page
}

// InitGeneral initializes all ui elements for the ui rendered by the grofer command
func (page *MainPage) InitGeneral(numCores int) {
	// Initialize Bar Graph for Memory Chart
	page.memoryChartWidget()
	// Initialize Table for Disk Chart
	page.diskChartWidget()
	// Initialize Plot for Network Chart
	page.networkChartWidget()
	// Initialize Graph for CPU Usage
	if numCores > 8 {
		page.avgCpuGraphWidget()
		page.cpuTableWidget(numCores)
	} else {
		page.cpuGraphWidget(numCores)
	}
	// Initialize Graph for Temperature Table
	page.temperatureTableWidget()
	// Get Terminal Dimensions
	w, h := ui.TerminalDimensions()
	if numCores > 8 {
		page.Grid.Set(
			ui.NewCol(0.4,
				ui.NewRow(0.5, page.AvgCPUGraph),
				ui.NewRow(0.5, page.CPUTable)),
			ui.NewCol(0.6,
				ui.NewRow(0.34, page.MemoryChart),
				ui.NewRow(0.34, ui.NewCol(0.5, page.NetworkChart), ui.NewCol(0.5, page.TemperatureTable)),
				ui.NewRow(0.34, page.DiskChart)),
		)
	} else {
		page.Grid.Set(
			ui.NewCol(0.4, page.CPUGraph),
			ui.NewCol(0.6,
				ui.NewRow(0.34, page.MemoryChart),
				ui.NewRow(0.34, ui.NewCol(0.5, page.NetworkChart), ui.NewCol(0.5, page.TemperatureTable)),
				ui.NewRow(0.34, page.DiskChart)),
		)
	}
	page.Grid.SetRect(0, 0, w, h)
}

func (page *MainPage) memoryChartWidget() {
	page.MemoryChart.Title = " Memory(RAM) "
	page.MemoryChart.BorderStyle.Fg = ui.ColorCyan
	page.MemoryChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.MemoryChart.ColResizer = func() {
		if page.MemoryChart.Inner.Dy() > 8 {
			page.MemoryChart.BarWidth = 2
		} else {
			page.MemoryChart.BarWidth = 1
		}
	}
}

func (page *MainPage) temperatureTableWidget() {
	page.TemperatureTable.Title = " Temp "
	page.TemperatureTable.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.TemperatureTable.BorderStyle.Fg = ui.ColorCyan
	page.TemperatureTable.HeaderStyle = ui.NewStyle(ui.ColorClear, ui.ColorClear, ui.ModifierUnderline)
	page.TemperatureTable.ColColor[1] = ui.ColorGreen
	page.TemperatureTable.ShowCursor = true
	page.TemperatureTable.ColResizer = func() {
		x := page.TemperatureTable.Inner.Dx()
		page.TemperatureTable.ColWidths = []int{x / 2, x / 2}
	}
}

func (page *MainPage) diskChartWidget() {
	page.DiskChart.Title = " Disk "
	page.DiskChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.DiskChart.ShowCursor = true
	page.DiskChart.BorderStyle.Fg = ui.ColorCyan
	page.DiskChart.ColResizer = func() {
		// Middle 4 columns are of fixed length
		x := page.DiskChart.Inner.Dx()
		page.DiskChart.ColWidths = []int{
			x / 6,
			x / 6,
			x / 6,
			x / 6,
			x / 6,
			x / 6,
		}
	}
}

func (page *MainPage) networkChartWidget() {
	page.NetworkChart.Title = " Network data "
	page.NetworkChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.NetworkChart.BorderStyle.Fg = ui.ColorCyan
	page.NetworkChart.Sparklines[0].TitleStyle.Fg = ui.ColorRed
	page.NetworkChart.Sparklines[0].LineColor = ui.ColorRed
	page.NetworkChart.Sparklines[1].TitleStyle.Fg = ui.ColorGreen
	page.NetworkChart.Sparklines[1].LineColor = ui.ColorGreen
	page.NetworkChart.Sparklines[1].Reverse = true
}

func (page *MainPage) cpuGraphWidget(numCores int) {
	page.CPUGraph.Title = " CPU Usage "
	page.CPUGraph.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.CPUGraph.HorizontalScale = 10
	page.CPUGraph.BorderStyle.Fg = ui.ColorCyan
	page.CPUGraph.DefaultLineColor = ui.ColorClear
	page.CPUGraph.Min.Y = 0
	page.CPUGraph.Max.Y = 100
	for i := 0; i < numCores; i++ {
		key := fmt.Sprintf("CPU%d", i)
		page.CPUGraph.LineColors[key] = ui.Color(i + 1)
		page.CPUGraph.Data[key] = []float64{0}
	}
}

func (page *MainPage) cpuTableWidget(numCores int) {
	page.CPUTable.Title = " CPU Table "
	page.CPUTable.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.CPUTable.BorderStyle.Fg = ui.ColorCyan
	page.CPUTable.NumCores = numCores
}

func (page *MainPage) avgCpuGraphWidget() {
	page.AvgCPUGraph.Title = " Average CPU Usage "
	page.AvgCPUGraph.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.AvgCPUGraph.HorizontalScale = 10
	page.AvgCPUGraph.BorderStyle.Fg = ui.ColorCyan
	page.AvgCPUGraph.DefaultLineColor = ui.ColorClear
	page.AvgCPUGraph.Min.Y = 0
	page.AvgCPUGraph.Max.Y = 100
	page.AvgCPUGraph.LineColors["Average CPU Load:"] = ui.ColorClear
	page.AvgCPUGraph.Data["Average CPU Load:"] = []float64{0}
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

	page.CPUChart.Title = " CPU Usage "
	page.CPUChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.CPUChart.BorderStyle = ui.NewStyle(ui.ColorCyan)
	page.CPUChart.TextStyle = ui.NewStyle(ui.ColorClear)
	page.CPUChart.TextAlignment = ui.AlignCenter
	page.CPUChart.RowSeparator = true
	page.CPUChart.ColumnResizer = func() {
		columnWidths := []int{}
		x := page.CPUChart.Inner.Dx()
		for i := 0; i < numCores; i++ {
			columnWidths = append(columnWidths, x/numCores)
		}

		page.CPUChart.ColumnWidths = columnWidths
	}

	page.CPUTable.Title = " CPU Usage "
	page.CPUTable.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.CPUTable.BorderStyle = ui.NewStyle(ui.ColorCyan)
	page.CPUTable.ColResizer = func() {
		x := page.CPUTable.Inner.Dx()

		page.CPUTable.ColWidths = []int{x / 2, x / 2}
	}
	page.CPUTable.Header = []string{"CPU", "Usage"}
	page.CPUTable.ShowCursor = true
	page.CPUTable.CursorColor = ui.ColorCyan

	if numCores > 8 {
		page.Grid.Set(
			ui.NewCol(0.3, page.CPUTable),
			ui.NewCol(0.7,
				ui.NewRow(0.125, page.UsrChart),
				ui.NewRow(0.125, page.NiceChart),
				ui.NewRow(0.125, page.SysChart),
				ui.NewRow(0.125, page.IowaitChart),
				ui.NewRow(0.125, page.IrqChart),
				ui.NewRow(0.125, page.SoftChart),
				ui.NewRow(0.125, page.IdleChart),
				ui.NewRow(0.125, page.StealChart),
			),
		)
	} else {
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
	}

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)

}
