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
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
)

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "about is a command that gives information about the project in a cute way",
	Run: func(cmd *cobra.Command, args []string) {

		if err := ui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		defer ui.Close()

		About := widgets.NewParagraph()
		About.Title = " Grofer "
		About.TitleStyle.Fg = ui.ColorCyan
		About.Border = true
		About.BorderStyle.Fg = ui.ColorBlue
		About.Text =
			"\nA system profiler written purely in golang!\n\n" +
				"Made with [♥](fg:red) by [PES Open Source](fg:green)\n\n"

		uiEvents := ui.PollEvents()
		t := time.NewTicker(100 * time.Millisecond)
		tick := t.C

		for {
			select {
			case e := <-uiEvents: // For keyboard events
				switch e.ID {
				case "q", "<C-c>": // q or Ctrl-C to quit
					return
				}
			case <-tick:
				ui.Clear()
				w, h := ui.TerminalDimensions()
				About.SetRect((w-35)/2, (h-10)/2, (w+35)/2, (h+10)/2)
				ui.Render(About)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}
