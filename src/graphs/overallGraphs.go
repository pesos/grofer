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

// RenderCharts handles plotting graphs and charts for system stats in general.
func RenderCharts(endChannel chan os.Signal, memChannel chan []float64, cpuChannel chan []float64, diskChannel chan [][]string, netChannel chan map[string][]float64, wg *sync.WaitGroup) {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// sinData := func() [][]float64 {
	// 	n := 220
	// 	data := make([][]float64, 2)
	// 	data[0] = make([]float64, n)
	// 	data[1] = make([]float64, n)
	// 	for i := 0; i < n; i++ {
	// 		data[0][i] = 1 + math.Sin(float64(i)/5)
	// 		data[1][i] = 1 + math.Cos(float64(i)/5)
	// 	}
	// 	return data
	// }()

	// Bar chart for Memory
	bc := widgets.NewBarChart()
	bc.Labels = []string{"Total", "Available", "Used"}
	bc.Title = " Memory (RAM) "
	bc.SetRect(35, 0, 65, 10)
	bc.BarWidth = 8
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	// Table for Disk Usage
	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.TextAlignment = ui.AlignCenter
	table.RowSeparator = false
	table.SetRect(35, 19, 80, 24)
	table.Title = " Disk "

	// Scatter Plot for Network stats
	sp := widgets.NewPlot()
	sp.Title = " Network Usage "
	sp.Marker = widgets.MarkerDot
	sp.Data = make([][]float64, 2)
	sp.Data[0] = make([]float64, 40)
	sp.Data[1] = make([]float64, 40)
	sp.SetRect(35, 10, 80, 19)
	sp.AxesColor = ui.ColorWhite
	sp.LineColors[0] = ui.ColorCyan
	sp.PlotType = widgets.ScatterPlot

	// Gauges for CPU core usage
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
				bc.Data = data
				ui.Render(bc)
			}

		case data := <-diskChannel: // Update disk values
			if run {
				table.Rows = data
				ui.Render(table)
			}

		case data := <-netChannel: // Update network stats & render dual sparkline
			if run {
				for _, value := range data {
					sp.Data[0] = append(sp.Data[0], value[0])
					sp.Data[0] = sp.Data[0][1:]

					sp.Data[1] = append(sp.Data[1], value[1])
					sp.Data[1] = sp.Data[1][1:]
				}

				ui.Render(sp)
			}

		case cpu_data := <-cpuChannel: // Update Gauge map with newer values
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
