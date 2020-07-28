package graphs

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type mainPage struct {
	Grid         *ui.Grid
	MemoryChart  *widgets.BarChart
	DiskChart    *widgets.Table
	NetworkChart *widgets.Plot
	CPUCharts    []*widgets.Gauge
	NetPara      *widgets.Paragraph
}

func newPage(numCores int) *mainPage {
	page := &mainPage{
		Grid:         ui.NewGrid(),
		MemoryChart:  widgets.NewBarChart(),
		DiskChart:    widgets.NewTable(),
		NetworkChart: widgets.NewPlot(),
		CPUCharts:    make([]*widgets.Gauge, 0),
		NetPara:      widgets.NewParagraph(),
	}
	page.init(numCores)
	return page
}

func (page *mainPage) init(numCores int) {

	// Initialize Bar Graph for Memory Chart
	page.MemoryChart.Title = " Memory (RAM) "
	page.MemoryChart.Labels = []string{"Total", "Available", "Used", "Free"}
	page.MemoryChart.BarWidth = 8
	page.MemoryChart.BarGap = 9
	page.MemoryChart.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	page.MemoryChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	page.MemoryChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	// Initialize Table for Disk Chart
	page.DiskChart.Title = " Disk "
	page.DiskChart.TextStyle = ui.NewStyle(ui.ColorWhite)
	page.DiskChart.TextAlignment = ui.AlignLeft
	page.DiskChart.RowSeparator = false
	page.DiskChart.ColumnWidths = []int{9, 9, 9, 9, 9, 11}

	// Initialize Plot for Network Chart
	page.NetworkChart.Title = " Network data(in mB) "
	page.NetworkChart.HorizontalScale = 1
	page.NetworkChart.AxesColor = ui.ColorWhite
	page.NetworkChart.LineColors[0] = ui.ColorCyan
	page.NetworkChart.LineColors[1] = ui.ColorRed
	page.NetworkChart.DrawDirection = 1
	page.NetworkChart.DataLabels = []string{"ip kB", "op kB"} //refer issue #214 for details

	//Initialize paragraph for NetPara
	page.NetPara.Text = "[Received(kB)](fg:cyan)\n\n[Sent(kB)](fg:red)"
	page.NetPara.Border = true
	page.NetPara.Title = " RX/TX "

	// Initialize Gauges for each CPU Core usage
	for i := 0; i < numCores; i++ {
		tempGauge := widgets.NewGauge()
		tempGauge.Title = " CPU " + strconv.Itoa(i) + " "
		tempGauge.Percent = 0
		tempGauge.BarColor = ui.ColorRed
		tempGauge.BorderStyle.Fg = ui.ColorWhite
		tempGauge.TitleStyle.Fg = ui.ColorCyan
		page.CPUCharts = append(page.CPUCharts, tempGauge)
	}

	// Initialize Grid layout
	if numCores == 8 {
		page.Grid.Set(
			ui.NewCol(0.54,
				ui.NewRow(0.125, page.CPUCharts[0]),
				ui.NewRow(0.125, page.CPUCharts[1]),
				ui.NewRow(0.125, page.CPUCharts[2]),
				ui.NewRow(0.125, page.CPUCharts[3]),
				ui.NewRow(0.125, page.CPUCharts[4]),
				ui.NewRow(0.125, page.CPUCharts[5]),
				ui.NewRow(0.125, page.CPUCharts[6]),
				ui.NewRow(0.125, page.CPUCharts[7]),
			),
			ui.NewCol(0.46,
				ui.NewRow(0.34, page.MemoryChart),
				ui.NewRow(0.34,
					ui.NewCol(0.25, page.NetPara),
					ui.NewCol(0.75, page.NetworkChart),
				),
				ui.NewRow(0.34, page.DiskChart),
			),
		)
	} else if numCores == 4 {
		page.Grid.Set(
			ui.NewCol(0.54,
				ui.NewRow(0.25, page.CPUCharts[0]),
				ui.NewRow(0.25, page.CPUCharts[1]),
				ui.NewRow(0.25, page.CPUCharts[2]),
				ui.NewRow(0.25, page.CPUCharts[3]),
			),
			ui.NewCol(0.46,
				ui.NewRow(0.34, page.MemoryChart),
				ui.NewRow(0.34,
					ui.NewCol(0.25, page.NetPara),
					ui.NewCol(0.75, page.NetworkChart),
				),
				ui.NewRow(0.34, page.DiskChart),
			),
		)
	}

	w, h := ui.TerminalDimensions()
	page.Grid.SetRect(0, 0, w, h)
}

