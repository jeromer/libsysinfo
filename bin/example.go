package main

import (
	"fmt"
	"github.com/jeromer/libsysinfo"
)

func main() {
	dumpSimple()
	dumpLsbRelease()
	dumpFileSystems()
	dumpCpuInfos()
	dumpNetworkInterfaces()
	dumpMemInfos()
}

func dumpSimple() {
	type F func() (string, error)
	funs := map[string]F{
		"Hostname": libsysinfo.Hostname,
		"Domain":   libsysinfo.Domain,
		"Fqdn":     libsysinfo.Fqdn,
		"HostId":   libsysinfo.HostId,
	}

	var out string
	var err error
	for legend, f := range funs {
		out, err = f()
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("%-15s : %s\n", legend, out)
	}
}

func dumpLsbRelease() {
	lsbr, err := libsysinfo.LsbRelease()
	if err != nil {
		panic(err)
	}

	format := "- %-14s : %s\n"
	fmt.Printf("\nLsbRelease\n---------\n")
	fmt.Printf(format, "codename", lsbr.Codename)
	fmt.Printf(format, "release", lsbr.Release)
	fmt.Printf(format, "description", lsbr.Description)
	fmt.Printf(format, "distributor Id", lsbr.DistributorId)
}

func dumpFileSystems() {
	fileSystems, err := libsysinfo.FileSystems()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nFileSystems\n---------\n")
	for _, fs := range fileSystems {
		fmt.Printf("- %s\n", fs)
	}
}

func dumpCpuInfos() {
	cpuInfos, err := libsysinfo.CpuInfos()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nCpuInfos\n---------\n")

	commonFormat := "- %-14s : "
	strFormat := commonFormat + "%s\n"
	intFormat := commonFormat + "%d\n"
	floatFormat := commonFormat + "%.3f\n"

	for _, cpu := range cpuInfos {
		fmt.Printf(strFormat, "Processor", cpu.Processor)
		fmt.Printf(strFormat, "VendorId", cpu.VendorId)
		fmt.Printf(strFormat, "CpuFamily", cpu.CpuFamily)
		fmt.Printf(strFormat, "Model", cpu.Model)
		fmt.Printf(strFormat, "ModelName", cpu.ModelName)
		fmt.Printf(intFormat, "Stepping", cpu.Stepping)
		fmt.Printf(floatFormat, "CPUMHz", cpu.CPUMHz)
		fmt.Printf(intFormat, "CacheSize", cpu.CacheSize)
		fmt.Printf(strFormat, "CacheSizeUnit", cpu.CacheSizeUnit)
		fmt.Printf(strFormat, "PhysicalId", cpu.PhysicalId)
		fmt.Printf(intFormat, "Siblings", cpu.Siblings)
		fmt.Printf(strFormat, "CoreId", cpu.CoreId)
		fmt.Printf(intFormat, "CpuCores", cpu.CpuCores)
		fmt.Printf(strFormat, "ApicId", cpu.ApicId)
		fmt.Printf(strFormat, "InitialApicId", cpu.InitialApicId)
		fmt.Printf(strFormat, "Fpu", cpu.Fpu)
		fmt.Printf(strFormat, "FpuException", cpu.FpuException)
		fmt.Printf(intFormat, "CpuIdLevel", cpu.CpuIdLevel)
		fmt.Printf(strFormat, "Wp", cpu.Wp)
		fmt.Printf(strFormat, "Flags", cpu.Flags)
		fmt.Printf(floatFormat, "Bogomips", cpu.Bogomips)
		fmt.Printf(intFormat, "ClflushSize", cpu.ClflushSize)
		fmt.Printf(intFormat, "CacheAlignment", cpu.CacheAlignment)
		fmt.Printf(strFormat, "AddressSizes", cpu.AddressSizes)
	}
}

func dumpNetworkInterfaces() {
	ifaces, err := libsysinfo.NetworkInterfaces()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("\nNetwork interfaces\n------------------\n")
	format := "- %-14s : %s\n"

	for _, nif := range ifaces {
		fmt.Printf(format, "Name", nif.Name)
		fmt.Printf(format, "V4Addr", nif.V4Addr)
		fmt.Printf(format, "V6Addr", nif.V6Addr)
		fmt.Printf(format, "MacAddr", nif.MacAddr)
		fmt.Printf(format, "BroadcastAddr", nif.BroadcastAddr)
		fmt.Printf(format, "NetMask", nif.NetMask)
		fmt.Println()
	}
}

func dumpMemInfos() {
	mi, err := libsysinfo.MemInfos()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nMemInfos\n---------\n")
	format := "- %-14s : %d %s\n"

	fmt.Printf(format, "MemTotal", mi.MemTotal, mi.UnitUsed)
	fmt.Printf(format, "MemFree", mi.MemFree, mi.UnitUsed)
	fmt.Printf(format, "Buffers", mi.Buffers, mi.UnitUsed)
	fmt.Printf(format, "Cached", mi.Cached, mi.UnitUsed)
	fmt.Printf(format, "SwapCached", mi.SwapCached, mi.UnitUsed)
	fmt.Printf(format, "SwapTotal", mi.SwapTotal, mi.UnitUsed)
	fmt.Printf(format, "SwapFree", mi.SwapFree, mi.UnitUsed)
}
