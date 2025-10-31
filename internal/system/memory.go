package system

import (
	"context"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryMetrics struct {
	Total       uint64
	Used        uint64
	Available   uint64
	Free        uint64
	UsedPercent float64
	Cached      uint64
	Buffers     uint64
}

func GetMemoryMetrics() (*MemoryMetrics, error) {
	ctx := context.Background()

	vmem, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return &MemoryMetrics{
		Total:       vmem.Total,
		Used:        vmem.Used,
		Available:   vmem.Available,
		Free:        vmem.Free,
		UsedPercent: vmem.UsedPercent,
		Cached:      vmem.Cached,
		Buffers:     vmem.Buffers,
	}, nil
}
