package graphs

import (
	"log"
	"math"
	"os"
	"strconv"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pesos/grofer/src/process"
)

var runProc = true

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func trim(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func getInMB(bytes uint64, precision int) float64 {
	temp := float64(bytes) / 1000000
	return trim(temp, precision)
}

// ProcVisuals renders graphs and charts for per-process stats.
func ProcVisuals(endChannel chan os.Signal, dataChannel chan *process.Process, wg *sync.WaitGroup) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	bc := widgets.NewBarChart()
	bc.Data = []float64{0, 0}
	bc.Labels = []string{"Volun", "Involun"}
	bc.Title = " Ctx switches "
	bc.SetRect(65, 0, 45, 10)
	bc.BarWidth = 8
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	cpuPercGauge := widgets.NewGauge()
	memPercGauge := widgets.NewGauge()

	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.TextAlignment = ui.AlignCenter
	table.RowSeparator = false
	table.SetRect(0, 6, 45, 20)
	table.Title = " PID "
	table.TitleStyle.Fg = ui.ColorCyan

	memStat := widgets.NewBarChart()
	memStat.Data = []float64{0, 0, 0, 0}
	memStat.Labels = []string{"RSS", "Data", "Stack", "Swap"}
	memStat.Title = " Mem Stats "
	memStat.SetRect(45, 10, 85, 20)
	memStat.BarWidth = 8
	memStat.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen, ui.ColorYellow, ui.ColorCyan}
	memStat.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	memStat.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	pageFaults := widgets.NewBarChart()
	pageFaults.Data = []float64{0, 0}
	pageFaults.Labels = []string{"Minor", "Major"}
	pageFaults.Title = " Page Faults "
	pageFaults.SetRect(85, 0, 65, 10)
	pageFaults.BarWidth = 8
	pageFaults.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	pageFaults.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	pageFaults.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	var statusMap map[string]string = map[string]string{
		"R": "Running",
		"S": "Sleep",
		"Z": "Zombie",
		"T": "Stop",
		"I": "Idle",
		"W": "Wait",
		"L": "Lock",
	}
	pause := func() {
		run = !run
		if run {
			bc.Title = " No. of context switches "

		} else {
			bc.Title = " No. of context switches (Stopped) "

		}
	}

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": //q or Ctrl-C to quit
				endChannel <- os.Kill
				ui.Close()
				wg.Done()
				return
			case "s": //s to pause
				pause()
			}
		case data := <-dataChannel:
			if runProc {
				// update ctx switches
				switches := []float64{float64(data.NumCtxSwitches.Voluntary), float64(data.NumCtxSwitches.Involuntary)}
				bc.Data = switches
				ui.Render(bc)

				// update cpu %
				cpuPercGauge.Title = " CPU % "
				cpuPercGauge.SetRect(0, 0, 45, 3)
				cpuPercGauge.Percent = int(data.CPUPercent)
				cpuPercGauge.BarColor = ui.ColorRed
				cpuPercGauge.BorderStyle.Fg = ui.ColorWhite
				cpuPercGauge.TitleStyle.Fg = ui.ColorCyan
				ui.Render(cpuPercGauge)

				// update mem %
				memPercGauge.Title = " Mem % "
				memPercGauge.SetRect(0, 3, 45, 6)
				memPercGauge.Percent = int(data.MemoryPercent)
				memPercGauge.BarColor = ui.ColorRed
				memPercGauge.BorderStyle.Fg = ui.ColorWhite
				memPercGauge.TitleStyle.Fg = ui.ColorCyan
				ui.Render(memPercGauge)

				// update proc info
				table.Rows = [][]string{
					[]string{"Attribute", "Value"},
					[]string{"", ""},
					[]string{"Name", data.Name},
					[]string{"Status", statusMap[data.Status] + " (" + data.Status + ")"},
					[]string{"Background", strconv.FormatBool(data.Background)},
					[]string{"Foreground", strconv.FormatBool(data.Foreground)},
					[]string{"Running", strconv.FormatBool(data.IsRunning)},
					[]string{"Creation Time", strconv.FormatInt(data.CreateTime, 10)},
					[]string{"Foreground", strconv.FormatBool(data.Foreground)},
					[]string{"Nice value", strconv.Itoa(int(data.Nice))},
					[]string{"Thread Count", strconv.Itoa(int(data.NumThreads))},
				}
				table.Title = " PID: " + strconv.Itoa(int(data.Proc.Pid)) + " "
				table.BorderStyle.Fg = ui.ColorWhite
				table.TitleStyle.Fg = ui.ColorCyan
				ui.Render(table)

				//update memory stats
				memData := []float64{getInMB(data.MemoryInfo.RSS, 1),
					getInMB(data.MemoryInfo.Data, 1),
					getInMB(data.MemoryInfo.Stack, 1),
					getInMB(data.MemoryInfo.Swap, 1),
				}
				memStat.Data = memData
				ui.Render(memStat)

				//update page faults
				faults := []float64{float64(data.PageFault.MinorFaults), float64(data.PageFault.MajorFaults)}
				pageFaults.Data = faults
				ui.Render(pageFaults)
			}
		}
	}
}
