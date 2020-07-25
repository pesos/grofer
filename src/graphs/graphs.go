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

func RenderCharts(endChannel chan os.Signal, memChannel chan []float64, cpuChannel chan []float64, diskChannel chan [][]string, netChannel chan map[string][]float64, wg *sync.WaitGroup) {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.TextAlignment = ui.AlignCenter
	table.RowSeparator = false
	table.SetRect(35, 19, 80, 24)
	table.Title = " Disk "

	bc := widgets.NewBarChart()
	bc.Data = []float64{3, 2}
	bc.Labels = []string{"Total", "Available", "Used"}
	bc.Title = " Memory (RAM) "
	bc.SetRect(35, 0, 65, 10)
	bc.BarWidth = 8
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	type gaugeMap map[int]*widgets.Gauge
	type netStatMap map[int]*widgets.Sparkline
	// var ipData, opData []float64

	ipData := []float64{}
	opData := []float64{}

	var prevIp float64 = 0
	var prevOp float64 = 0

	gMap := make(gaugeMap)

	// netMap := make(netStatMap)

	sl1 := widgets.NewSparkline()
	sl1.Title = "Bytes Sent"
	sl1.LineColor = ui.ColorRed

	sl2 := widgets.NewSparkline()
	sl2.Title = "Bytes Received"
	sl2.LineColor = ui.ColorMagenta

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
		case data := <-netChannel:
			if run {
				for _, value := range data {

					if value[0] == prevIp {
						ipData = append(ipData, 0)
					} else {
						ipData = append(ipData, value[0])
						prevIp = value[0]
					}

					if len(ipData) >= 200 {
						ipData = ipData[1:]
					}

					if value[1] == prevOp {
						opData = append(opData, 0)
					} else {
						opData = append(opData, value[1])
						prevOp = value[1]
					}

					if len(opData) >= 200 {
						opData = opData[1:]
					}
				}
				sl1.Data = ipData
				sl2.Data = opData
				slg1 := widgets.NewSparklineGroup(sl1, sl2)
				slg1.Title = " Network "
				slg1.SetRect(35, 10, 80, 19)
				ui.Render(slg1)
				// fmt.Println(ipData)
				// fmt.Println(opData)
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
