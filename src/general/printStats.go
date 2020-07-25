package general

import (
	"github.com/shirou/gopsutil/mem"
)

// PrintCPURates print the cpu rates
func PrintCPURates(cpuRates []float64, cpuChannel chan []float64) {
	cpuChannel <- cpuRates
}

// PrintMemRates prints stats about the memory
func PrintMemRates(memory *mem.VirtualMemoryStat, dataChannel chan []float64) {
	data := []float64{float64(memory.Available) / (1024 * 1024 * 1024), float64(memory.Used) / (1024 * 1024 * 1024)}
	dataChannel <- data
}
