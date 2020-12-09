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

package process

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pesos/grofer/src/utils"
)

// PerProcPage holds the ui elements rendered by the command grofer proc -p PID
type PerProcPage struct {
	Grid             *ui.Grid
	CPUChart         *widgets.Gauge
	MemChart         *widgets.Gauge
	PIDTable         *widgets.Table
	ChildProcsList   *widgets.List
	CTXSwitchesChart *utils.BarChart
	PageFaultsChart  *utils.BarChart
	MemStatsChart    *utils.BarChart
}

// NewProcPage initializes a new page from the PerProcPage struct and returns it
func NewPerProcPage() *PerProcPage {
	page := &PerProcPage{
		Grid:             ui.NewGrid(),
		CPUChart:         widgets.NewGauge(),
		MemChart:         widgets.NewGauge(),
		PIDTable:         widgets.NewTable(),
		ChildProcsList:   widgets.NewList(),
		CTXSwitchesChart: utils.NewBarChart(),
		PageFaultsChart:  utils.NewBarChart(),
		MemStatsChart:    utils.NewBarChart(),
	}
	page.InitPerProc()
	return page
}

// InitPerProc initializes and sets the ui and grid for grofer proc -p PID
func (page *PerProcPage) InitPerProc() {
	// Initialize Gauge for CPU Chart
	page.CPUChart.Title = " CPU % "
	page.CPUChart.LabelStyle.Fg = ui.ColorClear
	page.CPUChart.BarColor = ui.ColorGreen
	page.CPUChart.BorderStyle.Fg = ui.ColorCyan
	page.CPUChart.TitleStyle.Fg = ui.ColorClear

	// Initialize Gauge for Memory Chart
	page.MemChart.Title = " Mem % "
	page.MemChart.LabelStyle.Fg = ui.ColorClear
	page.MemChart.BarColor = ui.ColorGreen
	page.MemChart.BorderStyle.Fg = ui.ColorCyan
	page.MemChart.TitleStyle.Fg = ui.ColorClear

	// Initialize Table for PID Details Table
	page.PIDTable.TextStyle = ui.NewStyle(ui.ColorClear)
	page.PIDTable.TextAlignment = ui.AlignCenter
	page.PIDTable.RowSeparator = false
	page.PIDTable.Title = " PID "
	page.PIDTable.BorderStyle.Fg = ui.ColorCyan
	page.PIDTable.TitleStyle.Fg = ui.ColorClear

	// Initialize List for Child Processes list
	page.ChildProcsList.Title = " Child Processes "
	page.ChildProcsList.BorderStyle.Fg = ui.ColorCyan
	page.ChildProcsList.TitleStyle.Fg = ui.ColorClear
	page.ChildProcsList.TextStyle.Fg = ui.ColorClear

	// Initialize Bar Chart for CTX Swicthes Chart
	page.CTXSwitchesChart.Data = []float64{0, 0}
	page.CTXSwitchesChart.Labels = []string{"Volun", "Involun"}
	page.CTXSwitchesChart.Title = " Ctx switches "
	page.CTXSwitchesChart.BorderStyle.Fg = ui.ColorCyan
	page.CTXSwitchesChart.TitleStyle.Fg = ui.ColorClear
	page.CTXSwitchesChart.BarWidth = 9
	page.CTXSwitchesChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	page.CTXSwitchesChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorClear)}
	page.CTXSwitchesChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Bar Chart for Page Faults Chart
	page.PageFaultsChart.Data = []float64{0, 0}
	page.PageFaultsChart.Labels = []string{"minr", "mjr"}
	page.PageFaultsChart.Title = " Page Faults "
	page.PageFaultsChart.BorderStyle.Fg = ui.ColorCyan
	page.PageFaultsChart.TitleStyle.Fg = ui.ColorClear
	page.PageFaultsChart.BarWidth = 9
	page.PageFaultsChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	page.PageFaultsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorClear)}
	page.PageFaultsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Bar Chart for Memory Stats Chart
	page.MemStatsChart.Data = []float64{0, 0, 0, 0}
	page.MemStatsChart.Labels = []string{"RSS", "Data", "Stack", "Swap"}
	page.MemStatsChart.Title = " Mem Stats (mb) "
	page.MemStatsChart.BorderStyle.Fg = ui.ColorCyan
	page.MemStatsChart.TitleStyle.Fg = ui.ColorClear
	page.MemStatsChart.BarWidth = 9
	page.MemStatsChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorMagenta, ui.ColorYellow, ui.ColorCyan}
	page.MemStatsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorClear)}
	page.MemStatsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Grid layout
	page.Grid.Set(
		ui.NewCol(0.5,
			ui.NewRow(0.125, page.CPUChart),
			ui.NewRow(0.125, page.MemChart),
			ui.NewRow(0.35, page.PIDTable),
			ui.NewRow(0.4, page.ChildProcsList),
		),
		ui.NewCol(0.5,
			ui.NewRow(0.6,
				ui.NewCol(0.5, page.CTXSwitchesChart),
				ui.NewCol(0.5, page.PageFaultsChart),
			),
			ui.NewRow(0.4, page.MemStatsChart),
		),
	)

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)
}

// AllProcPage struct holds the ui elements rendered by the grofer proc command
type AllProcPage struct {
	Grid         *ui.Grid
	HeadingTable *widgets.Table
	BodyList     *widgets.List
}

// NewAllProcsPage initializes a new page from the AllProcPage struct and returns it
func NewAllProcsPage() *AllProcPage {
	page := &AllProcPage{
		Grid:         ui.NewGrid(),
		HeadingTable: widgets.NewTable(),
		BodyList:     widgets.NewList(),
	}
	page.InitAllProc()
	return page
}

// InitAllProc initializes and sets the ui and grid for grofer proc
func (page *AllProcPage) InitAllProc() {
	page.HeadingTable.TextStyle = ui.NewStyle(ui.ColorClear)
	page.HeadingTable.Rows = [][]string{[]string{" PID",
		" Command",
		" CPU",
		" Memory",
		" Status",
		" Foreground",
		" Creation Time",
		" Thread Count",
	}}
	page.HeadingTable.ColumnWidths = []int{10, 40, 10, 10, 8, 12, 23, 15}
	page.HeadingTable.TextAlignment = ui.AlignLeft
	page.HeadingTable.RowSeparator = false

	page.BodyList.TextStyle = ui.NewStyle(ui.ColorClear)
	page.BodyList.TitleStyle.Fg = ui.ColorCyan

	page.Grid.Set(
		ui.NewRow(0.12, page.HeadingTable),
		ui.NewRow(0.88, page.BodyList),
	)

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)
}
