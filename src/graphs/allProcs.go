package graphs

import (
	"log"
	"os"
	"strconv"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pesos/grofer/src/process"
)

var runAllProc = true

func getData(procs map[int32]*process.Process) []string {
	var data []string
	for pid, info := range procs {
		if info.Exe != "NA" {
			temp := strconv.Itoa(int(pid))
			for i := 0; i < 21-len(strconv.Itoa(int(pid))); i++ {
				temp = temp + " "
			}
			temp = temp + "[" + info.Exe + "](fg:green,bg:black)"
			data = append(data, temp)
		}
	}
	return data
}

func AllProcVisuals(dataChannel chan map[int32]*process.Process, endChannel chan os.Signal, wg *sync.WaitGroup) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	headerTable := widgets.NewTable()
	headerTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	headerTable.Rows = [][]string{[]string{"PID", "Command"}}
	headerTable.ColumnWidths = []int{20, 137}
	headerTable.SetRect(0, 0, 158, 3)
	headerTable.RowSeparator = false
	ui.Render(headerTable)

	procTable := widgets.NewList()
	procTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	procTable.TitleStyle.Fg = ui.ColorCyan
	procTable.SetRect(0, 3, 158, 38)

	pause := func() {
		run = !run
	}
	uiEvents := ui.PollEvents()

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
				procTable.ScrollDown()
			case "k", "<Up>":
				procTable.ScrollUp()
			case "<C-d>":
				procTable.ScrollHalfPageDown()
			case "<C-u>":
				procTable.ScrollHalfPageUp()
			case "<C-f>":
				procTable.ScrollPageDown()
			case "<C-b>":
				procTable.ScrollPageUp()
			case "g":
				if previousKey == "g" {
					procTable.ScrollTop()
				}
			case "<Home>":
				procTable.ScrollTop()
			case "G", "<End>":
				procTable.ScrollBottom()
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			procTable.Rows = getData(data)
			ui.Render(procTable)
		}
	}
}
