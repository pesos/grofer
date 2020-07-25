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

func RenderMemoryChart(endChannel chan os.Signal, dataChannel chan []float64, cpuChannel chan []float64, wg *sync.WaitGroup) {
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

	// g0 := widgets.NewGauge()
	// g0.Title = "CPU 0"
	// g0.SetRect(0, 15, 40, 18)
	// g0.Percent = 0
	// g0.BarColor = ui.ColorRed
	// g0.BorderStyle.Fg = ui.ColorWhite
	// g0.TitleStyle.Fg = ui.ColorCyan
	//
	// g1 := widgets.NewGauge()
	// g1.Title = "CPU 1"
	// g1.SetRect(40, 15, 80, 18)
	// g1.Percent = 0
	// g1.BarColor = ui.ColorRed
	// g1.BorderStyle.Fg = ui.ColorWhite
	// g1.TitleStyle.Fg = ui.ColorCyan
	//
	// g2 := widgets.NewGauge()
	// g2.Title = "CPU 2"
	// g2.SetRect(0, 19, 40, 22)
	// g2.Percent = 0
	// g2.BarColor = ui.ColorRed
	// g2.BorderStyle.Fg = ui.ColorWhite
	// g2.TitleStyle.Fg = ui.ColorCyan
	//
	// g3 := widgets.NewGauge()
	// g3.Title = "CPU 3"
	// g3.SetRect(40, 19, 80, 22)
	// g3.Percent = 0
	// g3.BarColor = ui.ColorRed
	// g3.BorderStyle.Fg = ui.ColorWhite
	// g3.TitleStyle.Fg = ui.ColorCyan

	pause := func() {
		run = !run
		if run {
			pc.Title = "Memory Usage"
			// g0.Title = "CPU Percentage"
			// g1.Title = g0.Title
			// g2.Title = g0.Title
			// g3.Title = g0.Title

		} else {
			pc.Title = "Pie Chart (Stopped)"
			// g0.Title = "Gauge (stopped)"
			// g1.Title = g0.Title
			// g2.Title = g0.Title
			// g3.Title = g0.Title
		}
		// ui.Render(pc)
		// ui.Render(g0)
		// ui.Render(g1)
		// ui.Render(g2)
		// ui.Render(g3)
	}

	ui.Render(pc)
	// ui.Render(g0)
	// ui.Render(g1)
	// ui.Render(g2)
	// ui.Render(g3)

	uiEvents := ui.PollEvents()
	// ticker := time.NewTicker(time.Second).C
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
