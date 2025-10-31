package exports

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/guicybercode/systui/internal/system"
)

type SystemReport struct {
	Timestamp time.Time             `json:"timestamp"`
	CPU       *system.CPUMetrics    `json:"cpu"`
	Memory    *system.MemoryMetrics `json:"memory"`
	Disk      *system.DiskMetrics   `json:"disk"`
	Processes []system.ProcessInfo  `json:"processes"`
	Services  []system.ServiceInfo  `json:"services"`
	Network   []system.NetworkStats `json:"network"`
}

func ExportJSON(filename string, report *SystemReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func GenerateReport() (*SystemReport, error) {
	cpu, err := system.GetCPUMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU metrics: %w", err)
	}

	mem, err := system.GetMemoryMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory metrics: %w", err)
	}

	disk, err := system.GetDiskMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get disk metrics: %w", err)
	}

	processes, err := system.GetProcesses()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	services, err := system.GetServices()
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	network, err := system.GetNetworkStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get network stats: %w", err)
	}

	return &SystemReport{
		Timestamp: time.Now(),
		CPU:       cpu,
		Memory:    mem,
		Disk:      disk,
		Processes: processes,
		Services:  services,
		Network:   network,
	}, nil
}
