package general

import (
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// PrintCPURates print the cpu rates :O
func PrintCPURates(cpuRates []float64) {
	var outputString string
	for index, rate := range cpuRates {
		outputString = outputString + "CPU" + strconv.Itoa(index) + ": " + strconv.Itoa(int(rate)) + "% "
	}
	fmt.Println(outputString)
}

// PrintMemRates prints stats about the memory
func PrintMemRates(memory *mem.VirtualMemoryStat) {
	fmt.Println("Total virtual memory:", float32(memory.Total)/(1024*1024*1024),
		"Available:", float32(memory.Available)/(1024*1024*1024),
		"Used:", float32(memory.Used)/(1024*1024*1024))
}

// PrintIdleTime prints idle time per CPU
func PrintIdleTime(cpuTimeStat []cpu.TimesStat) {
	fmt.Print("Idle time: ")
	for _, ind := range cpuTimeStat {
		fmt.Print(ind.CPU, ":", ind.Idle, " ")
	}
	fmt.Println()

}
