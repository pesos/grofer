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

package general

import (
	"bufio"
	"context"
	"os"
	"reflect"
	"strings"

	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/utils"
)

type BatteryInfo struct {
	Manufacturer         string `json:"manufacturer"`
	Technology           string `json:"technology"`
	ModelName            string `json:"model_name"`
	SerialNumber         string `json:"serial_number"`
	Capacity             string `json:"capacity"`
	CycleCount           string `json:"cycle_count"`
	EnergyFullDesign     string `json:"energy_full_design"`
	EnergyFull           string `json:"energy_full"`
	EnergyNow            string `json:"energy_now"`
	PowerNow             string `json:"power_now"`
	Status               string `json:"status"`
	ChargeStartThreshold string `json:"charge_start_threshold"`
	ChargeStopThreshold  string `json:"charge_stop_threshold"`
}

type BatteryData struct {
	FieldSet string
	Battery  [][]string
}

// NewBatteryInfo is a constructor for the BatteryInfo type.
func NewBatteryInfo() *BatteryInfo {
	return &BatteryInfo{}
}

// GetBatteryInfo updates the BatteryInfo struct and serves the data to the data channel.
func GetBatteryInfo(ctx context.Context, batteryInfo *BatteryInfo, dataChannel chan BatteryData, refreshRate uint64) error {
	return utils.TickUntilDone(ctx, refreshRate, func() error {
		var batteryFound bool = true
		var batteryData BatteryData = BatteryData{
			FieldSet: "NO BATTERY DATA",
		}
		err := batteryInfo.UpdateBatteryInfo()
		if err != nil {
			if err == core.ErrBatteryNotFound {
				batteryFound = false
			} else {
				return err
			}
		}

		if batteryFound {
			batteryData = batteryInfo.getBatteryData()
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- batteryData:
			return nil
		}
	})
}

// UpdateBatteryInfo updates fields of the type BatteryInfo
func (c *BatteryInfo) UpdateBatteryInfo() error {
	err := c.readBatteryInfo()
	if err != nil {
		return err
	}

	return nil
}

// ReadBatteryInfo reads files from /sys/class/power_supply/BAT0
// and returns battery specific stats
func (c *BatteryInfo) readBatteryInfo() error {
	_, err1 := os.Stat("/sys/class/power_supply/BAT0/manufacturer")
	_, err2 := os.Stat("/sys/class/power_supply/BAT0/technology")
	_, err3 := os.Stat("/sys/class/power_supply/BAT0/model_name")

	if err1 == nil && err2 == nil && err3 == nil {
		val := reflect.ValueOf(c).Elem()
		for i := 0; i < val.Type().NumField(); i++ {
			fileName := val.Type().Field(i).Tag.Get("json")
			file, err := os.Open("/sys/class/power_supply/BAT0/" + fileName)
			if err != nil {
				return err
			}
			defer file.Close()
			reader := bufio.NewReader(file)

			// Read first line containing load values
			data, err := reader.ReadBytes(byte('\n'))
			if err != nil {
				return err
			}
			vals := strings.Fields(string(data))
			val.FieldByName(val.Type().Field(i).Name).SetString(vals[0])
		}
	} else if os.IsNotExist(err1) || os.IsNotExist(err2) || os.IsNotExist(err3) {
		return core.ErrBatteryNotFound
	}
	return nil
}

// getBatteryData structures all the battery stats into the Battery data struct.
func (c *BatteryInfo) getBatteryData() BatteryData {
	var bData BatteryData = BatteryData{
		FieldSet: "BATTERY",
		Battery: [][]string{
			{"Stats", "Info"},
			{"manufacturer", c.Manufacturer},
			{"technology", c.Technology},
			{"model name", c.ModelName},
			{"serial number", c.SerialNumber},
			{"capacity", c.Capacity},
			{"cycle count", c.CycleCount},
			{"energy full design", c.EnergyFullDesign},
			{"energy full", c.EnergyFull},
			{"energy now", c.EnergyNow},
			{"power now", c.PowerNow},
			{"status", c.Status},
			{"charge start threshold", c.ChargeStartThreshold},
			{"charge stop threshold", c.ChargeStopThreshold},
		},
	}
	return bData
}
