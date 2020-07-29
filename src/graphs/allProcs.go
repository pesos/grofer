package graphs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pesos/grofer/src/process"
)

type allProcPage struct {
	Grid         *ui.Grid
	HeadingTable *widgets.Table
	BodyList     *widgets.List
}

func newProcsPage() *allProcPage {
	page := &allProcPage{
		Grid:         ui.NewGrid(),
		HeadingTable: widgets.NewTable(),
		BodyList:     widgets.NewList(),
	}
	page.init()
	return page
}

func (page *allProcPage) init() {
	page.HeadingTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	page.HeadingTable.Rows = [][]string{[]string{" PID", " Command", " CPU", " Memory", " Status", " Foreground", " Creation Time", " Thread Count"}}
	page.HeadingTable.ColumnWidths = []int{10, 40, 10, 10, 8, 12, 23, 15}
	page.HeadingTable.TextAlignment = ui.AlignLeft
	page.HeadingTable.RowSeparator = false

	page.BodyList.TextStyle = ui.NewStyle(ui.ColorWhite)
	page.BodyList.TitleStyle.Fg = ui.ColorCyan

	page.Grid.Set(
		ui.NewRow(0.12, page.HeadingTable),
		ui.NewRow(0.88, page.BodyList),
	)

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)
}

var runAllProc = true

func getData(procs map[int32]*process.Process) []string {
	var data []string
	for pid, info := range procs {
		if info.Exe != "NA" {
			temp := strconv.Itoa(int(pid))

			for i := 0; i < 12-len(strconv.Itoa(int(pid))); i++ {
				temp = temp + " "
			}

			commands := strings.Split(info.Exe, "/")
			command := commands[len(commands)-1]

			if len(command) > 40 {
				command = command[:40]
			} else {
				temp = temp + "[" + command + "](fg:green)"
				for i := 0; i < 41-len(command); i++ {
					temp = temp + " "
				}
			}

			cpuPercent := fmt.Sprintf("%.2f%s", info.CPUPercent, "%")
			temp = temp + cpuPercent
			for i := 0; i < 11-len(cpuPercent); i++ {
				temp = temp + " "
			}

			memPercent := fmt.Sprintf("%.2f%s", info.MemoryPercent, "%")
			temp = temp + memPercent
			for i := 0; i < 11-len(memPercent); i++ {
				temp = temp + " "
			}

			status := info.Status
			temp = temp + status
			for i := 0; i < 9-len(status); i++ {
				temp = temp + " "
			}

			if info.Foreground {
				temp = temp + "True"
				for i := 0; i < 9; i++ {
					temp = temp + " "
				}
			} else {
				temp = temp + "False"
				for i := 0; i < 8; i++ {
					temp = temp + " "
				}
			}

			ctime := info.CreateTime
			createTime := getDateFromUnix(ctime)
			temp = temp + createTime
			for i := 0; i < 24-len(createTime); i++ {
				temp = temp + " "
			}

			threads := info.NumThreads
			threadCount := strconv.FormatInt(int64(threads), 10)
			temp = temp + threadCount

			data = append(data, temp)
		}
	}
	return data
}

func AllProcVisuals(dataChannel chan map[int32]*process.Process, endChannel chan os.Signal, wg *sync.WaitGroup) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	myPage := newProcsPage()

	pause := func() {
		runAllProc = !runAllProc
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
				myPage.BodyList.ScrollDown()
			case "k", "<Up>":
				myPage.BodyList.ScrollUp()
			case "<C-d>":
				myPage.BodyList.ScrollHalfPageDown()
			case "<C-u>":
				myPage.BodyList.ScrollHalfPageUp()
			case "<C-f>":
				myPage.BodyList.ScrollPageDown()
			case "<C-b>":
				myPage.BodyList.ScrollPageUp()
			case "g":
				if previousKey == "g" {
					myPage.BodyList.ScrollTop()
				}
			case "<Home>":
				myPage.BodyList.ScrollTop()
			case "G", "<End>":
				myPage.BodyList.ScrollBottom()
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			if runAllProc {
				myPage.BodyList.Rows = getData(data)
			}

		case <-tick: // Update page with new values
			w, h := ui.TerminalDimensions()

			myPage.Grid.SetRect(0, 0, w, h)
			ui.Render(myPage.Grid)
		}
	}
}
