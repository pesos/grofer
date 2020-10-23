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

package process

import (
	"fmt"
	"os/exec"
	"strings"
)

type extendedInfo struct {
	SchedPolicy string `json:"sched_policy,omitempty"`
}

var policyMap map[string]string = map[string]string{
	"TS":  "SCHED_OTHER",
	"RR":  "SCHED_RR",
	"FF":  "SCHED_FIFO",
	"-":   "not reported",
	"B":   "SCHED_BATCH",
	"ISO": "SCHED_ISO",
	"IDL": "SCHED_IDL",
	"DLN": "SCHED_DEADLINE",
	"?":   "unknown value",
}

func getPolicy(pid int32) (string, error) {
	psOutput, _ := exec.Command("ps", "-o", "cls=", "-p", "128").Output()
	policy := strings.Trim(string(psOutput), " ")

	policy, exists := policyMap[policy]
	if !exists {
		return "", fmt.Errorf("value for scheduling policy for PID %d not recognized", pid)
	}

	return policyMap[policy], nil
}

func (extended *extendedInfo) updateExtendedInfo(pid int32) error {
	policy, err := getPolicy(pid)
	if err != nil {
		return err
	}

	extended.SchedPolicy = policy

	return nil
}
