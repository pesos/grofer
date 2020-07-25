package general

import (
	"fmt"
	"math"
	"strings"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// PrintCPURates print the cpu rates
func PrintCPURates(cpuRates []float64, cpuChannel chan []float64) {
	cpuChannel <- cpuRates
}

func roundOff(num uint64) float64 {
	x := float64(num) / (1024 * 1024 * 1024)
	return math.Round(x*10) / 10
}

// PrintMemRates prints stats about the memory
func PrintMemRates(memory *mem.VirtualMemoryStat, dataChannel chan []float64) {
	data := []float64{roundOff(memory.Total), roundOff(memory.Used)}
	dataChannel <- data
}

func PrintDiskRates(partitions []disk.PartitionStat, dataChannel chan [][]string) {
	rows := [][]string{[]string{"Mount", "Total", "Used", "Used %"}}
	for _, value := range partitions {
		usageVals, _ := disk.Usage(value.Mountpoint)
		stats := strings.Split(usageVals.String(), ",")[1]
		// fmt.Println(stats)
		if strings.Contains(stats, "ext") {

			path := usageVals.Path
			total := fmt.Sprintf("%.2f G", float64(usageVals.Total)/(1024*1024*1024))
			used := fmt.Sprintf("%.2f G", float64(usageVals.Used)/(1024*1024*1024))
			usedPercent := fmt.Sprintf("%.2f %s", usageVals.UsedPercent, "%")
			row := []string{path, total, used, usedPercent}
			rows = append(rows, row)
		}
	}
	dataChannel <- rows
}
