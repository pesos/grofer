package graphs

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var run = true

// OverallVisuals handles plotting graphs and charts for system stats in general.
func OverallVisuals(endChannel chan os.Signal, dataChannel chan []float64, cpuChannel chan []float64, wg *sync.WaitGroup) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	nproc := len(cpuChannel)
	fmt.Println(nproc)

	pc := widgets.NewPieChart()
	pc.Title = " Memory Usage "
	pc.SetRect(0, 0, 25, 15)
	pc.Data = []float64{0, 0}
	pc.AngleOffset = -.5 * math.Pi
	pc.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f", v)
	}

	type gaugeMap map[int]*widgets.Gauge

	gMap := make(gaugeMap)

	pause := func() {
		run = !run
		if run {
			pc.Title = "Memory Usage"

		} else {
			pc.Title = "Pie Chart (Stopped)"
		}
	}

	ui.Render(pc)

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": //q or Ctrl-C to quit
				endChannel <- os.Kill
				wg.Done()
				return
			case "s": //s to stop
				pause()
			}
		case data := <-dataChannel:
			if run {
				pc.Data = data
				ui.Render(pc)
			}
		case cpu_data := <-cpuChannel:

			if run {
				for index, rate := range cpu_data {
					tempGauge := widgets.NewGauge()
					tempGauge.Title = "CPU " + strconv.Itoa(index)
					tempGauge.SetRect(25, 0+(index*3), 60, 0+((index+1)*3))
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
