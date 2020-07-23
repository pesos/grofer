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
		fmt.Println("Usage statement here")
		os.Exit(1)
	}

	arg, _ := strconv.Atoi(os.Args[1])
	pid := int32(arg)
	myProcess, err := process.NewProcess(pid)
	if err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Println(myProcess.CPUPercent())
	}

}
