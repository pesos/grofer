/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	containerGraph "github.com/pesos/grofer/src/display/container"

	"github.com/pesos/grofer/src/container"
	"github.com/pesos/grofer/src/general"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	defaultCid                  = ""
	defaultContainerRefreshRate = 3000
)

// containerCmd represents the container command
var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"containers"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("the proc command should have no arguments, see grofer proc --help for further info")
		}

		cid, _ := cmd.Flags().GetString("container-id")
		containerRefreshRate, _ := cmd.Flags().GetUint64("refresh")

		if containerRefreshRate < 1000 {
			return fmt.Errorf("invalid refresh rate: minimum refresh rate is 1000(ms)")
		}

		if cid != defaultCid {

		} else {
			dataChannel := make(chan container.ContainerMetrics, 1)

			eg, ctx := errgroup.WithContext(context.Background())

			eg.Go(func() error {
				return container.Serve(dataChannel, ctx, int64(containerRefreshRate))
			})
			eg.Go(func() error {
				return containerGraph.OverallVisuals(ctx, dataChannel, containerRefreshRate)
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
	rootCmd.AddCommand(containerCmd)

	containerCmd.Flags().StringP(
		"container-id",
		"c",
		"",
		"specify container ID",
	)

	containerCmd.Flags().Uint64P(
		"refresh",
		"r",
		defaultContainerRefreshRate,
		"Container information UI refreshes rate in milliseconds greater than 1000",
	)
}
