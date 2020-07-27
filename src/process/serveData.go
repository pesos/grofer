package process

import (
	"os"
	"sync"
	"time"
)

var mut sync.Mutex

// Serve serves data on a per process basis
func Serve(processes map[int32]*Process, pid int32, dataChannel chan *Process, endChannel chan os.Signal, wg *sync.WaitGroup) {
	for {
		select {
		case <-endChannel:
			wg.Done() // Stop execution if end signal received
			return

		default:
			func() {
				// mut.Lock()
				processes[pid].UpdateProcInfo()
				dataChannel <- processes[pid]
				// mut.Unlock()
			}()
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
			func() {
				// mut.Lock()
				procs, err := InitAllProcs()
				if err == nil {
					for _, info := range procs {
						info.UpdateProcInfo()
					}
					dataChannel <- procs
				}
				// mut.Unlock()
			}()
		}
		time.Sleep(10 * time.Millisecond)
	}
}
