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

	exportGeneral "github.com/pesos/grofer/src/export/general"
	exportProc "github.com/pesos/grofer/src/export/proc"

	"github.com/spf13/cobra"
)

const (
	defaultExportRefreshRate = 1000
	defaultExportIterations  = 10
	defaultExportFileName    = "grofer_profile"
	defaultExportType        = "json"
	defaultExportPid         = -1
)

// Maintain a map of extensions provided by grofer.
// If grofer were to support a config file which
// enforces certain types to be explicitly disabled then
// a map would prove useful.
var providedExportTypes = map[string]bool{
	"json": true,
}

func hasValidExtension(filename, exportType string) error {
	filename = strings.ToLower(filename)

	var hasProvidedExtension bool = false

	// Check if any one of the allowed export types is a suffix for the
	// file name provided.
	for exportType, allowed := range providedExportTypes {
		if allowed {
			hasType := strings.HasSuffix(filename, "."+exportType)
			hasProvidedExtension = hasProvidedExtension || hasType
		}
	}
	// If en extension which is supported by grofer is provided
	// then check if it matches with the export type specified
	// in the command. If not then return an error
	if hasProvidedExtension {
		validExtension := strings.HasSuffix(filename, exportType)
		if validExtension {
			return nil
		}
		return fmt.Errorf("invaid file extension")
	}

	// If the file extension is something that grofer does not recognise
	// then it assumes that it is a valid type and trusts the user on the sme.
	return nil
}

func validateFileName(filename, exportType string) error {
	isValid := hasValidExtension(filename, exportType)
	return isValid
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
		if err != nil {
			return err
		}

		exportType = strings.ToLower(exportType)
		if validExportType := providedExportTypes[exportType]; !validExportType {
			return fmt.Errorf("export type not supported")
		}

		exportPid, err := cmd.Flags().GetInt32("pid")
		if err != nil {
			return err
		}

		filename, err := cmd.Flags().GetString("filename")
		if err != nil {
			return err
		}

		err = validateFileName(filename, exportType)
		if err != nil {
			return err
		}

		if exportPid == defaultExportPid {
			switch exportType {
			case "json":
				return exportGeneral.ExportJSON(filename, iter, refreshRate)

			default:
				return fmt.Errorf("invalid export type, see grofer export --help")
			}
		} else {
			switch exportType {
			case "json":
				return exportProc.ExportPidJSON(exportPid, filename, iter, refreshRate)

			default:
				return fmt.Errorf("invalid export type, see grofer export --help")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// add flags for the export command
	exportCmd.Flags().Uint32P(
		"iter",
		"i",
		defaultExportIterations,
		"specify the number of iterations to run profiler",
	)
	exportCmd.Flags().StringP(
		"filename",
		"f",
		defaultExportFileName,
		"specify the name of the export file",
	)
	exportCmd.Flags().Uint64P(
		"refresh",
		"r",
		defaultExportRefreshRate,
		"specify frequency of data fetch in milliseconds",
	)
	exportCmd.Flags().StringP(
		"type",
		"t",
		defaultExportType,
		"specify the output format of the profiling result (json or csv)",
	)
	exportCmd.Flags().Int32P(
		"pid",
		"p",
		defaultExportPid,
		"specify pid of process to profile, ignore to profile all processes",
	)
}
