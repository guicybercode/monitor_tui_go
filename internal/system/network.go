package system

import (
	"context"

	"github.com/shirou/gopsutil/v3/net"
)

type NetworkStats struct {
	Interface   string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
	ErrorsIn    uint64
	ErrorsOut   uint64
	DropIn      uint64
	DropOut     uint64
}

type NetworkConnection struct {
	Status     string
	LocalAddr  string
	LocalPort  uint32
	RemoteAddr string
	RemotePort uint32
	PID        int32
}

func GetNetworkStats() ([]NetworkStats, error) {
	ctx := context.Background()

	ioCounters, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		return nil, err
	}

	var stats []NetworkStats
	for _, io := range ioCounters {
		stats = append(stats, NetworkStats{
			Interface:   io.Name,
			BytesSent:   io.BytesSent,
			BytesRecv:   io.BytesRecv,
			PacketsSent: io.PacketsSent,
			PacketsRecv: io.PacketsRecv,
			ErrorsIn:    io.Errin,
			ErrorsOut:   io.Errout,
			DropIn:      io.Dropin,
			DropOut:     io.Dropout,
		})
	}

	return stats, nil
}

func GetNetworkConnections() ([]NetworkConnection, error) {
	ctx := context.Background()

	connections, err := net.ConnectionsWithContext(ctx, "inet")
	if err != nil {
		return nil, err
	}

	var netConns []NetworkConnection
	for _, conn := range connections {
		netConns = append(netConns, NetworkConnection{
			Status:     conn.Status,
			LocalAddr:  conn.Laddr.IP,
			LocalPort:  conn.Laddr.Port,
			RemoteAddr: conn.Raddr.IP,
			RemotePort: conn.Raddr.Port,
			PID:        conn.Pid,
		})
	}

	return netConns, nil
}
