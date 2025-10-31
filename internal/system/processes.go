package system

import (
	"context"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	PID         int32
	Name        string
	CPUPercent  float64
	MemPercent  float32
	Status      string
	User        string
	CreateTime  int64
	CommandLine string
	Nice        int32
}

func GetProcesses() ([]ProcessInfo, error) {
	ctx := context.Background()

	pids, err := process.PidsWithContext(ctx)
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo
	for _, pid := range pids {
		proc, err := process.NewProcessWithContext(ctx, pid)
		if err != nil {
			continue
		}

		name, _ := proc.NameWithContext(ctx)
		cpuPercent, _ := proc.CPUPercentWithContext(ctx)
		memPercent, _ := proc.MemoryPercentWithContext(ctx)
		status, _ := proc.StatusWithContext(ctx)
		username, _ := proc.UsernameWithContext(ctx)
		createTime, _ := proc.CreateTimeWithContext(ctx)
		cmdline, _ := proc.CmdlineWithContext(ctx)
		nice, _ := proc.NiceWithContext(ctx)

		processes = append(processes, ProcessInfo{
			PID:         pid,
			Name:        name,
			CPUPercent:  cpuPercent,
			MemPercent:  memPercent,
			Status:      status[0],
			User:        username,
			CreateTime:  createTime,
			CommandLine: cmdline,
			Nice:        nice,
		})
	}

	return processes, nil
}

func KillProcess(pid int32, signal syscall.Signal) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}

func ReniceProcess(pid int32, nice int) error {
	cmd := exec.Command("renice", strconv.Itoa(nice), strconv.Itoa(int(pid)))
	return cmd.Run()
}
