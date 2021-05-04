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
	"log"

	containerGraph "github.com/pesos/grofer/src/display/container"
	"github.com/pesos/grofer/src/utils"

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
	Use:     "container",
	Short:   "container command is used to get information related to docker containers",
	Long:    `container command is used to get information related to docker containers. It provides both overall and per container metrics.`,
	Aliases: []string{"containers", "docker"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("the container command should have no arguments, see grofer container --help for further info")
		}

		cid, _ := cmd.Flags().GetString("container-id")
		containerRefreshRate, _ := cmd.Flags().GetUint64("refresh")

		if containerRefreshRate < 1000 {
			return fmt.Errorf("invalid refresh rate: minimum refresh rate is 1000(ms)")
		}

		if cid != defaultCid {
			dataChannel := make(chan container.PerContainerMetrics, 1)

			eg, ctx := errgroup.WithContext(context.Background())

			eg.Go(func() error {
				return container.ServeContainer(ctx, cid, dataChannel, int64(containerRefreshRate))
			})
			eg.Go(func() error {
				return containerGraph.ContainerVisuals(ctx, dataChannel, containerRefreshRate)
			})

			if err := eg.Wait(); err != nil {
				if err == general.ErrInvalidContainer {
					utils.ErrorMsg("cid")
				}
				if err != general.ErrCanceledByUser {
					log.Fatalf("Error: %v\n", err)
				}
			}
		} else {
			dataChannel := make(chan container.ContainerMetrics, 1)

			eg, ctx := errgroup.WithContext(context.Background())

			eg.Go(func() error {
				return container.Serve(ctx, dataChannel, int64(containerRefreshRate))
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
