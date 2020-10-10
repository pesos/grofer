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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	overallGraph "github.com/pesos/grofer/src/display/general"
	"github.com/pesos/grofer/src/general"
	info "github.com/pesos/grofer/src/general"
	"github.com/pesos/grofer/src/utils"
	"golang.org/x/sync/errgroup"
)

const (
	defaultOverallRefreshRate = 1000
	defaultConfigFileLocation = ""
	defaultCPUBehavior        = false
	defaultHelpMessageToggle  = false
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grofer",
	Short: "grofer is a system and resource monitor written in golang",
	Long: `grofer is a system and resource monitor written in golang.

While using a TUI based command, press ? to get information about keybindigns (if any) for that command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		overallRefreshRate, _ := cmd.Flags().GetUint64("refresh")
		if overallRefreshRate < 1000 {
			return fmt.Errorf("invalid refresh rate: minimum refresh rate is 1000(ms)")
		}

		cpuLoadFlag, _ := cmd.Flags().GetBool("cpuinfo")
		if cpuLoadFlag {
			cpuLoad := info.NewCPULoad()
			dataChannel := make(chan *info.CPULoad, 1)

			eg, ctx := errgroup.WithContext(context.Background())

			eg.Go(func() error {
				return info.GetCPULoad(ctx, cpuLoad, dataChannel, uint64(4*overallRefreshRate/5))
			})

			eg.Go(func() error {
				return overallGraph.RenderCPUinfo(ctx, dataChannel, overallRefreshRate)
			})

			if err := eg.Wait(); err != nil {
				if err != general.ErrCanceledByUser {
					fmt.Printf("Error: %v\n", err)
				}
			}

		} else {
			dataChannel := make(chan utils.DataStats, 1)

			eg, ctx := errgroup.WithContext(context.Background())

			eg.Go(func() error {
				return general.GlobalStats(ctx, dataChannel, uint64(4*overallRefreshRate/5))
			})
			eg.Go(func() error {
				return overallGraph.RenderCharts(ctx, dataChannel, overallRefreshRate)
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		defaultConfigFileLocation,
		"config file (default is $HOME/.grofer.yaml)",
	)

	rootCmd.Flags().Uint64P(
		"refresh",
		"r",
		defaultOverallRefreshRate,
		"Overall stats UI refreshes rate in milliseconds greater than 1000",
	)

	rootCmd.Flags().BoolP(
		"cpuinfo",
		"c",
		defaultCPUBehavior,
		"Info about the CPU Load over all CPUs",
	)

	rootCmd.Flags().BoolP(
		"toggle",
		"t",
		defaultHelpMessageToggle,
		"Help message for toggle",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".grofer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".grofer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
