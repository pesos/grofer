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
	"strings"

	export "github.com/pesos/grofer/src/export/general"
	"github.com/spf13/cobra"
)

const (
	DefaultExportRefreshRate = 1000
	DefaultExportIterations  = 10
	DefaultExportFileName    = "grofer_profile"
	DefaultExportType        = "json"
)

func validateFileName(fileName string, exportType string) error {
	split := strings.Split(fileName, ".")
	if split[len(split)-1] != "json" && split[len(split)-1] != "csv" {
		return fmt.Errorf("invalid file extension")
	}
	if split[len(split)-1] != exportType {
		return fmt.Errorf("mismatch of export type and file extension")
	}

	return nil
}

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Used to export profiled data.",
	Long:  `the export command can be used to export profiled data to a specific file format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iter, err := cmd.Flags().GetUint32("iter")
		if err != nil {
			return err
		}

		refreshRate, err := cmd.Flags().GetUint64("refresh")
		if err != nil {
			return err
		}
		exportType, err := cmd.Flags().GetString("type")
		exportType = strings.ToLower(exportType)
		if err != nil {
			return err
		}

		fileName, err := cmd.Flags().GetString("fileName")
		if err != nil {
			return err
		}
		if fileName == DefaultExportFileName {
			fileName = fileName + "." + exportType
		}
		err = validateFileName(fileName, exportType)
		if err != nil {
			return err
		}

		switch exportType {
		case "json":
			return export.ExportJSON(fileName, iter, refreshRate)
		// TODO: add csv export functionality
		default:
			return fmt.Errorf("invalid export type, see grofer export --help")
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// add flags for the export command
	exportCmd.Flags().Uint32P(
		"iter",
		"i",
		DefaultExportIterations,
		"specify the number of iterations to run profiler",
	)
	exportCmd.Flags().StringP(
		"fileName",
		"f",
		DefaultExportFileName,
		"specify the name of the export file",
	)
	exportCmd.Flags().Uint64P(
		"refresh",
		"r",
		DefaultExportRefreshRate,
		"specify frequency of data fetch in milliseconds",
	)
	exportCmd.Flags().StringP(
		"type",
		"t",
		DefaultExportType,
		"specify the output format of the profiling result (json or csv)",
	)
}
