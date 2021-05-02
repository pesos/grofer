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
	DetailsTable *utils.Table
}

// NewOverallContainerPage initializes a new page from the OverallContainerPage struct and returns it
func NewOverallContainerPage() *OverallContainerPage {
	page := &OverallContainerPage{
		Grid:         ui.NewGrid(),
		CPUChart:     widgets.NewGauge(),
		MemChart:     widgets.NewGauge(),
		NetChart:     utils.NewBarChart(),
		BlkChart:     utils.NewBarChart(),
		DetailsTable: utils.NewTable(),
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
	page.DetailsTable.Title = " Details "
	page.DetailsTable.BorderStyle.Fg = ui.ColorCyan
	page.DetailsTable.TitleStyle.Fg = ui.ColorClear
	page.DetailsTable.ColResizer = func() {
		x := page.DetailsTable.Inner.Dx() - (12 + 10 + 10 + 17 + 23)
		page.DetailsTable.ColWidths = []int{
			12,
			ui.MaxInt(15, int(x*3/13)),
			ui.MaxInt(20, int(x*4/13)),
			ui.MaxInt(20, int(x*4/13)),
			ui.MaxInt(10, int(x*2/13)),
			10, 10, 17, 23,
		}
	}
	page.DetailsTable.Header = []string{"ID", "Image", "Name", "Status", "State", "CPU", "Memory", "Net I/O", "Block I/O "}
	page.DetailsTable.ShowCursor = true
	page.DetailsTable.CursorColor = ui.ColorCyan

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
		ui.NewRow(0.6, page.DetailsTable),
	)

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)
}

type PerContainerPage struct {
	Grid          *ui.Grid
	DetailsTable  *utils.Table
	CPUChart      *widgets.Gauge
	MemChart      *widgets.Gauge
	NetChart      *utils.BarChart
	BlkChart      *utils.BarChart
	MountTable    *utils.Table
	NetworkTable  *utils.Table
	CPUUsageTable *utils.Table
	PortMapTable  *utils.Table
	ProcTable     *utils.Table
}

// NewPerContainerPage initializes a new page from the PerContainerPage struct and returns it
func NewPerContainerPage() *PerContainerPage {
	page := &PerContainerPage{
		Grid:          ui.NewGrid(),
		DetailsTable:  utils.NewTable(),
		CPUChart:      widgets.NewGauge(),
		MemChart:      widgets.NewGauge(),
		NetChart:      utils.NewBarChart(),
		BlkChart:      utils.NewBarChart(),
		MountTable:    utils.NewTable(),
		NetworkTable:  utils.NewTable(),
		CPUUsageTable: utils.NewTable(),
		PortMapTable:  utils.NewTable(),
		ProcTable:     utils.NewTable(),
	}
	page.InitPerContainer()
	return page
}

