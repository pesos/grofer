package graphs

import (
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pesos/grofer/src/process"
)

var (
	K = math.Pow(10, 3)
	M = math.Pow(10, 6)
	G = math.Pow(10, 9)
	T = math.Pow(10, 12)
)

func roundOff(num float64, divisor float64) float64 {
	x := num / divisor
	return math.Round(x*10) / 10
}

func roundValues(num1, num2 float64) ([]float64, string) {
	nums := []float64{}
	var units string
	var n float64
	if num1 > num2 {
		n = num1
	} else {
		n = num2
	}

	switch {
	case n < K:
		nums = append(nums, num1)
		nums = append(nums, num2)
		units = " "

	case n < M:
		nums = append(nums, roundOff(num1, K))
		nums = append(nums, roundOff(num2, K))
		units = " per thousand "

	case n < G:
		nums = append(nums, roundOff(num1, M))
		nums = append(nums, roundOff(num2, M))
		units = " per million "

	case n < T:
		nums = append(nums, roundOff(num1, G))
		nums = append(nums, roundOff(num2, G))
		units = " per trillion "
	}

	return nums, units

}

type proccessPage struct {
	Grid             *ui.Grid
	CPUChart         *widgets.Gauge
	MemChart         *widgets.Gauge
	PIDTable         *widgets.Table
	ChildProcsList   *widgets.List
	CTXSwitchesChart *widgets.BarChart
	PageFaultsChart  *widgets.BarChart
	MemStatsChart    *widgets.BarChart
}

func newProcPage() *proccessPage {
	page := &proccessPage{
		Grid:             ui.NewGrid(),
		CPUChart:         widgets.NewGauge(),
		MemChart:         widgets.NewGauge(),
		PIDTable:         widgets.NewTable(),
		ChildProcsList:   widgets.NewList(),
		CTXSwitchesChart: widgets.NewBarChart(),
		PageFaultsChart:  widgets.NewBarChart(),
		MemStatsChart:    widgets.NewBarChart(),
	}
	page.init()
	return page
}

func (page *proccessPage) init() {
	// Initialize Gauge for CPU Chart
	page.CPUChart.Title = " CPU % "
	page.CPUChart.BarColor = ui.ColorGreen
	page.CPUChart.BorderStyle.Fg = ui.ColorCyan
	page.CPUChart.TitleStyle.Fg = ui.ColorWhite

	// Initialize Gauge for Memory Chart
	page.MemChart.Title = " Mem % "
	page.MemChart.BarColor = ui.ColorGreen
	page.MemChart.BorderStyle.Fg = ui.ColorCyan
	page.MemChart.TitleStyle.Fg = ui.ColorWhite

	// Initialize Table for PID Details Table
	page.PIDTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	page.PIDTable.TextAlignment = ui.AlignCenter
	page.PIDTable.RowSeparator = false
	page.PIDTable.Title = " PID "
	page.PIDTable.BorderStyle.Fg = ui.ColorCyan
	page.PIDTable.TitleStyle.Fg = ui.ColorWhite

	// Initialize List for Child Processes list
	page.ChildProcsList.Title = " Child Processes "
	page.ChildProcsList.BorderStyle.Fg = ui.ColorCyan
	page.ChildProcsList.TitleStyle.Fg = ui.ColorWhite

	// Initialize Bar Chart for CTX Swicthes Chart
	page.CTXSwitchesChart.Data = []float64{0, 0}
	page.CTXSwitchesChart.Labels = []string{"Volun", "Involun"}
	page.CTXSwitchesChart.Title = " Ctx switches "
	page.CTXSwitchesChart.BorderStyle.Fg = ui.ColorCyan
	page.CTXSwitchesChart.TitleStyle.Fg = ui.ColorWhite
	page.CTXSwitchesChart.BarWidth = 10
	page.CTXSwitchesChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	page.CTXSwitchesChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	page.CTXSwitchesChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Bar Chart for Page Faults Chart
	page.PageFaultsChart.Data = []float64{0, 0}
	page.PageFaultsChart.Labels = []string{"minr", "mjr"}
	page.PageFaultsChart.Title = " Page Faults "
	page.PageFaultsChart.BorderStyle.Fg = ui.ColorCyan
	page.PageFaultsChart.TitleStyle.Fg = ui.ColorWhite
	page.PageFaultsChart.BarWidth = 10
	page.PageFaultsChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	page.PageFaultsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	page.PageFaultsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Bar Chart for Memory Stats Chart
	page.MemStatsChart.Data = []float64{0, 0, 0, 0}
	page.MemStatsChart.Labels = []string{"RSS", "Data", "Stack", "Swap"}
	page.MemStatsChart.Title = " Mem Stats (mb) "
	page.MemStatsChart.BorderStyle.Fg = ui.ColorCyan
	page.MemStatsChart.TitleStyle.Fg = ui.ColorWhite
	page.MemStatsChart.BarWidth = 10
	page.MemStatsChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorMagenta, ui.ColorYellow, ui.ColorCyan}
	page.MemStatsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	page.MemStatsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Grid layout
	page.Grid.Set(
		ui.NewCol(0.5,
			ui.NewRow(0.125, page.CPUChart),
			ui.NewRow(0.125, page.MemChart),
			ui.NewRow(0.35, page.PIDTable),
			ui.NewRow(0.4, page.ChildProcsList),
		),
		ui.NewCol(0.5,
			ui.NewRow(0.6,
				ui.NewCol(0.5, page.CTXSwitchesChart),
				ui.NewCol(0.5, page.PageFaultsChart),
			),
			ui.NewRow(0.4, page.MemStatsChart),
		),
	)

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)
}

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

