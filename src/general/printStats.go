package general

import (
	"math"

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
