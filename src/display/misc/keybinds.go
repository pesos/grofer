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

package help

var PROC_KEYBINDS = []string{
	"Quit: q or <C-c>",
	"",
	"[Process navigation](fg:white)",
	"  - k and <Up>: up",
	"  - j and <Down>: down",
	"  - <C-u>: half page up",
	"  - <C-d>: half page down",
	"  - <C-b>: full page up",
	"  - <C-f>: full page down",
	"  - gg and <Home>: jump to top",
	"  - G and <End>: jump to bottom",
	"",
	"[Sorting](fg:white)",
	"  - Use column number to sort ascending.",
	"  - Use <F-column number> to sort descending.",
	"  - Eg: 1 to sort ascedning on 1st Col and F1 for descending",
	"  - 0: Disable Sort",
	"",
	"[Process actions](fg:white)",
	"  - K and <F9>: Kill process",
}

var CONT_KEYBINDS = []string{
	"Quit: q or <C-c>",
	"",
	"[Container navigation](fg:white)",
	"  - k and <Up>: up",
	"  - j and <Down>: down",
	"  - <C-u>: half page up",
	"  - <C-d>: half page down",
	"  - <C-b>: full page up",
	"  - <C-f>: full page down",
	"  - gg and <Home>: jump to top",
	"  - G and <End>: jump to bottom",
	"",
	"[Sorting](fg:white)",
	"  - Use column number to sort ascending.",
	"  - Use <F-column number> to sort descending.",
	"  - Eg: 1 to sort ascedning on 1st Col and F1 for descending",
	"  - 0: Disable Sort",
	// "",
	// "[Container actions](fg:white)",
	// "  - K and <F9>: Kill Container",
}

var CONT_CID_KEYBINDS = []string{
	"Quit: q or <C-c>",
	"",
	"[Table Selection](fg:white)",
	"  - 1: MountTable",
	"  - 2: NetworkTable",
	"  - 3: CPUUsageTable",
	"  - 4: PortMapTable",
	"  - 5: ProcTable",
	"",
	"[Table navigation](fg:white)",
	"  - k and <Up>: up",
	"  - j and <Down>: down",
	"  - <C-u>: half page up",
	"  - <C-d>: half page down",
	"  - <C-b>: full page up",
	"  - <C-f>: full page down",
	"  - gg and <Home>: jump to top",
	"  - G and <End>: jump to bottom",
}

var MAIN_KEYBINDS = []string{
	"Quit: q or <C-c>",
}

var PROC_PID_KEYBINDS = []string{
	"Quit: q or <C-c>",
}
