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

	"github.com/pesos/grofer/src/utils"

	procGraph "github.com/pesos/grofer/src/graphs/process"
	"github.com/pesos/grofer/src/process"
	"github.com/spf13/cobra"
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

			go process.Serve(proc, dataChannel, endChannel, &wg)
			go procGraph.ProcVisuals(endChannel, dataChannel, &wg)
			wg.Wait()
		} else {
			dataChannel := make(chan map[int32]*process.Process, 1)
			endChannel := make(chan os.Signal, 1)

			wg.Add(2)

			go process.ServeProcs(dataChannel, endChannel, &wg)
			go procGraph.AllProcVisuals(dataChannel, endChannel, &wg)
			wg.Wait()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(procCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// procCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// procCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	procCmd.Flags().Int32P("pid", "p", -1, "specify pid of process")
}
