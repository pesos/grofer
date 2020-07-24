package general

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// GlobalStats gets stats about the mem and the CPUs and prints it.
func GlobalStats(endChannel chan os.Signal, wg *sync.WaitGroup) {
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
			PrintIdleTime(cpuTimeStat)
			PrintCPURates(cpuUsageRates)
			PrintMemRates(memoryStat)
			fmt.Println()
		}
	}
}
