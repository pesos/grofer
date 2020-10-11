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
package process

import (
	"context"

	"github.com/pesos/grofer/src/utils"
	proc "github.com/shirou/gopsutil/process"
)

// Serve serves data on a per process basis
func Serve(process *Process, dataChannel chan *Process, ctx context.Context, refreshRate int64) error {
	return utils.TickUntilDone(ctx, refreshRate, func() error {
		process.UpdateProcInfo()
		dataChannel <- process

		return nil
	})
}

func ServeProcs(dataChannel chan []*proc.Process, ctx context.Context, refreshRate int64) error {
	return utils.TickUntilDone(ctx, refreshRate, func() error {
		procs, err := proc.Processes()
		if err == nil {
			dataChannel <- procs
		}

		return nil
	})
}
