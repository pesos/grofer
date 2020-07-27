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
	table.SetRect(35, 26, 80, 31)
	table.Title = " Disk "

	ipData := make([]float64, 40)
	opData := make([]float64, 40)
	var nproc int

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

		case data := <-diskChannel: // Update disk values
			if run {
				table.Rows = data
				ui.Render(table)
			}

		case data := <-netChannel: // Update network stats & render braille plots
			if run {
				for _, value := range data {

					ipData = append(ipData, value[0])
					ipData = ipData[1:]

					opData = append(opData, value[1])
					opData = opData[1:]
				}

				pl := widgets.NewPlot()
				pl2 := widgets.NewPlot()

				temp := [][]float64{}
				temp = append(temp, opData)
				pl.Data = temp
				pl.HorizontalScale = 1
				pl.AxesColor = ui.ColorWhite
				pl.LineColors[0] = ui.ColorCyan
				pl.Title = " I/P Data "
				pl.SetRect(35, 0, 70, 13)

				temp2 := [][]float64{}
				temp2 = append(temp2, opData)
				pl2.Data = temp
				pl2.HorizontalScale = 1
                pl2.LineColors[0] = ui.ColorCyan
				//pl2.LineColors[1] = ui.ColorRed
				pl2.AxesColor = ui.ColorWhite
				pl2.Title = " O/P Data "
				pl2.SetRect(35, 13, 70, 26)

				ui.Render(pl)
				ui.Render(pl2)

			}

		case cpu_data := <-cpuChannel: // Update Gauge map with newer values
			nproc = len(cpu_data)
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

		case data := <-memChannel: // Update memory values
			if run {
				_,term_breadth := ui.TerminalDimensions()

                bc.Data = data
				bc.SetRect(0, nproc*3, 35, term_breadth)
				ui.Render(bc)
			}

		}
	}
}
