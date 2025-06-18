package metrics

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func FormatBytes(bytes uint64) string {
	const (
		KB = 1 << (10 * 1)
		MB = 1 << (10 * 2)
		GB = 1 << (10 * 3)
		TB = 1 << (10 * 4)
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

type HostINFO struct {
	Hostname string `json:"hostname"`
	Platform string `json:"platform"`
	Uptime   int    `json:"uptime"`
}

type MemoryINFO struct {
	TotalMem     string `json:"total"`
	AvailableMem string `json:"available"`
	UsedMem      string `json:"used"`
	FreeMem      string `json:"free"`
	SwapTotal    string `json:"swapTotal"`
	SwapFree     string `json:"swapFree"`
	SwapUsed     string `json:"swapUsed"`
}

type CPUInfo struct {
	Vendor    string  `json:"vendor"`
	ModelName string  `json:"modelName"`
	Cores     int     `json:"cores"`
	Threads   int     `json:"threads"`
	Usage     float64 `json:"usage"`
}

type BasicMetrics struct {
	HostINFO HostINFO   `json:"host"`
	CPUUsage CPUInfo    `json:"cpu"`
	RAMUsage MemoryINFO `json:"ram"`
}

func GetBasicMetrics() BasicMetrics {
	cpuUsage, _ := cpu.Percent(0, false)
	memUsage, _ := mem.VirtualMemory()
	cpuCores, _ := cpu.Counts(false)
	cpuThreads, _ := cpu.Counts(true)
	cpuInfo, _ := cpu.Info()
	hostInfo, _ := host.Info()

	hostInfoS := HostINFO{
		Hostname: hostInfo.Hostname,
		Platform: hostInfo.Platform,
		Uptime:   int(hostInfo.Uptime),
	}

	cpuInfoS := CPUInfo{
		Cores:   cpuCores,
		Threads: cpuThreads,
		Usage:   float64(cpuUsage[0]),
	}

	for _, info := range cpuInfo {
		cpuInfoS.ModelName = info.ModelName
		cpuInfoS.Vendor = info.VendorID
	}

	ram := MemoryINFO{
		TotalMem:     FormatBytes(memUsage.Total),
		UsedMem:      FormatBytes(memUsage.Used),
		FreeMem:      FormatBytes(memUsage.Free),
		AvailableMem: FormatBytes(memUsage.Available),
		SwapTotal:    FormatBytes(memUsage.SwapTotal),
		SwapUsed:     FormatBytes(memUsage.SwapCached),
		SwapFree:     FormatBytes(memUsage.SwapFree),
	}

	metrics := BasicMetrics{
		HostINFO: hostInfoS,
		CPUUsage: cpuInfoS,
		RAMUsage: ram,
	}

	return metrics
}
