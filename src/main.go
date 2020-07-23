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
		cpu_percent, err := myProcess.CPUPercent()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(cpu_percent)
	}

}
