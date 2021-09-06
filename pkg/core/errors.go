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

import (
	"errors"
)

var (
	// ErrCanceledByUser is used when the UI is closed by the user
	ErrCanceledByUser = errors.New("canceled by user")
	// ErrInvalidPID is used when the user provided PID does not match a running process
	ErrInvalidPID = errors.New("PID does not exist")
	// ErrInvalidContainer is used when the user provided Container ID does not match an existing container
	ErrInvalidContainer = errors.New("container does not exist")
	// ErrBatteryNotFound is used when the host does not have a `/sys/class/power_supply/BAT0` directory tor ead battery info from
	ErrBatteryNotFound = errors.New("could not read from /sys/class/power_supply/BAT0")
)
