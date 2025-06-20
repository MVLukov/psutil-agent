package metrics

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

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

type OS struct {
	PrettyName string `json:"prettyName"`
	ID         string `json:"id"`
}

type HostINFO struct {
	Hostname string `json:"hostname"`
	OS       OS     `json:"OS"`
	Platform string `json:"platform"`
	Uptime   string `json:"uptime"`
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
	CPUInfo  CPUInfo    `json:"cpu"`
	RAMInfo  MemoryINFO `json:"ram"`
}

func GetBasicMetrics() BasicMetrics {
	cpuUsage, _ := cpu.Percent(0, false)
	memUsage, _ := mem.VirtualMemory()
	cpuCores, _ := cpu.Counts(false)
	cpuThreads, _ := cpu.Counts(true)
	cpuInfo, _ := cpu.Info()
	hostInfo, _ := host.Info()
	os, err := getOsType()

	if err != nil {
		fmt.Println(err.Error())
	}

	hostInfoS := HostINFO{
		Hostname: hostInfo.Hostname,
		Platform: hostInfo.Platform,
		Uptime:   formatUptime(int(hostInfo.Uptime)),
		OS:       os,
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
		CPUInfo:  cpuInfoS,
		RAMInfo:  ram,
	}

	return metrics
}

func getLinuxDistro() (OS, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return OS{}, err
	}
	defer file.Close()

	osS := OS{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "PRETTY_NAME=") {
			// return strings.Trim(line[13:], `"`), nil
			osS.PrettyName = strings.Trim(line[13:], `"`)
		}

		if strings.HasPrefix(line, "ID=") {
			id := strings.TrimPrefix(line, "ID=")
			id = strings.Trim(id, `"`) // remove surrounding quotes

			osS.ID = id
		}
	}

	return osS, nil
}

func formatUptime(seconds int) string {
	d := time.Duration(seconds) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func getOsType() (OS, error) {
	osName := runtime.GOOS

	if osName == "windows" {
		windows, err := GetWindowsVersion()
		if err != nil {
			return OS{}, err
		}

		return OS{
			PrettyName: fmt.Sprintf("%s %s (Build %d)", windows[0], windows[1]),
			ID:         strings.ToLower(windows[0]),
		}, nil
	}

	if osName == "linux" {
		osS, err := getLinuxDistro()
		if err != nil {
			fmt.Println(err)
			return osS, err
		}

		return osS, nil
	}

	if osName == "darwin" {
		return OS{
			PrettyName: "unknown",
			ID:         "unknown",
		}, nil
	}

	return OS{}, nil
}
