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

package misc

// HelpKeybindingType is the type of the keybinding that
// a help page will use for a specific command.
type HelpKeybindingType int

const (
	// RootCommand is the keybinding identifier
	// for the "main" command of grofer, i.e. `grofer`.
	RootCommand HelpKeybindingType = iota
	// ProcCommand is the keybinding identifier
	// for the `grofer proc` command.
	ProcCommand
	// PerProcCommand is the keybinding identifier
	// for the `grofer proc -p <pid>` command.
	PerProcCommand
	// ContainerCommand is the keybinding identifier
	// for the `grofer container` command.
	ContainerCommand
	// PerContainerCommand is the keybinding identifier
	// for the `grofer container -c <CID>` command.
	PerContainerCommand
)

// getHelpKeybindingsForCommand returns the help keybinding for a specific command.
func getHelpKeybindingsForCommand(forCommand HelpKeybindingType) [][]string {
	switch forCommand {
	case RootCommand:
		return getMainCommandKeybindings()
	case ProcCommand:
		return getProcCommandKeybindings()
	case PerProcCommand:
		return getPerProcCommandKeybindings()
	case ContainerCommand:
		return getContainerCommandKeybindings()
	case PerContainerCommand:
		return getPerContainerCommandKeybindings()
	default:
		return getDefaultHelpKeybinding()
	}
}

func getErrorKeybindings() [][]string {
	return getDefaultHelpKeybinding()
}

func getDefaultHelpKeybinding() [][]string {
	return [][]string{
		{""},
		{"To close this prompt: <Esc>"},
	}
}

func getMainCommandKeybindings() [][]string {
	return [][]string{
		{"Quit: q or <C-c>"},
		{"Pause Rendering: s"},
		{""},
		{"Table Navigation"},
		{"  - <Left>/h: select table to left "},
		{"  - <Right>/l: select table to right"},
		{""},
		{"Table Scrolling"},
		{"  - <Up>/k: scroll up"},
		{"  - <Down>/j: scroll down"},
		{""},
		{"Enable CPU Table: t"},
		{""},
		{"To close this prompt: <Esc>"},
	}
}

func getProcCommandKeybindings() [][]string {
	return [][]string{
		{"Quit: q or <C-c>"},
		{"Pause Rendering: s"},
		{""},
		{"Process table navigation"},
		{"  - k and <Up>: scroll up"},
		{"  - j and <Down>: scroll down"},
		{"  - <C-u>: half page up"},
		{"  - <C-d>: half page down"},
		{"  - <C-b>: full page up"},
		{"  - <C-f>: full page down"},
		{"  - gg and <Home>: jump to top"},
		{"  - G and <End>: jump to bottom"},
		{""},
		{"Sorting"},
		{"  - Use column number to sort ascending."},
		{"  - Use <F-column number> to sort descending."},
		{"  - Eg: 1 to sort ascending on 1st Col and F1 for descending"},
		{"  - 0: Disable Sort"},
		{""},
		{"Process actions"},
		{"  - K and <F9>: Open signal selector menu"},
		{""},
		{"Signal selection"},
		{"  - K and <F9>: Send SIGTERM to selected process. Kills the process"},
		{"  - k and <Up>: up"},
		{"  - j and <Down>: down"},
		{"  - 0-9: navigate by numeric index"},
		{"  - <Enter>: send highlighted signal to process"},
		{"  - <Esc>: close signal selector"},
		{""},
		{"To close this prompt: <Esc>"},
	}
}

func getPerProcCommandKeybindings() [][]string {
	return [][]string{
		{"Quit: q or <C-c>"},
		{"Pause Rendering: s"},
		{""},
		{""},
		{"Table navigation"},
		{"  - k and <Up>: scroll up"},
		{"  - j and <Down>: scroll down"},
		{"  - <C-u>: half page up"},
		{"  - <C-d>: half page down"},
		{"  - <C-b>: full page up"},
		{"  - <C-f>: full page down"},
		{"  - gg and <Home>: jump to top"},
		{"  - G and <End>: jump to bottom"},
		{""},
		{"To close this prompt: <Esc>"},
	}
}

func getContainerCommandKeybindings() [][]string {
	return [][]string{
		{"Quit: q or <C-c>"},
		{"Pause Rendering: s"},
		{""},
		{"Container table navigation"},
		{"  - k and <Up>: scroll up"},
		{"  - j and <Down>: scroll down"},
		{"  - <C-u>: half page up"},
		{"  - <C-d>: half page down"},
		{"  - <C-b>: full page up"},
		{"  - <C-f>: full page down"},
		{"  - gg and <Home>: jump to top"},
		{"  - G and <End>: jump to bottom"},
		{""},
		{"Sorting"},
		{"  - Use column number to sort ascending."},
		{"  - Use <F-column number> to sort descending."},
		{"  - Eg: 1 to sort ascending on 1st Col and F1 for descending"},
		{"  - 0: Disable Sort"},
		{""},
		{"Container actions"},
		{"  - <Enter>: Open action selector menu"},
		{""},
		{"Action selection"},
		{"  - k and <Up>: up"},
		{"  - j and <Down>: down"},
		{"  - <Enter>: perform highlighted action"},
		{"  - <Esc>: close action selector"},
		{""},
		{"To close this prompt: <Esc>"},
	}
}

func getPerContainerCommandKeybindings() [][]string {
	return [][]string{
		{"Quit: q or <C-c>"},
		{"Pause Rendering: s"},
		{""},
		{"Table Selection"},
		{"  - 1: Details Table"},
		{"  - 2: Mount Table"},
		{"  - 3: Network Table"},
		{"  - 4: CPU Usage Table"},
		{"  - 5: Port Map Table"},
		{"  - 6: Proccess Table"},
		{""},
		{"Table navigation"},
		{"  - k and <Up>: scroll up"},
		{"  - j and <Down>: scroll down"},
		{"  - <C-u>: half page up"},
		{"  - <C-d>: half page down"},
		{"  - <C-b>: full page up"},
		{"  - <C-f>: full page down"},
		{"  - gg and <Home>: jump to top"},
		{"  - G and <End>: jump to bottom"},
		{""},
		{"To close this prompt: <Esc>"},
	}
}
