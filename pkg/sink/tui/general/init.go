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
	viz "github.com/pesos/grofer/pkg/utils/visualization"
)

// MainPage contains the UI widgets for the UI rendered by the grofer command
type MainPage struct {
	Grid             *ui.Grid
	MemoryChart      *viz.HorizontalBarChart
	DiskChart        *viz.Table
	NetworkChart     *viz.SparklineGroup
	CPUTable         *viz.CPUTableChart
	CPUChart         *viz.LineGraph
	TemperatureTable *viz.Table
	InfoTable        *viz.Table
	BatteryGauge     *widgets.Gauge
	selectedTable    int
	cpuTableVisible  bool
}

// CPUPage contains the UI widgets for the UI rendered by the grofer -c command
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
	CPUTable    *viz.Table
}

// NewPage returns a new page initialized from the MainPage struct
func NewPage(numCores int) *MainPage {
	rxSparkLine := viz.NewSparkline()
	rxSparkLine.Data = []float64{}
	txSparkLine := viz.NewSparkline()
	txSparkLine.Data = []float64{}

	page := &MainPage{
		Grid:             ui.NewGrid(),
		MemoryChart:      viz.NewHorizontalBarChart(),
		DiskChart:        viz.NewTable(),
		NetworkChart:     viz.NewSparklineGroup(rxSparkLine, txSparkLine),
		CPUTable:         viz.NewCPUTableChart(),
		CPUChart:         viz.NewLineGraph(),
		TemperatureTable: viz.NewTable(),
		InfoTable:        viz.NewTable(),
		BatteryGauge:     widgets.NewGauge(),
	}
	page.init(numCores)
	return page
}

// NewCPUPage is a constructor for the CPUPage type
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
		CPUTable:    viz.NewTable(),
	}
	page.init(numCores)
	return page
}

// init initializes all ui elements for the ui rendered by the grofer command
func (page *MainPage) init(numCores int) {
	if numCores > 8 {
		page.cpuTableVisible = true
	}

	// Initialize Bar Graph for Memory Chart
	page.initMemoryChartWidget()
	// Initialize Table for Disk Chart
	page.initDiskChartWidget()
	// Initialize Plot for Network Chart
	page.initNetworkChartWidget()
	// Initialize Table for Info Table
	page.initInfoTableWidget()
	// Initialise Battery Gauge
	page.initBatteryGauge()

	if page.cpuTableVisible {
		page.initCPUTableWidget(numCores)
	} else {
		page.initCPUChartWidget(numCores)
	}
	// Initialize Graph for Temperature Table
	page.initTemperatureTableWidget()
	// Set page grid
	page.initPageGrid()

}

// ToggleCPUWidget helps toggle the widget on grid used to display CPU usage
func (page *MainPage) ToggleCPUWidget() viz.ScrollableWidget {
	defer page.initPageGrid()
	if page.cpuTableVisible {
		page.cpuTableVisible = false
		return page.DiskChart
	}
	page.cpuTableVisible = true
	return page.CPUTable
}

func (page *MainPage) initPageGrid() {
	page.Grid = ui.NewGrid()
	if page.cpuTableVisible {
		page.Grid.Set(
			ui.NewCol(
				0.4,
				ui.NewRow(0.1, page.BatteryGauge),
				ui.NewRow(0.2, page.InfoTable),
				ui.NewRow(0.3, page.TemperatureTable),
				ui.NewRow(0.4, page.MemoryChart),
			),
			ui.NewCol(
				0.6,
				ui.NewRow(0.4, page.CPUTable),
				ui.NewRow(0.3, page.NetworkChart),
				ui.NewRow(0.3, page.DiskChart),
			),
		)
	} else {
		page.Grid.Set(
			ui.NewCol(
				0.4,
				ui.NewRow(0.1, page.BatteryGauge),
				ui.NewRow(0.2, page.InfoTable),
				ui.NewRow(0.3, page.TemperatureTable),
				ui.NewRow(0.4, page.MemoryChart),
			),
			ui.NewCol(
				0.6,
				ui.NewRow(0.4, page.CPUChart),
				ui.NewRow(0.3, page.NetworkChart),
				ui.NewRow(0.3, page.DiskChart),
			),
		)
	}
}

func (page *MainPage) initBatteryGauge() {
	page.BatteryGauge.Title = " Battery Not Found "
	page.BatteryGauge.BorderStyle.Fg = ui.ColorCyan
	page.BatteryGauge.BarColor = ui.ColorGreen
	page.BatteryGauge.Percent = 0
}

func (page *MainPage) initInfoTableWidget() {
	page.InfoTable.Title = " System Info "
	page.InfoTable.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.InfoTable.BorderStyle.Fg = ui.ColorCyan
	page.InfoTable.HeaderStyle = ui.NewStyle(ui.ColorClear, ui.ColorClear, ui.ModifierBold)
	page.InfoTable.ShowCursor = false
	page.InfoTable.ColResizer = func() {
		x := page.InfoTable.Inner.Dx()
		page.InfoTable.ColWidths = []int{x / 2, x / 2}
	}
}

