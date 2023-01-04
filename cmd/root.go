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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/metrics/factory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultOverallRefreshRate = 1000
	defaultConfigFileLocation = ""
	defaultCPUBehavior        = false
	defaultBatteryBehaviour   = false
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grofer",
	Short: "grofer is a system and resource monitor written in golang",
	Long: `grofer is a system and resource monitor written in golang.

While using a TUI based command, press ? to get information about key bindings (if any) for that command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// validate and extract flags.
		rootCmd, err := constructRootCommand(cmd, args)
		if err != nil {
			return err
		}

		// construct a system wide metric specific MetricScraper.
		systemWideMetricScraper, err := factory.
			NewMetricScraperFactory().
			ForCommand(core.RootCommand).
			WithScrapeInterval(rootCmd.refreshRate).
			Construct()

		if err != nil {
			return err
		}

		if !rootCmd.batteryInfo {
			err = systemWideMetricScraper.Serve(factory.WithCPUInfoAs(rootCmd.cpuInfo))
			if err != nil && err != core.ErrCanceledByUser {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			err = systemWideMetricScraper.Serve(factory.WithBatteryInfoAs(rootCmd.batteryInfo))
			if err != nil && err != core.ErrCanceledByUser {
				fmt.Printf("Error: %v\n", err)
			}
		}

		return nil
	},
}

type rootCommand struct {
	refreshRate uint64
	cpuInfo     bool
	batteryInfo bool
}

func constructRootCommand(cmd *cobra.Command, args []string) (*rootCommand, error) {
	refreshRate, err := cmd.Flags().GetUint64("refresh")
	if err != nil {
		return nil, err
	}

	if refreshRate < 1000 {
		return nil, fmt.Errorf("invalid refresh rate: minimum refresh rate is 1000(ms)")
	}

	cpuInfo, err := cmd.Flags().GetBool("cpuinfo")
	if err != nil {
		return nil, err
	}

	batteryInfo, err := cmd.Flags().GetBool("battery")
	if err != nil {
		return nil, err
	}

	return &rootCommand{
		refreshRate: refreshRate,
		cpuInfo:     cpuInfo,
		batteryInfo: batteryInfo,
	}, nil
}

// Execute runs the root command
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
		"battery",
		"b",
		defaultBatteryBehaviour,
		"All stats about the battery.",
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
