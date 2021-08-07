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
	"errors"
	"fmt"

	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/metrics/factory"
	"github.com/spf13/cobra"
)

const (
	defaultProcRefreshRate = 3000
	defaultProcPid         = ""
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
		// validate args and extract flags.
		procCmd, err := constructProcCommand(cmd, args)
		if err != nil {
			return err
		}

		// create a metric scraper factory that will help construct
		// a process metric specific MetricScraper.
		metricScraperFactory := factory.
			NewMetricScraperFactory().
			ForCommand(core.ProcCommand).
			WithScrapeInterval(procCmd.refreshRate)

		if procCmd.isPerProcess() {
			metricScraperFactory = metricScraperFactory.ForSingularEntity(procCmd.pid)
		}

		processMetricScraper, err := metricScraperFactory.Construct()
		if err != nil {
			return err
		}

		err = processMetricScraper.Serve()
		if err != nil && err != core.ErrCanceledByUser {
			fmt.Printf("Error: %v\n", err)
		}

		return nil
	},
}

type procCommand struct {
	pid         string
	refreshRate uint64
}

func constructProcCommand(cmd *cobra.Command, args []string) (*procCommand, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("the proc command should have no arguments, see grofer proc --help for further info")
	}

	pid, err := cmd.Flags().GetString("pid")
	if err != nil {
		return nil, errors.New("error extracting --pid flag")
	}
	procRefreshRate, err := cmd.Flags().GetUint64("refresh")
	if err != nil {
		return nil, errors.New("error extracting --refresh flag")
	}
	if procRefreshRate < 1000 {
		return nil, fmt.Errorf("invalid refresh rate: minimum refresh rate is 1000(ms)")
	}

	return &procCommand{
		refreshRate: procRefreshRate,
		pid:         pid,
	}, nil
}

func (pc *procCommand) isPerProcess() bool {
	return pc.pid != defaultProcPid
}

func init() {
	rootCmd.AddCommand(procCmd)

	procCmd.Flags().Uint64P(
		"refresh",
		"r",
		defaultProcRefreshRate,
		"Process information UI refreshes rate in milliseconds greater than 1000",
	)

	procCmd.Flags().StringP(
		"pid",
		"p",
		defaultProcPid,
		"specify PID of process. Passing PID 0 lists all the processes (same as not using the -p flag).",
	)
}
