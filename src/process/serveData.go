package process

import (
	"os"
	"sync"
)

func Serve(processes []*Process, endChannel chan os.Signal, wg *sync.WaitGroup) {
	for _, proc := range processes {
		select {
		case <-endChannel: // Stop execution if end signal received
			wg.Done()
			return

		default:
			proc.UpdateProcInfo()
			proc.PrintStats()
			//print it in some structured way
		}
	}
	return
}
