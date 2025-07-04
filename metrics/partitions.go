package metrics

import (
	"strconv"

	"github.com/shirou/gopsutil/v4/disk"
)

type diskPartition struct {
	MountPoint  string `json:"mountPoint"`
	Used        string `json:"used"`
	Free        string `json:"free"`
	Total       string `json:"total"`
	UsedPercent string `json:"usedPercent"`
	FsType      string `json:"fsType"`
}

type PartitionMetrics struct {
	Partitions []diskPartition `json:"partitions"`
}

func GetDisksMetrics() PartitionMetrics {
	diskPartitions := []diskPartition{}

	partitions, _ := disk.Partitions(false)
	for _, p := range partitions {
		path := p.Mountpoint
		if path == "/etc/hosts" || path == "/etc/hostname" || path == "/etc/resolv.conf" {
			continue
		}

		part, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		diskPartitions = append(diskPartitions, diskPartition{
			MountPoint:  p.Mountpoint,
			Used:        FormatBytes(part.Used),
			Free:        FormatBytes(part.Free),
			Total:       FormatBytes(part.Total),
			UsedPercent: strconv.Itoa(int(part.UsedPercent)),
			FsType:      part.Fstype,
		})
	}

	metrics := PartitionMetrics{
		Partitions: diskPartitions,
	}

	return metrics
}
