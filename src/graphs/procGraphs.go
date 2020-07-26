package graphs

import (
	"log"
	"os"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pesos/grofer/src/process"
)

var runProc = true

// ProcVisuals renders graphs and charts for per-process stats.
func ProcVisuals(endChannel chan os.Signal, dataChannel chan *process.Process, wg *sync.WaitGroup) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	//fmt.Println(len(dataChannel))
	bc := widgets.NewBarChart()
	bc.Title = "Mem stats"
	bc.SetRect(50, 0, 100, 10)
	bc.Labels = []string{"", "v", "iv"}
	bc.BarColors[0] = ui.ColorGreen
	bc.NumStyles[0] = ui.NewStyle(ui.ColorBlack)

	pause := func() {
		runProc = !runProc
		if runProc {
			bc.Title = "No. of ctx switches"

		} else {
			bc.Title = "No. of ctx switches (Stopped)"
		}
	}
	ui.Render(bc)

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
			if runProc {
				bc.Data = []float64{0, float64(data.NumCtxSwitches.Voluntary), float64(data.NumCtxSwitches.Involuntary)}
				ui.Render(bc)
			}
		}
	}
}
