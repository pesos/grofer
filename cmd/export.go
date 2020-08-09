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
	export "github.com/pesos/grofer/src/export/general"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Used to export profiled data.",
	Long:  `the export command can be used to export profiled data to a specific file format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var file string = "test.json"
		var iter int64 = 5
		var interval uint64 = 10
		return export.ExportJSON(file, iter, interval)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	// TODO: Add flags Issue #34
}
