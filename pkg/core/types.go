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

package core

// Command represents a command or a sub-command of the `grofer` CLI.
type Command int

const (
	// RootCommand is the root command of grofer, i.e.
	// `grofer`.
	RootCommand Command = iota
	// ProcCommand is `grofer proc` and its variants.
	ProcCommand
	// ContainerCommand is `grofer container` and its
	// variants.
	ContainerCommand
	// ExportCommand is `grofer export` and its variants.
	ExportCommand
)

// Sink represents any entity that consumes generated metrics.
type Sink int

// Different Sinks here can be added depending on how `grofer` is
// extended. Sinks exist mainly to decouple entities that produce
// metrics, and entities that consume them, allowing for addition
// of Sinks independent of the metric producing entity.
const (
	// TUI reprents the terminal UI that consumes the metrics
	// generated.
	TUI Sink = iota
)
