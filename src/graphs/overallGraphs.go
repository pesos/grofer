package graphs

import (
	"log"
	"os"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type gaugeMap map[int]*widgets.Gauge

type mainPage struct {
	Grid        *ui.Grid
	MemoryChart *widgets.BarChart
	// DiskChart   *widgets.Table
	// NetworkChart *widgets.SparklineGroup
	// CPUCharts gaugeMap
}

func newPage() *mainPage {
	page := &mainPage{
		Grid:        ui.NewGrid(),
		MemoryChart: widgets.NewBarChart(),
		// DiskChart:   widgets.NewTable(),
		// NetworkChart: widgets.NewSparklineGroup(),
		// CPUCharts: make(gaugeMap),
	}
	return page
}

func (page *mainPage) init() {
	page.MemoryChart.Title = " Memory (RAM) "
	page.MemoryChart.Labels = []string{"Total", "Available", "Used"}
	page.MemoryChart.BarWidth = 4
	page.MemoryChart.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	page.MemoryChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	page.MemoryChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	// page.DiskChart.Title = " Disk "
	// page.DiskChart.TextStyle = ui.NewStyle(ui.ColorWhite)
	// page.DiskChart.TextAlignment = ui.AlignCenter
	// page.DiskChart.RowSeparator = false

	page.Grid.Set(
		ui.NewRow(.5,
			// ui.NewCol(0.5, page.DiskChart),
			ui.NewCol(0.5, page.MemoryChart),
		),
	)
}

var run = true

// RenderCharts handles plotting graphs and charts for system stats in general.
func RenderCharts(endChannel chan os.Signal, memChannel chan []float64, cpuChannel chan []float64, diskChannel chan [][]string, netChannel chan map[string][]float64, wg *sync.WaitGroup) {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	myPage := newPage()

	myPage.init()

	pause := func() {
		run = !run
		if run {
			myPage.MemoryChart.Title = " Memory (RAM) "

		} else {
			myPage.MemoryChart.Title = " Memory (Stopped) "

		}
	}

	uiEvents := ui.PollEvents()

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
				ui.Render(myPage.Grid)
			}

			// case <-diskChannel: // Update disk values
			// 	if run {
			// 		data := [][]string{[]string{"Mount", "Total", "Used", "Used %"}, []string{"Mount", "Total", "Used", "Used %"}}
			// 		myPage.DiskChart.Rows = data
			// 		ui.Render(myPage.Grid)
			// 	}

		}
	}
}
