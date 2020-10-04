package export

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	procInfo "github.com/pesos/grofer/src/process"
	"github.com/pesos/grofer/src/utils"
)

type contextSwitches struct {
	Voluntary   int64 `json:"voluntary"`
	InVoluntary int64 `json:"involuntary"`
}

type pageFaults struct {
	Minor uint64 `json:"minor"`
	Major uint64 `json:"major"`
}

type memPidStats struct {
	RSS   uint64 `json:"RSS"`
	Data  uint64 `json:"Data"`
	Stack uint64 `json:"Stack"`
	Swap  uint64 `json:"Swap"`
}

type ProcDetails struct {
	Name              string `json:"Name"`
	Command           string `json:"Command"`
	Status            string `json:"Status"`
	Background        string `json:"Background"`
	Running           string `json:"Running"`
	CreationTime      string `json:"CreationTime"`
	NiceValue         string `json:"NiceValue"`
	ThreadCount       string `json:"ThreadCount"`
	ChildProcessCount string `json:"ChildProcessCount"`
}

type PidStats struct {
	Cpu         float64         `json:"cpu"`
	Mem         float32         `json:"mem"`
	PidDetails  ProcDetails     `json:"pid"`
	CtxSwitches contextSwitches `json:"ctxSwitches"`
	PageFaults  pageFaults      `json:"pageFaults"`
	MemStats    memPidStats     `json:"memStats"`
	ChildProcs  map[int]string  `json:"childProcs"`
}

var statusMap map[string]string = map[string]string{
	"R": "Running",
	"S": "Sleep",
	"Z": "Zombie",
	"T": "Stop",
	"I": "Idle",
	"W": "Wait",
	"L": "Lock",
}

func getPidDataJSON(pid int32) (PidStats, error) {
	proc, err := procInfo.NewProcess(pid)
	if err != nil {
		return PidStats{}, err
	}
	proc.UpdateProcInfo()

	pidDetails := ProcDetails{
		Name:              proc.Name,
		Command:           proc.Exe,
		Status:            statusMap[proc.Status],
		Background:        strconv.FormatBool(proc.Background),
		Running:           strconv.FormatBool(proc.IsRunning),
		CreationTime:      utils.GetDateFromUnix(proc.CreateTime),
		NiceValue:         strconv.Itoa(int(proc.Nice)),
		ThreadCount:       strconv.Itoa(int(proc.NumThreads)),
		ChildProcessCount: strconv.Itoa(len(proc.Children)),
	}

	ctxSwitches := contextSwitches{
		Voluntary:   proc.NumCtxSwitches.Voluntary,
		InVoluntary: proc.NumCtxSwitches.Involuntary,
	}

	pgFaults := pageFaults{
		Minor: proc.PageFault.MinorFaults,
		Major: proc.PageFault.MajorFaults,
	}

	memStats := memPidStats{
		RSS:   proc.MemoryInfo.RSS,
		Data:  proc.MemoryInfo.Data,
		Stack: proc.MemoryInfo.Stack,
		Swap:  proc.MemoryInfo.Swap,
	}

	pidData := PidStats{
		Cpu:         proc.CPUPercent,
		Mem:         proc.MemoryPercent,
		PidDetails:  pidDetails,
		CtxSwitches: ctxSwitches,
		PageFaults:  pgFaults,
		MemStats:    memStats,
	}
	return pidData, nil
}

func ExportPidJSON(pid int32, filename string, iter uint32, refreshRate uint64) error {
	toWrite, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer toWrite.Close()
	encoder := json.NewEncoder(toWrite)
	var jsonData []PidStats
	for i := uint32(0); i < iter; i++ {
		data, err := getPidDataJSON(pid)
		if err != nil {
			break
		}
		jsonData = append(jsonData, data)
		time.Sleep(time.Duration(refreshRate) * time.Millisecond)
	}

	err := encoder.Encode(jsonData)
	if err != nil {
		return err
	}
	return nil
}
