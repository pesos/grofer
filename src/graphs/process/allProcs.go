/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package process

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/pesos/grofer/src/process"
	"github.com/pesos/grofer/src/utils"
)

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
			createTime := utils.GetDateFromUnix(ctime)
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

func AllProcVisuals(dataChannel chan map[int32]*process.Process,
	endChannel chan os.Signal,
	wg *sync.WaitGroup) {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	myPage := NewAllProcsPage()

	pause := func() {
		runAllProc = !runAllProc
	}

	uiEvents := ui.PollEvents()
	tick := time.Tick(100 * time.Millisecond)

	previousKey := ""
	selectedStyle := ui.NewStyle(ui.ColorYellow, ui.ColorClear, ui.ModifierBold)

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
			myPage.BodyList.SelectedRowStyle = selectedStyle
			if runAllProc {

				myPage.BodyList.Rows = getData(data)
			}
		case <-tick: // Update page with new values
			w, h := ui.TerminalDimensions()
			ui.Clear()
			myPage.Grid.SetRect(0, 0, w, h)
			ui.Render(myPage.Grid)
		}
	}
}
