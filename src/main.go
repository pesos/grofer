package main

import (
	"os"
	"sync"
	"time"

	"github.com/pesos/grofer/src/general"
	proc "github.com/pesos/grofer/src/process"
)

func main() {
	procs, err := proc.InitAllProcs()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	endProcessChannel := make(chan os.Signal) // Channel to signal end of routine
	endGeneralChannel := make(chan os.Signal) // Channel to signal end of routine

	wg.Add(2) // Increment semaphore by 2 to allow new routines

	go general.GlobalStats(endGeneralChannel, &wg) // Launch routine
	go proc.Serve(procs, endProcessChannel, &wg)

	time.Sleep(5 * time.Second) // Replace with Termui code

	// signal.Notify(endProcessChannel, os.Kill) // Doesn't work for some reason
	endProcessChannel <- os.Kill
	endGeneralChannel <- os.Kill

	wg.Wait()

	// signal.Notify(endChannel, os.Interrupt)
	// <-endChannel

}
