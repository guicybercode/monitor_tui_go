package system

import (
	"context"

	"github.com/shirou/gopsutil/v3/disk"
)

type DiskMetrics struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

type DiskInfo struct {
	Device      string
	Mountpoint  string
	Fstype      string
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func GetDiskMetrics() (*DiskMetrics, error) {
	ctx := context.Background()

	usage, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return nil, err
	}

	return &DiskMetrics{
		Total:       usage.Total,
		Used:        usage.Used,
		Free:        usage.Free,
		UsedPercent: usage.UsedPercent,
	}, nil
}

func GetDiskPartitions() ([]DiskInfo, error) {
	ctx := context.Background()

	partitions, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		return nil, err
	}

	var diskInfos []DiskInfo
	for _, part := range partitions {
		usage, err := disk.UsageWithContext(ctx, part.Mountpoint)
		if err != nil {
			continue
		}

		diskInfos = append(diskInfos, DiskInfo{
			Device:      part.Device,
			Mountpoint:  part.Mountpoint,
			Fstype:      part.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	return diskInfos, nil
}
