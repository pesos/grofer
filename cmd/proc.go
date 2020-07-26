/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

	"github.com/pesos/grofer/src/graphs"
	"github.com/pesos/grofer/src/process"
	"github.com/spf13/cobra"
)

// procCmd represents the proc command
var procCmd = &cobra.Command{
	Use:   "proc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("the proc command should have no arguments, see grofer proc --help for further info")
		}

		pid, _ := cmd.Flags().GetInt32("pid")

		var wg sync.WaitGroup
		endChannel := make(chan os.Signal, 1)
		dataChannel := make(chan *process.Process, 1)

		wg.Add(2)

		procs, err := process.InitAllProcs()
		if err != nil {
			return err
		}
		go process.Serve(procs, pid, dataChannel, endChannel, &wg)
		//time.Sleep(2 * time.Second)
		go graphs.ProcVisuals(endChannel, dataChannel, &wg)

		wg.Wait()
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
	procCmd.Flags().Int32P("pid", "p", 1, "specify pid of process")
}
