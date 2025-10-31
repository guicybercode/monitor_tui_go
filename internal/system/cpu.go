package system

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUMetrics struct {
	Usage   float64
	PerCore []float64
	Cores   int
	Model   string
	LoadAvg []float64
}

func GetCPUMetrics() (*CPUMetrics, error) {
	ctx := context.Background()

	percent, err := cpu.PercentWithContext(ctx, time.Second, true)
	if err != nil {
		return nil, err
	}

	avg, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		return nil, err
	}

	info, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	count, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		return nil, err
	}

	model := "Unknown"
	if len(info) > 0 {
		model = info[0].ModelName
	}

	loadAvg := []float64{}
	if len(avg) > 0 {
		loadAvg = []float64{avg[0]}
	}

	perCore := make([]float64, len(percent))
	copy(perCore, percent)

	var totalUsage float64
	for _, p := range percent {
		totalUsage += p
	}
	if len(percent) > 0 {
		totalUsage /= float64(len(percent))
	}

	return &CPUMetrics{
		Usage:   totalUsage,
		PerCore: perCore,
		Cores:   count,
		Model:   model,
		LoadAvg: loadAvg,
	}, nil
}
