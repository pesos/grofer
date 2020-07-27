package process

import (
	"os"
	"sync"
)

var mut sync.Mutex

func Serve(processes map[int32]*Process, pid int32, dataChannel chan *Process, endChannel chan os.Signal, wg *sync.WaitGroup) {

	for {
		select {
		case <-endChannel:
			wg.Done() // Stop execution if end signal received
			return

		default:
			func() {
				mut.Lock()
				processes[pid].UpdateProcInfo()
				dataChannel <- processes[pid]
				mut.Unlock()
			}()
		}
	}

}
