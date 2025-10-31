package exports

import (
	"fmt"
	"os"
	"time"
)

func ExportMarkdown(filename string, report *SystemReport) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# System Report\n\n")
	fmt.Fprintf(file, "Generated: %s\n\n", report.Timestamp.Format(time.RFC3339))

	fmt.Fprintf(file, "## CPU Metrics\n\n")
	fmt.Fprintf(file, "- **Usage**: %.2f%%\n", report.CPU.Usage)
	fmt.Fprintf(file, "- **Cores**: %d\n", report.CPU.Cores)
	fmt.Fprintf(file, "- **Model**: %s\n\n", report.CPU.Model)

	fmt.Fprintf(file, "## Memory Metrics\n\n")
	fmt.Fprintf(file, "- **Total**: %s\n", formatBytes(report.Memory.Total))
	fmt.Fprintf(file, "- **Used**: %s (%.2f%%)\n", formatBytes(report.Memory.Used), report.Memory.UsedPercent)
	fmt.Fprintf(file, "- **Available**: %s\n", formatBytes(report.Memory.Available))
	fmt.Fprintf(file, "- **Free**: %s\n\n", formatBytes(report.Memory.Free))

	fmt.Fprintf(file, "## Disk Metrics\n\n")
	fmt.Fprintf(file, "- **Total**: %s\n", formatBytes(report.Disk.Total))
	fmt.Fprintf(file, "- **Used**: %s (%.2f%%)\n", formatBytes(report.Disk.Used), report.Disk.UsedPercent)
	fmt.Fprintf(file, "- **Free**: %s\n\n", formatBytes(report.Disk.Free))

	fmt.Fprintf(file, "## Processes\n\n")
	fmt.Fprintf(file, "| PID | Name | CPU %% | Memory %% | User |\n")
	fmt.Fprintf(file, "|-----|------|-------|----------|------|\n")
	for _, proc := range report.Processes {
		if len(proc.Name) > 20 {
			proc.Name = proc.Name[:20] + "..."
		}
		fmt.Fprintf(file, "| %d | %s | %.2f | %.2f | %s |\n",
			proc.PID, proc.Name, proc.CPUPercent, proc.MemPercent, proc.User)
	}
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "## Services\n\n")
	fmt.Fprintf(file, "| Name | State | Description |\n")
	fmt.Fprintf(file, "|------|-------|-------------|\n")
	for _, svc := range report.Services {
		if len(svc.Name) > 30 {
			svc.Name = svc.Name[:30] + "..."
		}
		fmt.Fprintf(file, "| %s | %s | %s |\n", svc.Name, svc.State, svc.Description)
	}
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "## Network Interfaces\n\n")
	fmt.Fprintf(file, "| Interface | Bytes Sent | Bytes Received |\n")
	fmt.Fprintf(file, "|-----------|-----------|----------------|\n")
	for _, net := range report.Network {
		fmt.Fprintf(file, "| %s | %s | %s |\n",
			net.Interface, formatBytes(net.BytesSent), formatBytes(net.BytesRecv))
	}

	return nil
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
