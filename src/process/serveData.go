package process

import (
	"os"
	"sync"
)

func Serve(processes map[int32]*Process, pid int32, dataChannel chan *Process, endChannel chan os.Signal, wg *sync.WaitGroup) {

	for {
		select {
		case <-endChannel:
			wg.Done() // Stop execution if end signal received
			return

		default:
			for {
				processes[pid].UpdateProcInfo()
				dataChannel <- processes[pid]
			}
		}
	}

}
