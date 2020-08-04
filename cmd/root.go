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

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/pesos/grofer/src/general"
	overallGraph "github.com/pesos/grofer/src/graphs/general"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grofer",
	Short: "grofer is a system profiler written in golang",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		endChannel := make(chan os.Signal, 1)
		memChannel := make(chan []float64, 1)
		cpuChannel := make(chan []float64, 1)
		diskChannel := make(chan [][]string, 1)
		netChannel := make(chan map[string][]float64, 1)

		wg.Add(2)

		go general.GlobalStats(endChannel, cpuChannel, memChannel, diskChannel, netChannel, &wg)
		go overallGraph.RenderCharts(endChannel, memChannel, cpuChannel, diskChannel, netChannel, &wg)

		wg.Wait()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.grofer.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
