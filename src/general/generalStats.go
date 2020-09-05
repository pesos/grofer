/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package general

import (
	"os"
	"sync"
	"time"

	"github.com/pesos/grofer/src/utils"
)

// GlobalStats gets stats about the mem and the CPUs and prints it.
func GlobalStats(endChannel chan os.Signal,
	dataChannel chan utils.DataStats, refreshRate int32,
	wg *sync.WaitGroup) {

	for {
		select {
		case <-endChannel: // Stop execution if end signal received
			wg.Done()
			return

		default: // Get Memory and CPU rates per core periodically

			go ServeCPURates(dataChannel)
			go ServeMemRates(dataChannel)
			go ServeDiskRates(dataChannel)
			ServeNetRates(dataChannel)
			time.Sleep(time.Duration(refreshRate) * time.Millisecond)
		}
	}
}
