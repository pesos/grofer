package process

import (
	"os"
	"sync"
	"time"
)

var mut sync.Mutex

// Serve serves data on a per process basis
func Serve(process *Process, dataChannel chan *Process, endChannel chan os.Signal, wg *sync.WaitGroup) {
	for {
		select {
		case <-endChannel:
			wg.Done() // Stop execution if end signal received
			return

		default:
			process.UpdateProcInfo()
			dataChannel <- process
		}
	}

}

func ServeProcs(dataChannel chan map[int32]*Process, endChannel chan os.Signal, wg *sync.WaitGroup) {
	for {
		select {
		case <-endChannel:
			wg.Done()
			return

		default:
			procs, err := InitAllProcs()
			if err == nil {
				for _, info := range procs {
					info.UpdateProcForVisual()
				}
				dataChannel <- procs
				time.Sleep(1 * time.Second)
			}
		}
	}
}