func (page *MainPage) initMemoryChartWidget() {
	page.MemoryChart.Title = " Memory(RAM) "
	page.MemoryChart.BorderStyle.Fg = ui.ColorCyan
	page.MemoryChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.MemoryChart.BarGap = 1
	page.MemoryChart.ColResizer = func() {
		if page.MemoryChart.Inner.Dy() > 12 {
			page.MemoryChart.BarWidth = 2
		} else {
			page.MemoryChart.BarWidth = 1
		}
	}
}

func (page *MainPage) initTemperatureTableWidget() {
	page.TemperatureTable.Title = " Temp "
	page.TemperatureTable.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.TemperatureTable.BorderStyle.Fg = ui.ColorCyan
	page.TemperatureTable.HeaderStyle = ui.NewStyle(ui.ColorClear, ui.ColorClear, ui.ModifierBold)
	page.TemperatureTable.ColColor[1] = ui.ColorGreen
	page.TemperatureTable.ShowCursor = false
	page.TemperatureTable.ColResizer = func() {
		x := page.TemperatureTable.Inner.Dx()
		page.TemperatureTable.ColWidths = []int{2 * x / 3, x / 3}
	}
}

func (page *MainPage) initDiskChartWidget() {
	page.DiskChart.Title = " Disk "
	page.DiskChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.DiskChart.HeaderStyle = ui.NewStyle(ui.ColorClear, ui.ColorClear, ui.ModifierBold)
	page.DiskChart.ShowCursor = true
	page.DiskChart.BorderStyle.Fg = ui.ColorClear
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

func (page *MainPage) initNetworkChartWidget() {
	page.NetworkChart.Title = " Network data "
	page.NetworkChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.NetworkChart.BorderStyle.Fg = ui.ColorCyan
	page.NetworkChart.Sparklines[0].TitleStyle.Fg = ui.ColorRed
	page.NetworkChart.Sparklines[0].LineColor = ui.ColorRed
	page.NetworkChart.Sparklines[1].TitleStyle.Fg = ui.ColorGreen
	page.NetworkChart.Sparklines[1].LineColor = ui.ColorGreen
	page.NetworkChart.Sparklines[1].Reverse = true
}

func (page *MainPage) initCPUChartWidget(numCores int) {
	page.CPUChart.Title = " CPU Usage "
	page.CPUChart.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.CPUChart.BorderStyle.Fg = ui.ColorCyan
	page.CPUChart.DefaultLineColor = ui.ColorClear
	page.CPUChart.MaxVal = 100
	for i := 0; i < numCores; i++ {
		cpu := fmt.Sprintf("CPU %d", i)
		page.CPUChart.Data[cpu] = []float64{0}
		page.CPUChart.LineColors[cpu] = ui.SelectColor(ui.StandardColors, i)
	}
}

func (page *MainPage) initCPUTableWidget(numCores int) {
	page.CPUTable.Title = " Per CPU Usage "
	page.CPUTable.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.CPUTable.BorderStyle.Fg = ui.ColorCyan
	page.CPUTable.NumCores = numCores
}

func (page *MainPage) switchTableLeft(cpuTableVisible bool) viz.ScrollableWidget {
	if cpuTableVisible {
		scrollableWidgets := []viz.ScrollableWidget{page.InfoTable, page.CPUTable, page.DiskChart, page.TemperatureTable}
		page.selectedTable = (page.selectedTable + 1) % len(scrollableWidgets)
		return scrollableWidgets[page.selectedTable]
	}

	scrollableWidgets := []viz.ScrollableWidget{page.InfoTable, page.DiskChart, page.TemperatureTable}
	page.selectedTable = (page.selectedTable + 1) % len(scrollableWidgets)
	return scrollableWidgets[page.selectedTable]
}

func (page *MainPage) switchTableRight(cpuTableVisible bool) viz.ScrollableWidget {
	if cpuTableVisible {
		scrollableWidgets := []viz.ScrollableWidget{page.InfoTable, page.TemperatureTable, page.DiskChart, page.CPUTable}
		page.selectedTable = (page.selectedTable - 1)
		if page.selectedTable < 0 {
			page.selectedTable = len(scrollableWidgets) - 1
		}
		return scrollableWidgets[page.selectedTable]
	}

	scrollableWidgets := []viz.ScrollableWidget{page.InfoTable, page.DiskChart, page.TemperatureTable}
	page.selectedTable = (page.selectedTable - 1)
	if page.selectedTable < 0 {
		page.selectedTable = len(scrollableWidgets) - 1
	}
	return scrollableWidgets[page.selectedTable]
}

func (page *CPUPage) init(numCores int) {
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
