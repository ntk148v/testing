package main

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

// DiskInfo represents disk partition information
type DiskInfo struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mount_point"`
	FSType      string  `json:"fs_type"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

func main() {
	fmt.Println("========== System Information ==========")

	// --- OS / Host Info ---
	hostInfo, err := host.Info()
	if err != nil || hostInfo == nil {
		fmt.Println("Failed to get host info:", err)
	} else {
		fmt.Printf("OS: %s %s (%s)\n", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.KernelArch)
		fmt.Printf("Hostname: %s\n", hostInfo.Hostname)
		fmt.Printf("Uptime: %d seconds\n", hostInfo.Uptime)
	}

	// --- CPU Info ---
	cpuInfo, err := cpu.Info()
	if err != nil || len(cpuInfo) == 0 {
		fmt.Println("Failed to get CPU info:", err)
	} else {
		fmt.Printf("\nCPU: %s\n", cpuInfo[0].ModelName)
	}
	fmt.Printf("Cores: %d\n", runtime.NumCPU())

	// --- Memory Info ---
	vmStat, err := mem.VirtualMemory()
	if err != nil || vmStat == nil {
		fmt.Println("Failed to get memory info:", err)
	} else {
		fmt.Printf("\nTotal RAM: %.2f GB\n", float64(vmStat.Total)/(1024*1024*1024))
		fmt.Printf("Used RAM:  %.2f GB (%.2f%%)\n",
			float64(vmStat.Used)/(1024*1024*1024),
			vmStat.UsedPercent)
	}

	// --- Disk Info ---
	fmt.Println("\n--- Disk Partitions ---")

	partitions, err := disk.Partitions(true)
	if err != nil {
		fmt.Println("Failed to get disk partitions:", err)
	} else {
		for _, p := range partitions {
			usage, err := disk.Usage(p.Mountpoint)
			if err != nil || usage == nil {
				fmt.Printf("Skipping %s (%s): %v\n", p.Mountpoint, p.Device, err)
				continue
			}

			diskInfo := DiskInfo{
				Device:      p.Device,
				MountPoint:  p.Mountpoint,
				FSType:      p.Fstype,
				Total:       usage.Total,
				Used:        usage.Used,
				Free:        usage.Free,
				UsedPercent: usage.UsedPercent,
			}

			fmt.Printf("Device: %s\n", diskInfo.Device)
			fmt.Printf("Mount Point: %s\n", diskInfo.MountPoint)
			fmt.Printf("FS Type: %s\n", diskInfo.FSType)
			fmt.Printf("Total: %.2f GB\n", float64(diskInfo.Total)/(1024*1024*1024))
			fmt.Printf("Used: %.2f GB (%.2f%%)\n",
				float64(diskInfo.Used)/(1024*1024*1024),
				diskInfo.UsedPercent)
			fmt.Printf("Free: %.2f GB\n", float64(diskInfo.Free)/(1024*1024*1024))
			fmt.Println()
		}
	}
}
