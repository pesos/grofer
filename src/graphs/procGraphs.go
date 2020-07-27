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
	page.CPUChart.BarColor = ui.ColorRed
	page.CPUChart.BorderStyle.Fg = ui.ColorWhite
	page.CPUChart.TitleStyle.Fg = ui.ColorCyan

	// Initialize Gauge for Memory Chart
	page.MemChart.Title = " Mem % "
	page.MemChart.BarColor = ui.ColorRed
	page.MemChart.BorderStyle.Fg = ui.ColorWhite
	page.MemChart.TitleStyle.Fg = ui.ColorCyan

	// Initialize Table for PID Details Table
	page.PIDTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	page.PIDTable.TextAlignment = ui.AlignCenter
	page.PIDTable.RowSeparator = false
	page.PIDTable.Title = " PID "
	page.PIDTable.TitleStyle.Fg = ui.ColorCyan

	// Initialize List for Child Processes list
	page.ChildProcsList.Title = " Child Processes "
	page.ChildProcsList.BorderStyle.Fg = ui.ColorWhite
	page.ChildProcsList.TitleStyle.Fg = ui.ColorCyan

	// Initialize Bar Chart for CTX Swicthes Chart
	page.CTXSwitchesChart.Data = []float64{0, 0}
	page.CTXSwitchesChart.Labels = []string{"Volun", "Involun"}
	page.CTXSwitchesChart.Title = " Ctx switches "
	page.CTXSwitchesChart.BorderStyle.Fg = ui.ColorWhite
	page.CTXSwitchesChart.TitleStyle.Fg = ui.ColorCyan
	page.CTXSwitchesChart.BarWidth = 10
	page.CTXSwitchesChart.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	page.CTXSwitchesChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	page.CTXSwitchesChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Bar Chart for Page Faults Chart
	page.PageFaultsChart.Data = []float64{0, 0}
	page.PageFaultsChart.Labels = []string{"minr", "mjr"}
	page.PageFaultsChart.Title = " Page Faults "
	page.PageFaultsChart.BorderStyle.Fg = ui.ColorWhite
	page.PageFaultsChart.TitleStyle.Fg = ui.ColorCyan
	page.PageFaultsChart.BarWidth = 10
	page.PageFaultsChart.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen, ui.ColorYellow, ui.ColorCyan}
	page.PageFaultsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	page.PageFaultsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Bar Chart for Memory Stats Chart
	page.MemStatsChart.Data = []float64{0, 0, 0, 0}
	page.MemStatsChart.Labels = []string{"RSS", "Data", "Stack", "Swap"}
	page.MemStatsChart.Title = " Mem Stats (mb) "
	page.MemStatsChart.BorderStyle.Fg = ui.ColorWhite
	page.MemStatsChart.TitleStyle.Fg = ui.ColorCyan
	page.MemStatsChart.BarWidth = 10
	page.MemStatsChart.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen, ui.ColorYellow, ui.ColorCyan}
	page.MemStatsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	page.MemStatsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Initialize Grid layout
	page.Grid.Set(
		ui.NewCol(0.5,
			ui.NewRow(0.1, page.CPUChart),
			ui.NewRow(0.1, page.MemChart),
			ui.NewRow(0.4, page.PIDTable),
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
			temp = temp + "[" + exe + "](fg:green,bg:black)"
			childProcs = append(childProcs, temp)
		} else {
			childProcs = append(childProcs, strconv.Itoa(int(proc.Pid))+"            "+"NA")
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
	tick := time.Tick(100 * time.Millisecond)

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
				switches := []float64{float64(data.NumCtxSwitches.Voluntary), float64(data.NumCtxSwitches.Involuntary)}
				myPage.CTXSwitchesChart.Data = switches

				// update cpu %
				myPage.CPUChart.Percent = int(data.CPUPercent)

				// update mem %
				myPage.MemChart.Percent = int(data.MemoryPercent)

				// update proc info
				myPage.PIDTable.Rows = [][]string{
					[]string{"Name", data.Name},
					[]string{"Command", data.Exe},
					[]string{"Status", statusMap[data.Status] + " (" + data.Status + ")"},
					[]string{"Background", strconv.FormatBool(data.Background)},
					[]string{"Foreground", strconv.FormatBool(data.Foreground)},
					[]string{"Running", strconv.FormatBool(data.IsRunning)},
					[]string{"Creation Time", getDateFromUnix(data.CreateTime)},
					[]string{"Nice value", strconv.Itoa(int(data.Nice))},
					[]string{"Thread count", strconv.Itoa(int(data.NumThreads))},
					[]string{"Child process count", strconv.Itoa(len(data.Children))},
				}
				myPage.PIDTable.Title = " PID: " + strconv.Itoa(int(data.Proc.Pid)) + " "
				myPage.PIDTable.BorderStyle.Fg = ui.ColorWhite
				myPage.PIDTable.TitleStyle.Fg = ui.ColorCyan

				//update memory stats
				memData := []float64{getInMB(data.MemoryInfo.RSS, 1),
					getInMB(data.MemoryInfo.Data, 1),
					getInMB(data.MemoryInfo.Stack, 1),
					getInMB(data.MemoryInfo.Swap, 1),
				}
				myPage.MemStatsChart.Data = memData

				//update page faults
				faults := []float64{float64(data.PageFault.MinorFaults),
					float64(data.PageFault.MajorFaults),
				}
				myPage.PageFaultsChart.Data = faults
				myPage.ChildProcsList.Rows = getChildProcs(data)
			}

		case <-tick:
			w, h := ui.TerminalDimensions()

			myPage.Grid.SetRect(0, 0, w, h)
			ui.Render(myPage.Grid)
		}
	}
}
