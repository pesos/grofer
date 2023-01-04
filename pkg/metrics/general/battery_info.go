package general

import (
	"bufio"
	"context"
	"os"
	"reflect"
	"strings"

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
	Battery [][]string
}

// NewBatteryInfo is a constructor for the BatteryInfo type.
func NewBatteryInfo() *BatteryInfo {
	return &BatteryInfo{}
}

// GetBatteryInfo updates the BatteryInfo struct and serves the data to the data channel.
func GetBatteryInfo(ctx context.Context, batteryInfo *BatteryInfo, dataChannel chan BatteryData, refreshRate uint64) error {
	return utils.TickUntilDone(ctx, refreshRate, func() error {
		err := batteryInfo.UpdateBatteryInfo()
		batteryData := batteryInfo.getBatteryData()
		if err != nil {
			return err
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
	return nil
}

// getBatteryData structures all the battery stats into the Battery data struct.
func (c *BatteryInfo) getBatteryData() BatteryData {
	var bData BatteryData = BatteryData{
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
