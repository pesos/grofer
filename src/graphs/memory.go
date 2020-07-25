package graphs

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var run = true

func renderMemoryChart(endChannel chan os.Signal, dataChannel chan []float64, wg *sync.WaitGroup) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	pc := widgets.NewPieChart()
	pc.Title = " Memory Usage "
	pc.SetRect(0, 0, 25, 15)
	pc.Data = []float64{0, 0}
	pc.AngleOffset = -.5 * math.Pi
	pc.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f", v)
	}

	pause := func() {
		run = !run
		if run {
			pc.Title = "Memory Usage"
		} else {
			pc.Title = "Pie Chart (Stopped)"
		}
		ui.Render(pc)
	}

	ui.Render(pc)

	uiEvents := ui.PollEvents()
	// ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				endChannel <- os.Kill
				wg.Done()
				return
			case "s":
				pause()
			}
		case data := <-dataChannel:
			if run {
				pc.Data = data
				ui.Render(pc)
			}
		}
	}
}