// InitPerContainer initializes and sets the ui and grid for grofer proc -p PID
func (page *PerContainerPage) InitPerContainer() {
	// Initialize Gauge for CPU Chart
	page.CPUChart.Title = " CPU % "
	page.CPUChart.BarColor = ui.ColorGreen
	page.CPUChart.LabelStyle.Fg = ui.ColorClear
	page.CPUChart.BorderStyle.Fg = ui.ColorCyan
	page.CPUChart.TitleStyle.Fg = ui.ColorClear

	// Initialize Gauge for Memory Chart
	page.MemChart.Title = " Mem % "
	page.MemChart.LabelStyle.Fg = ui.ColorClear
	page.MemChart.BarColor = ui.ColorGreen
	page.MemChart.BorderStyle.Fg = ui.ColorCyan
	page.MemChart.TitleStyle.Fg = ui.ColorClear

	// Intialise Bar Chart for Net Chart
	page.NetChart.Data = []float64{0, 0}
	page.NetChart.Title = " Network I/O "
	page.NetChart.Labels = []string{"RX", "TX"}
	page.NetChart.BorderStyle.Fg = ui.ColorCyan
	page.NetChart.TitleStyle.Fg = ui.ColorClear
	page.NetChart.BarWidth = 9
	page.NetChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	page.NetChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorClear)}
	page.NetChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Intialise Bar Chart for Blk Chart
	page.BlkChart.Data = []float64{0, 0}
	page.BlkChart.Title = " Blk I/O "
	page.BlkChart.Labels = []string{"Read", "Write"}
	page.BlkChart.BorderStyle.Fg = ui.ColorCyan
	page.BlkChart.TitleStyle.Fg = ui.ColorClear
	page.BlkChart.BarWidth = 9
	page.BlkChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	page.BlkChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorClear)}
	page.BlkChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Table for Container Details Table
	page.DetailsTable.Title = " Details "
	page.DetailsTable.BorderStyle.Fg = ui.ColorCyan
	page.DetailsTable.TitleStyle.Fg = ui.ColorClear
	page.DetailsTable.ColResizer = func() {
		x := page.DetailsTable.Inner.Dx()
		page.DetailsTable.ColWidths = []int{
			x / 2,
			x / 2,
		}
	}

	// Initialize Table for Mount Table
	page.MountTable.Title = " Mounts "
	page.MountTable.BorderStyle.Fg = ui.ColorCyan
	page.MountTable.TitleStyle.Fg = ui.ColorClear
	page.MountTable.ColGap = 1
	page.MountTable.ColResizer = func() {
		x := page.MountTable.Inner.Dx()
		page.MountTable.ColWidths = []int{
			4 * x / 10,
			4 * x / 10,
			2 * x / 10,
		}
	}
	page.MountTable.Header = []string{"SRC", "DST", "Mode"}
	page.MountTable.CursorColor = ui.ColorCyan

	// Initialize Table for Network Table
	page.NetworkTable.Title = " Networks "
	page.NetworkTable.BorderStyle.Fg = ui.ColorCyan
	page.NetworkTable.TitleStyle.Fg = ui.ColorClear
	page.NetworkTable.ColGap = 1
	page.NetworkTable.ColResizer = func() {
		x := page.NetworkTable.Inner.Dx()
		page.NetworkTable.ColWidths = []int{
			2 * x / 10,
			2 * x / 10,
			3 * x / 10,
			x / 10,
		}
	}
	page.NetworkTable.Header = []string{"Name", "Driver", "IP", "Ingress"}
	page.NetworkTable.CursorColor = ui.ColorCyan

	// Initialize Table for CPU Usage Table
	page.CPUUsageTable.Title = " Per CPU "
	page.CPUUsageTable.BorderStyle.Fg = ui.ColorCyan
	page.CPUUsageTable.TitleStyle.Fg = ui.ColorClear
	page.CPUUsageTable.ColResizer = func() {
		x := page.CPUUsageTable.Inner.Dx()
		page.CPUUsageTable.ColWidths = []int{
			x / 2,
			x / 2,
		}
	}
	page.CPUUsageTable.Header = []string{"CPU", "Usage"}
	page.CPUUsageTable.CursorColor = ui.ColorCyan

	// Initialize Table for Port Map Table
	page.PortMapTable.Title = " Port Mappings "
	page.PortMapTable.BorderStyle.Fg = ui.ColorCyan
	page.PortMapTable.TitleStyle.Fg = ui.ColorClear
	page.PortMapTable.ColGap = 1
	page.PortMapTable.ColResizer = func() {
		x := page.PortMapTable.Inner.Dx()
		page.PortMapTable.ColWidths = []int{
			3 * x / 10,
			3 * x / 10,
			3 * x / 10,
		}
	}
	page.PortMapTable.Header = []string{"Host", "Container", "Type"}
	page.PortMapTable.CursorColor = ui.ColorCyan

	// Initialize Table for procs Table
	page.ProcTable.Title = " Processes "
	page.ProcTable.BorderStyle.Fg = ui.ColorCyan
	page.ProcTable.TitleStyle.Fg = ui.ColorClear
	page.ProcTable.ColResizer = func() {
		x := page.ProcTable.Inner.Dx()
		page.ProcTable.ColWidths = []int{
			2 * x / 10,
			3 * x / 10,
			5 * x / 10,
		}
	}
	page.ProcTable.Header = []string{"PID", "UID", "CMD"}
	page.ProcTable.CursorColor = ui.ColorCyan

	// Initialize Grid layout
	page.Grid.Set(
		ui.NewRow(0.3,
			ui.NewCol(0.3, page.DetailsTable),
			ui.NewCol(0.45, page.MountTable),
			ui.NewCol(0.25, page.BlkChart),
		),
		ui.NewRow(0.3,
			ui.NewCol(0.3,
				ui.NewRow(0.5, page.MemChart),
				ui.NewRow(0.5, page.CPUChart),
			),
			ui.NewCol(0.45, page.NetworkTable),
			ui.NewCol(0.25, page.NetChart),
		),
		ui.NewRow(0.4,
			ui.NewCol(0.2, page.CPUUsageTable),
			ui.NewCol(0.4, page.PortMapTable),
			ui.NewCol(0.4, page.ProcTable),
		),
	)

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)
}
