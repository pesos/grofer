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
	"context"
	"fmt"

	proc "github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"

	procGraph "github.com/pesos/grofer/src/display/process"
	"github.com/pesos/grofer/src/general"
	"github.com/pesos/grofer/src/process"
	"github.com/pesos/grofer/src/utils"
	"golang.org/x/sync/errgroup"
)

const (
	defaultProcRefreshRate = 3000
	defaultProcPid         = 0
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
		procRefreshRate, _ := cmd.Flags().GetUint64("refresh")

		if procRefreshRate < 1000 {
			return fmt.Errorf("invalid refresh rate: minimum refresh rate is 1000(ms)")
		}

		if pid != defaultProcPid {
			dataChannel := make(chan *process.Process, 1)

			eg, ctx := errgroup.WithContext(context.Background())

			proc, err := process.NewProcess(pid)
			if err != nil {
				utils.ErrorMsg("pid")
				return fmt.Errorf("invalid pid")
			}

			eg.Go(func() error {
				return process.Serve(proc, dataChannel, ctx, int64(4*procRefreshRate/5))
			})
			eg.Go(func() error {
				return procGraph.ProcVisuals(ctx, dataChannel, procRefreshRate)
			})

			if err := eg.Wait(); err != nil {
				if err != general.ErrCanceledByUser {
					fmt.Printf("Error: %v\n", err)
				}
			}
		} else {
			dataChannel := make(chan []*proc.Process, 1)

			eg, ctx := errgroup.WithContext(context.Background())

			eg.Go(func() error {
				return process.ServeProcs(dataChannel, ctx, int64(4*procRefreshRate/5))
			})
			eg.Go(func() error {
				return procGraph.AllProcVisuals(dataChannel, ctx, procRefreshRate)
			})

			if err := eg.Wait(); err != nil {
				if err != general.ErrCanceledByUser {
					fmt.Printf("Error: %v\n", err)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(procCmd)

	procCmd.Flags().Uint64P(
		"refresh",
		"r",
		defaultProcRefreshRate,
		"Process information UI refreshes rate in milliseconds greater than 1000",
	)

	procCmd.Flags().Int32P(
		"pid",
		"p",
		defaultProcPid,
		"specify PID of process. Passing PID 0 lists all the processes (same as not using the -p flag).",
	)
}
