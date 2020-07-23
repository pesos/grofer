package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shirou/gopsutil/process"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("PID not entered!")
		os.Exit(1)
	}

	arg, _ := strconv.Atoi(os.Args[1])
	pid := int32(arg)

	myProcess, err := process.NewProcess(pid)
	if err != nil {
		log.Fatal(err)
	}

	for {
		status, _ := myProcess.IsRunning()

		if status == true {

			cpu_percent, err := myProcess.CPUPercent()
			if err != nil {
				log.Fatal(err)
			}

			mem_percent, err := myProcess.MemoryPercent()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("CPU Percent: ", cpu_percent, " Memory Percent: ", mem_percent)
		}

	}
}