var run = true

// RenderCharts handles plotting graphs and charts for system stats in general.
func RenderCharts(endChannel chan os.Signal, memChannel chan []float64, cpuChannel chan []float64, diskChannel chan [][]string, netChannel chan map[string][]float64, wg *sync.WaitGroup) {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var totalBytesRecv float64
	var totalBytesSent float64

	// Get number of cores in machine
	numCores := runtime.NumCPU()

	if numCores != 4 && numCores != 8 { // Commit die!
		endChannel <- os.Kill
		wg.Done()
		return
	}

	// Create new page
	myPage := newPage(numCores)

	// Initialize slices for Network Data
	ipData := make([]float64, 5)
	opData := make([]float64, 5)

	// Pause to pause updating data
	pause := func() {
		run = !run
	}

	uiEvents := ui.PollEvents()
	tick := time.Tick(100 * time.Millisecond)
	for {
		select {
		case e := <-uiEvents: // For keyboard events
			switch e.ID {
			case "q", "<C-c>": // q or Ctrl-C to quit
				endChannel <- os.Kill
				wg.Done()
				return

			case "s": // s to stop
				pause()
			}

		case data := <-memChannel: // Update memory values
			if run {
				myPage.MemoryChart.Data = data
			}

		case data := <-diskChannel: // Update disk values
			if run {
				myPage.DiskChart.Rows = data
			}

		case data := <-netChannel: // Update network stats & render braille plots
			if run {

				var curBytesRecv, curBytesSent float64

				for _, netInterface := range data {
					curBytesRecv += netInterface[1]
					curBytesSent += netInterface[0]
				}

				var recentBytesRecv, recentBytesSent float64

				if totalBytesRecv != 0 {
					recentBytesRecv = curBytesRecv - totalBytesRecv
					recentBytesSent = curBytesSent - totalBytesSent

					if int(recentBytesRecv) < 0 {
						recentBytesRecv = 0
					}
					if int(recentBytesSent) < 0 {
						recentBytesSent = 0
					}

					ipData = append(ipData, recentBytesRecv)
					opData = append(opData, recentBytesSent)
				}

				totalBytesRecv = curBytesRecv
				totalBytesSent = curBytesSent

				titles := make([]string, 2)

				for i := 0; i < 2; i++ {
					if i == 0 {
						titles[i] = fmt.Sprintf("[Total RX](fg:cyan): %5.1f %s\n\n", totalBytesRecv/1024, "mB")
					} else {
						titles[i] = fmt.Sprintf("[Total TX](fg:red): %5.1f %s", totalBytesSent/1024, "mB")
					}

				}

				myPage.NetPara.Text = titles[0] + titles[1]

				temp := [][]float64{}
				temp = append(temp, ipData)
				temp = append(temp, opData)
				myPage.NetworkChart.Data = temp
			}

		case cpu_data := <-cpuChannel: // Update Gauge map with newer values
			if run {
				for index, rate := range cpu_data {
					myPage.CPUCharts[index].Title = " CPU " + strconv.Itoa(index) + " "
					myPage.CPUCharts[index].Percent = int(rate)
					myPage.CPUCharts[index].BarColor = ui.ColorRed
					myPage.CPUCharts[index].BorderStyle.Fg = ui.ColorWhite
					myPage.CPUCharts[index].TitleStyle.Fg = ui.ColorCyan
				}
			}

		case <-tick: // Update page with new values
			w, h := ui.TerminalDimensions()

			myPage.Grid.SetRect(0, 0, w, h)
			ui.Render(myPage.Grid)

		}
	}
}
