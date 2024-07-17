//go:build darwin
// +build darwin

package hstat

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/viper"
	"n2x.dev/x-api-go/grpc/resources/topology"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/hstat/alert"
)

func (hs *hstats) updateSys(nr *topology.NodeReq) {
	nodeName := viper.GetString("nodeName")

	hostInfo, err := host.Info()
	if err != nil {
		xlog.Warnf("[hstats] Unable to get host info: %v", err)
		return
	}
	cpuCount, err := cpu.Counts(true)
	if err != nil {
		xlog.Warnf("[hstats] Unable to get cpu counts: %v", err)
		return
	}
	// cpuPercent, err := cpu.Percent(0, false)
	// if err != nil {
	// 	xlog.Warnf("[hstats] Unable to get cpu info: %v", err)
	// 	return
	// }
	loadAvg, err := load.Avg()
	if err != nil {
		xlog.Warnf("[hstats] Unable to get load info: %v", err)
		return
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		xlog.Warnf("[hstats] Unable to get mem info: %v", err)
		return
	}
	diskInfo, err := disk.Usage("/")
	if err != nil {
		xlog.Warnf("[hstats] Unable to get disk info: %v", err)
		return
	}

	// hs.host.OS = runtime.GOOS
	// hs.host.OS = hostInfo.OS
	hs.host.OS = fmt.Sprintf("macOS %s", hostInfo.PlatformVersion)
	hs.host.OSType = getOSType()
	hs.host.Arch = runtime.GOARCH

	if hostInfo.Uptime > 0 {
		hs.host.Uptime = uptimeStr(hostInfo.Uptime)
		if hostInfo.Uptime < 300 { // uptime < 300 seconds
			go alert.HostUptimeAlert(nr, nodeName, hs.host.Uptime)
		}
	}

	hs.host.LoadAvg = loadAvg.Load5
	if (loadAvg.Load5 / float64(cpuCount)) > 2.00 {
		hs.host.CpuPressure = true
		go alert.HostCPUHighAlert(nr, nodeName, fmt.Sprintf("%f", loadAvg.Load5/float64(cpuCount)))
	} else {
		hs.host.CpuPressure = false
		go alert.HostCPULowAlert(nr, nodeName, fmt.Sprintf("%f", loadAvg.Load5/float64(cpuCount)))
	}

	// hs.host.CpuUsage = uint64(cpuPercent[0])
	// if hs.host.CpuUsage > 90 {
	// 	hs.host.CpuPressure = true
	// 	go alert.HostCPUHighAlert(nr, nodeName, fmt.Sprintf("%d%%", hs.host.CpuUsage))
	// } else {
	// 	hs.host.CpuPressure = false
	// 	go alert.HostCPULowAlert(nr, nodeName, fmt.Sprintf("%d%%", hs.host.CpuUsage))
	// }

	hs.host.MemoryUsage = uint64(memInfo.UsedPercent)
	if hs.host.MemoryUsage > 90 {
		hs.host.MemoryPressure = true
		go alert.HostMemHighAlert(nr, nodeName, fmt.Sprintf("%d%%", hs.host.MemoryUsage))
	} else {
		hs.host.MemoryPressure = false
		go alert.HostMemLowAlert(nr, nodeName, fmt.Sprintf("%d%%", hs.host.MemoryUsage))
	}

	hs.host.DiskUsage = uint64(diskInfo.UsedPercent)
	if hs.host.DiskUsage > 90 {
		hs.host.DiskPressure = true
		go alert.HostDiskHighAlert(nr, nodeName, fmt.Sprintf("%d%%", hs.host.DiskUsage))
	} else {
		hs.host.DiskPressure = false
		go alert.HostDiskLowAlert(nr, nodeName, fmt.Sprintf("%d%%", hs.host.DiskUsage))
	}
}