func getChildProcs(proc *process.Process) []string {
	childProcs := []string{"PID                   Command"}
	for _, proc := range proc.Children {
		exe, err := proc.Exe()
		if err == nil {
			temp := strconv.Itoa(int(proc.Pid))
			for i := 0; i < 22-len(strconv.Itoa(int(proc.Pid))); i++ {
				temp = temp + " "
			}
			temp = temp + "[" + exe + "](fg:green)"
			childProcs = append(childProcs, temp)
		} else {
			childProcs = append(childProcs, "["+strconv.Itoa(int(proc.Pid))+"](fg:yellow)"+"            "+"NA")
		}
	}

	return childProcs
}

func getDateFromUnix(createTime int64) string {
	t := time.Unix(createTime, 0)
	date := t.Format(time.RFC822)
	return date
}

// ProcVisuals renders graphs and charts for per-process stats.
func ProcVisuals(endChannel chan os.Signal, dataChannel chan *process.Process, wg *sync.WaitGroup) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	// Create new page
	myPage := newProcPage()

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
	}

	uiEvents := ui.PollEvents()
	tick := time.Tick(1 * time.Second)

	previousKey := ""
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
			case "j", "<Down>":
				myPage.ChildProcsList.ScrollDown()
			case "k", "<Up>":
				myPage.ChildProcsList.ScrollUp()
			case "<C-d>":
				myPage.ChildProcsList.ScrollHalfPageDown()
			case "<C-u>":
				myPage.ChildProcsList.ScrollHalfPageUp()
			case "<C-f>":
				myPage.ChildProcsList.ScrollPageDown()
			case "<C-b>":
				myPage.ChildProcsList.ScrollPageUp()
			case "g":
				if previousKey == "g" {
					myPage.ChildProcsList.ScrollTop()
				}
			case "<Home>":
				myPage.ChildProcsList.ScrollTop()
			case "G", "<End>":
				myPage.ChildProcsList.ScrollBottom()
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			if runProc {
				// update ctx switches
				switches, units := roundValues(float64(data.NumCtxSwitches.Voluntary), float64(data.NumCtxSwitches.Involuntary))

				myPage.CTXSwitchesChart.Data = switches
				myPage.CTXSwitchesChart.Title = " CTX Switches" + units

				// update cpu %
				myPage.CPUChart.Percent = int(data.CPUPercent)

				// update mem %
				myPage.MemChart.Percent = int(data.MemoryPercent)

				// update proc info
				myPage.PIDTable.Rows = [][]string{
					[]string{"[Name](fg:yellow)", data.Name},
					[]string{"[Command](fg:yellow)", data.Exe},
					[]string{"[Status](fg:yellow)", statusMap[data.Status] + " (" + data.Status + ")"},
					[]string{"[Background](fg:yellow)", strconv.FormatBool(data.Background)},
					[]string{"[Foreground](fg:yellow)", strconv.FormatBool(data.Foreground)},
					[]string{"[Running](fg:yellow)", strconv.FormatBool(data.IsRunning)},
					[]string{"[Creation Time](fg:yellow)", getDateFromUnix(data.CreateTime)},
					[]string{"[Nice value](fg:yellow)", strconv.Itoa(int(data.Nice))},
					[]string{"[Thread count](fg:yellow)", strconv.Itoa(int(data.NumThreads))},
					[]string{"[Child process count](fg:yellow)", strconv.Itoa(len(data.Children))},
				}
				myPage.PIDTable.Title = " PID: " + strconv.Itoa(int(data.Proc.Pid)) + " "

				//update memory stats
				memData := []float64{getInMB(data.MemoryInfo.RSS, 1),
					getInMB(data.MemoryInfo.Data, 1),
					getInMB(data.MemoryInfo.Stack, 1),
					getInMB(data.MemoryInfo.Swap, 1),
				}
				myPage.MemStatsChart.Data = memData

				//update page faults
				faults, units := roundValues(float64(data.PageFault.MinorFaults), float64(data.PageFault.MajorFaults))

				myPage.PageFaultsChart.Data = faults
				myPage.PageFaultsChart.Title = " Page Faults" + units
				myPage.ChildProcsList.Rows = getChildProcs(data)
			}

		case <-tick:
			w, h := ui.TerminalDimensions()
			ui.Clear()
			myPage.Grid.SetRect(0, 0, w, h)
			ui.Render(myPage.Grid)
		}
	}
}
