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
package cmd

import (
	"fmt"
	"os"
	"sync"

	proc "github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"

	procGraph "github.com/pesos/grofer/src/display/process"
	"github.com/pesos/grofer/src/process"
	"github.com/pesos/grofer/src/utils"
)

const (
	defaultProcRefreshRate = 3000
	defaultProcPid         = -1
)

// procCmd represents the proc command
var procCmd = &cobra.Command{
	Use:   "proc",
	Short: "proc command is used to get per-process information",
	Long: `proc command is used to get information about each running process in the system.

Syntax:
  grofer proc

To get information about a particular process whose PID is known the -p or --pid flag can be used.

Syntax:
  grofer proc -p [PID]`,
	Aliases: []string{"process", "processess"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("the proc command should have no arguments, see grofer proc --help for further info")
		}

		pid, _ := cmd.Flags().GetInt32("pid")
		procRefreshRate, _ := cmd.Flags().GetInt32("refresh")

		if procRefreshRate < 1000 {
			return fmt.Errorf("invalid refresh rate: minimum refresh rate is 1000(ms)")
		}

		var wg sync.WaitGroup

		if pid != -1 {
			endChannel := make(chan os.Signal, 1)
			dataChannel := make(chan *process.Process, 1)

			wg.Add(2)

			proc, err := process.NewProcess(pid)
			if err != nil {
				utils.ErrorMsg()
				return fmt.Errorf("invalid pid")
			}

			go process.Serve(proc, dataChannel, endChannel, int32(4*procRefreshRate/5), &wg)
			go procGraph.ProcVisuals(endChannel, dataChannel, procRefreshRate, &wg)
			wg.Wait()
		} else {
			dataChannel := make(chan []*proc.Process, 1)
			endChannel := make(chan os.Signal, 1)

			wg.Add(2)

			go process.ServeProcs(dataChannel, endChannel, int32(4*procRefreshRate/5), &wg)
			go procGraph.AllProcVisuals(dataChannel, endChannel, procRefreshRate, &wg)
			wg.Wait()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(procCmd)

	procCmd.Flags().Int32P(
		"refresh",
		"r",
		defaultProcRefreshRate,
		"Process information UI refreshes rate in milliseconds greater than 1000",
	)

	procCmd.Flags().Int32P(
		"pid",
		"p",
		defaultProcPid,
		"specify pid of process",
	)
}
