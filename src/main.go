package main

import (
	"os"
	"sync"
	"time"

	proc "github.com/pesos/grofer/src/process"
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

	time.Sleep(5 * time.Second) // Replace with Termui code

	// signal.Notify(endChannel, os.Kill) // Doesn't work for some reason
	endChannel <- os.Kill

	wg.Wait()

	// signal.Notify(endChannel, os.Interrupt)
	// <-endChannel

}
