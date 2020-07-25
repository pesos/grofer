package graphs

import (
	"log"
	"os"
	"strconv"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var run = true

func RenderCharts(endChannel chan os.Signal, memChannel chan []float64, cpuChannel chan []float64, diskChannel chan [][]string, wg *sync.WaitGroup) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.TextAlignment = ui.AlignCenter
	table.RowSeparator = false
	table.SetRect(35, 15, 80, 20)
	table.Title = " Disk "

	bc := widgets.NewBarChart()
	bc.Data = []float64{3, 2}
	bc.Labels = []string{"Total", "Used"}
	bc.Title = " Memory (RAM) "
	bc.SetRect(35, 0, 70, 10)
	bc.BarWidth = 7
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	type gaugeMap map[int]*widgets.Gauge

	gMap := make(gaugeMap)

	pause := func() {
		run = !run
		if run {
			bc.Title = " Memory (RAM) "

		} else {
			bc.Title = " Memory (Stopped) "

		}
	}

	// ui.Render(pc)

	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": // q or Ctrl-C to quit
				endChannel <- os.Kill
				wg.Done()
				return
			case "s": // s to stop
				pause()
			}
		case data := <-memChannel:
			if run {
				bc.Data = data
				ui.Render(bc)
			}
		case data := <-diskChannel:

			if run {
				table.Rows = data
				ui.Render(table)
			}
		case cpu_data := <-cpuChannel:

			if run {
				for index, rate := range cpu_data {
					tempGauge := widgets.NewGauge()
					tempGauge.Title = " CPU " + strconv.Itoa(index)
					tempGauge.SetRect(0, 0+(index*3), 35, 0+((index+1)*3))
					tempGauge.Percent = int(rate)
					tempGauge.BarColor = ui.ColorRed
					tempGauge.BorderStyle.Fg = ui.ColorWhite
					tempGauge.TitleStyle.Fg = ui.ColorCyan
					gMap[index] = tempGauge
					ui.Render(gMap[index])
				}
			}
		}
	}
}
