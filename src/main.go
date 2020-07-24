package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func printCPURates(cpuRates []float64) {
	var outputString string
	for index, rate := range cpuRates {
		outputString = outputString + "CPU" + strconv.Itoa(index) + ": " + strconv.Itoa(int(rate)) + "% "
	}
	fmt.Println(outputString)
}

func printMemRates(memory *mem.VirtualMemoryStat) {
	fmt.Println("Total virtual memory:", float32(memory.Total)/(1024*1024*1024), "Available:", float32(memory.Available)/(1024*1024*1024), "Used:", float32(memory.Used)/(1024*1024*1024))
}

func printIdleTime(cpuTimeStat []cpu.TimesStat) {
	fmt.Print("Idle time: ")
	for _, ind := range cpuTimeStat {
		fmt.Print(ind.CPU, ":", ind.Idle, " ")
	}
	fmt.Println()

}

// Function to print out CPU usages and memory values
func globalStats(endChannel chan int, wg *sync.WaitGroup) {
	for {
		select {
		case <-endChannel: // Stop execution if end signal received

			wg.Done()
			return

		default: // Get Memory and CPU rates per core for every 1 second

			cpuUsageRates, err := cpu.Percent(1*time.Second, true)
			if err != nil {
				log.Fatal(err)
			}

			//Times() used here for idle time, true used for per CPU idle time
			cpuTimeStat, err := cpu.Times(true)
			if err != nil {
				log.Fatal(err)
			}

			memoryStat, err := mem.VirtualMemory()
			if err != nil {
				log.Fatal(err)
			}

			// Do something with values, just printing for now
			printIdleTime(cpuTimeStat)
			printCPURates(cpuUsageRates)
			printMemRates(memoryStat)
			println()
		}
	}
}

func main() {

	var wg sync.WaitGroup

	endChannel := make(chan int, 1) // Channel to signal end of routine

	wg.Add(1) // Increment semaphore by 1 to allow new routine

	go globalStats(endChannel, &wg) // Launch routine

	time.Sleep(10 * time.Second) // A galeej way to keep the main routine busy
	endChannel <- 1              // Send signal for routine to stop

	wg.Wait()

	// if len(os.Args) != 2 {
	// 	fmt.Println("PID not entered!")
	// 	os.Exit(1)
	// }
	// arg, _ := strconv.Atoi(os.Args[1])
	// pid := int32(arg)

	// myProcess, err := process.NewProcess(pid)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// for {
	// 	status, _ := myProcess.IsRunning()
	//
	// 	if status == true {
	//
	// 		cpu_percent, err := myProcess.CPUPercent()
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	//
	// 		mem_percent, err := myProcess.MemoryPercent()
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		fmt.Println("CPU Percent: ", cpu_percent, " Memory Percent: ", mem_percent)
	// 	}
	//
	// }
}
