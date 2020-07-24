package main

import (
	"os"
	"os/signal"
	"sync"

	proc "github.com/pesos/go-htop/src/process"
)

func main() {
	procs, err := proc.InitAllProcs()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	endChannel := make(chan os.Signal) // Channel to signal end of routine

	wg.Add(1) // Increment semaphore by 1 to allow new routine

	//go general.GlobalStats(endChannel, &wg) // Launch routine
	go proc.Serve(procs, endChannel, &wg)

	wg.Wait()

	signal.Notify(endChannel, os.Interrupt)
	signal.Notify(endChannel, os.Kill)

	<-endChannel
}
