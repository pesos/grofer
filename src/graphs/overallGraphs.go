package graphs

import (
	"log"
	"os"
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
	page.MemoryChart.Labels = []string{"Total", "Available", "Used"}
	page.MemoryChart.BarWidth = 8
	page.MemoryChart.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	page.MemoryChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	page.MemoryChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	page.DiskChart.Title = " Disk "
	page.DiskChart.TextStyle = ui.NewStyle(ui.ColorWhite)
	page.DiskChart.TextAlignment = ui.AlignCenter
	page.DiskChart.RowSeparator = false
	page.DiskChart.ColumnWidths = []int{8, 8, 8, 8}

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

	page.Grid.Set(
		ui.NewCol(0.46,
			ui.NewRow(0.34, page.MemoryChart),
			ui.NewRow(0.34, page.DiskChart),
			ui.NewRow(0.34, page.NetworkChart),
		),
		ui.NewCol(0.54,
			ui.NewRow(0.125, page.CPUCharts[0]),
			ui.NewRow(0.125, page.CPUCharts[1]),
		),
	)
	// page.Grid.Items[1].Entry = page.CPUCharts
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

	// numCores := 8

	myPage := newPage(8)

	// myPage.init()

	ipData := make([]float64, 40)
	opData := make([]float64, 40)

	pause := func() {
		run = !run
		if run {
			myPage.MemoryChart.Title = " Memory (RAM) "

		} else {
			myPage.MemoryChart.Title = " Memory (Stopped) "

		}
	}

	uiEvents := ui.PollEvents()
	tick := time.Tick(time.Second)
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

		case <-tick:
			w, h := ui.TerminalDimensions()

			myPage.Grid.SetRect(0, 0, w, h)
			ui.Render(myPage.Grid)
		}
	}
}
