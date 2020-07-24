package process

import (
	"os"
	"sync"
)

func Serve(processes []*Process, endChannel chan os.Signal, wg *sync.WaitGroup) {

	for {
		select {
		case <-endChannel:
			wg.Done() // Stop execution if end signal received
			return

		default:
			for _, proc := range processes {
				proc.UpdateProcInfo()
				proc.PrintStats() // print it in some structured way
			}
		}
	}

}
