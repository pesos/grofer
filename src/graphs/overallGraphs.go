package graphs

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type gaugeMap map[int]*widgets.Gauge

type mainPage struct {
	Grid         *ui.Grid
	MemoryChart  *widgets.BarChart
	DiskChart    *widgets.Table
	NetworkChart *widgets.Plot
	CPUCharts    []*widgets.Gauge
}

func newPage(numCores int) *mainPage {
	page := &mainPage{
		Grid:         ui.NewGrid(),
		MemoryChart:  widgets.NewBarChart(),
		DiskChart:    widgets.NewTable(),
		NetworkChart: widgets.NewPlot(),
		CPUCharts:    make([]*widgets.Gauge, 0),
	}
	page.init(numCores)
	return page
}

func (page *mainPage) init(numCores int) {
	page.MemoryChart.Title = " Memory (RAM) "
	page.MemoryChart.Labels = []string{"Total", "Available", "Used", "Free"}
	page.MemoryChart.BarWidth = 8
	page.MemoryChart.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	page.MemoryChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	page.MemoryChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	page.DiskChart.Title = " Disk "
	page.DiskChart.TextStyle = ui.NewStyle(ui.ColorWhite)
	page.DiskChart.TextAlignment = ui.AlignLeft
	page.DiskChart.RowSeparator = false
	page.DiskChart.ColumnWidths = []int{9, 9, 9, 9, 9, 11}
	// page.DiskChart.FillRow = true

	page.NetworkChart.Title = " Network data "
	page.NetworkChart.HorizontalScale = 1
	page.NetworkChart.AxesColor = ui.ColorWhite
	page.NetworkChart.LineColors[0] = ui.ColorCyan
	page.NetworkChart.LineColors[1] = ui.ColorRed
	page.NetworkChart.DrawDirection = 1
	page.NetworkChart.DataLabels = []string{"ip kB", "op kB"}

	for i := 0; i < numCores; i++ {
		tempGauge := widgets.NewGauge()
		tempGauge.Title = " CPU " + strconv.Itoa(i) + " "
		// tempGauge.SetRect(0, 0+(i*3), 35, 0+((i+1)*3))
		tempGauge.Percent = 0
		tempGauge.BarColor = ui.ColorRed
		tempGauge.BorderStyle.Fg = ui.ColorWhite
		tempGauge.TitleStyle.Fg = ui.ColorCyan
		page.CPUCharts = append(page.CPUCharts, tempGauge)
	}

	// var rows *ui.GridItem
	// rows.Entry = page.CPUCharts

	if numCores == 8 {
		page.Grid.Set(
			ui.NewCol(0.54,
				ui.NewRow(0.125, page.CPUCharts[0]),
				ui.NewRow(0.125, page.CPUCharts[1]),
				ui.NewRow(0.125, page.CPUCharts[2]),
				ui.NewRow(0.125, page.CPUCharts[3]),
				ui.NewRow(0.125, page.CPUCharts[4]),
				ui.NewRow(0.125, page.CPUCharts[5]),
				ui.NewRow(0.125, page.CPUCharts[6]),
				ui.NewRow(0.125, page.CPUCharts[7]),
				// ui.NewRow(0.125, page.CPUCharts),
			),
			ui.NewCol(0.46,
				ui.NewRow(0.34, page.MemoryChart),
				ui.NewRow(0.34, page.DiskChart),
				ui.NewRow(0.34, page.NetworkChart),
			),
		)
	} else if numCores == 4 {
		page.Grid.Set(
			ui.NewCol(0.54,
				ui.NewRow(0.25, page.CPUCharts[0]),
				ui.NewRow(0.25, page.CPUCharts[1]),
				ui.NewRow(0.25, page.CPUCharts[2]),
				ui.NewRow(0.25, page.CPUCharts[3]),
			),
			ui.NewCol(0.46,
				ui.NewRow(0.34, page.MemoryChart),
				ui.NewRow(0.34, page.DiskChart),
				ui.NewRow(0.34, page.NetworkChart),
			),
		)
	}

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)
}

var run = true

// RenderCharts handles plotting graphs and charts for system stats in general.
func RenderCharts(endChannel chan os.Signal, memChannel chan []float64, cpuChannel chan []float64, diskChannel chan [][]string, netChannel chan map[string][]float64, wg *sync.WaitGroup) {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	numCores := runtime.NumCPU()

	if numCores != 4 && numCores != 8 { // Commit die!
		endChannel <- os.Kill
		wg.Done()
		return
	}

	myPage := newPage(numCores)
	myPage.MemoryChart.BarGap = 13

	ipData := make([]float64, 65)
	opData := make([]float64, 65)

	// Bar chart for Memory
	bc := widgets.NewBarChart()
	bc.Labels = []string{"Total", "Available", "Used"}
	bc.Title = " Memory (RAM) "
	bc.BarWidth = 10
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	// Table for Disk Usage
	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.TextAlignment = ui.AlignCenter
	table.RowSeparator = false
	table.SetRect(35, 16, 80, 24)
	table.Title = " Disk "

	pause := func() {
		run = !run
		if run {
			myPage.MemoryChart.Title = " Memory (RAM) "

		} else {
			myPage.MemoryChart.Title = " Memory (Stopped) "

		}
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
				for _, value := range data {

					ipData = append(ipData, value[0])
					ipData = ipData[1:]

					opData = append(opData, value[1])
					opData = opData[1:]
				}

				temp := [][]float64{}
				temp = append(temp, ipData)
				temp = append(temp, opData)
				myPage.NetworkChart.Data = temp
			}

		case cpu_data := <-cpuChannel: // Update Gauge map with newer values
			// nproc = len(cpu_data)
			if run {
				for index, rate := range cpu_data {
					myPage.CPUCharts[index].Title = " CPU " + strconv.Itoa(index) + " "
					myPage.CPUCharts[index].Percent = int(rate)
					myPage.CPUCharts[index].BarColor = ui.ColorRed
					myPage.CPUCharts[index].BorderStyle.Fg = ui.ColorWhite
					myPage.CPUCharts[index].TitleStyle.Fg = ui.ColorCyan
				}
			}

		case <-tick:
			w, h := ui.TerminalDimensions()

			myPage.Grid.SetRect(0, 0, w, h)
			ui.Render(myPage.Grid)

		case data := <-memChannel: // Update memory values
			if run {

				bc.Data = data
			}
		}
	}
}
