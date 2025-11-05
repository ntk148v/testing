
package main

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	fmt.Println("========== System Information ==========")

	// --- OS / Host Info ---
	hostInfo, _ := host.Info()
	fmt.Printf("OS: %s %s (%s)\n", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.KernelArch)
	fmt.Printf("Hostname: %s\n", hostInfo.Hostname)
	fmt.Printf("Uptime: %d seconds\n", hostInfo.Uptime)

	// --- CPU Info ---
	cpuInfo, _ := cpu.Info()
	fmt.Printf("\nCPU: %s\n", cpuInfo[0].ModelName)
	fmt.Printf("Cores: %d\n", runtime.NumCPU())

	// --- Memory Info ---
	vmStat, _ := mem.VirtualMemory()
	fmt.Printf("\nTotal RAM: %.2f GB\n", float64(vmStat.Total)/(1024*1024*1024))
	fmt.Printf("Used RAM:  %.2f GB (%.2f%%)\n",
		float64(vmStat.Used)/(1024*1024*1024),
		vmStat.UsedPercent)

	// --- Disk Info ---
	diskStat, _ := disk.Usage("/")
	fmt.Printf("\nDisk Total: %.2f GB\n", float64(diskStat.Total)/(1024*1024*1024))
	fmt.Printf("Disk Used:  %.2f GB (%.2f%%)\n",
		float64(diskStat.Used)/(1024*1024*1024),
		diskStat.UsedPercent)
}

