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

package container

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pesos/grofer/src/utils"
)

type OverallContainerPage struct {
	Grid         *ui.Grid
	CPUChart     *widgets.Gauge
	MemChart     *widgets.Gauge
	NetChart     *utils.BarChart
	BlkChart     *utils.BarChart
	HeadingTable *widgets.Table
	BodyList     *widgets.List
}

// NewOverallContainerPage initializes a new page from the OverallContainerPage struct and returns it
func NewOverallContainerPage() *OverallContainerPage {
	page := &OverallContainerPage{
		Grid:         ui.NewGrid(),
		CPUChart:     widgets.NewGauge(),
		MemChart:     widgets.NewGauge(),
		NetChart:     utils.NewBarChart(),
		BlkChart:     utils.NewBarChart(),
		HeadingTable: widgets.NewTable(),
		BodyList:     widgets.NewList(),
	}
	page.InitOverallContainer()
	return page
}

// InitOverallContainer initializes and sets the ui and grid for grofer proc -p PID
func (page *OverallContainerPage) InitOverallContainer() {
	// Initialize Gauge for CPU Chart
	page.CPUChart.Title = " Total CPU % "
	page.CPUChart.LabelStyle.Fg = ui.ColorClear
	page.CPUChart.BarColor = ui.ColorGreen
	page.CPUChart.BorderStyle.Fg = ui.ColorCyan
	page.CPUChart.TitleStyle.Fg = ui.ColorClear

	// Initialize Gauge for Memory Chart
	page.MemChart.Title = " Total Mem % "
	page.MemChart.LabelStyle.Fg = ui.ColorClear
	page.MemChart.BarColor = ui.ColorGreen
	page.MemChart.BorderStyle.Fg = ui.ColorCyan
	page.MemChart.TitleStyle.Fg = ui.ColorClear

	// Intialise Bar Chart for Net Chart
	page.NetChart.Data = []float64{0, 0}
	page.NetChart.Title = " Total Network I/O "
	page.NetChart.Labels = []string{"RX", "TX"}
	page.NetChart.BorderStyle.Fg = ui.ColorCyan
	page.NetChart.TitleStyle.Fg = ui.ColorClear
	page.NetChart.BarWidth = 9
	page.NetChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	page.NetChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorClear)}
	page.NetChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Intialise Bar Chart for Blk Chart
	page.BlkChart.Data = []float64{0, 0}
	page.BlkChart.Title = " Total Blk I/O "
	page.BlkChart.Labels = []string{"Read", "Write"}
	page.BlkChart.BorderStyle.Fg = ui.ColorCyan
	page.BlkChart.TitleStyle.Fg = ui.ColorClear
	page.BlkChart.BarWidth = 9
	page.BlkChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	page.BlkChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorClear)}
	page.BlkChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Table for Container Details Table
	page.HeadingTable.TextStyle = ui.NewStyle(ui.ColorClear)
	page.HeadingTable.Rows = [][]string{{
		" ID",
		" Image",
		" Name",
		" Status",
		" State",
		" CPU",
		" Memory",
		" Net I/O",
		" Block I/O ",
	}}
	page.HeadingTable.ColumnWidths = []int{15, 16, 20, 15, 15, 10, 10, 17, 17}
	page.HeadingTable.TextAlignment = ui.AlignLeft
	page.HeadingTable.RowSeparator = false

	// Initialize List for Conatiner list
	page.BodyList.TextStyle = ui.NewStyle(ui.ColorClear)
	page.BodyList.TextStyle.Fg = ui.ColorClear
	page.BodyList.TitleStyle.Fg = ui.ColorClear
	page.BodyList.BorderStyle.Fg = ui.ColorCyan

	// Initialize Grid layout
	page.Grid.Set(
		ui.NewRow(0.4,
			ui.NewCol(0.5,
				ui.NewRow(0.5, page.CPUChart),
				ui.NewRow(0.5, page.MemChart),
			),
			ui.NewCol(0.25, page.NetChart),
			ui.NewCol(0.25, page.BlkChart),
		),
		ui.NewRow(0.6,
			ui.NewRow(0.2, page.HeadingTable),
			ui.NewRow(0.8, page.BodyList),
		),
	)

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)
}
