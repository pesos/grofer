package general

import (
	"fmt"
	//"strconv"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// PrintCPURates print the cpu rates 
func PrintCPURates(cpuRates []float64,cpuChannel chan []float64) {
	//var outputString string
	//for index, rate := range cpuRates {
	//	outputString = outputString + "CPU" + strconv.Itoa(index) + ": " + strconv.Itoa(int(rate)) + "% "
	//}
	
    //fmt.Println(outputString)
    
    var cpu_data []float64
    for _, rate := range cpuRates {
        cpu_data = append(cpu_data,float64(rate))
    }
    cpuChannel <- cpu_data
}

// PrintMemRates prints stats about the memory
func PrintMemRates(memory *mem.VirtualMemoryStat, dataChannel chan []float64) {
	// fmt.Println("Total virtual memory:", float64(memory.Total)/(1024*1024*1024),
	// "Available:", float64(memory.Available)/(1024*1024*1024),
	// "Used:", float64(memory.Used)/(1024*1024*1024))
	data := []float64{float64(memory.Available) / (1024 * 1024 * 1024), float64(memory.Used) / (1024 * 1024 * 1024)}
	dataChannel <- data
}

// PrintIdleTime prints idle time per CPU
func PrintIdleTime(cpuTimeStat []cpu.TimesStat) {
	fmt.Print("Idle time: ")
	for _, ind := range cpuTimeStat {
		fmt.Print(ind.CPU, ":", ind.Idle, " ")
	}
	fmt.Println()

}
