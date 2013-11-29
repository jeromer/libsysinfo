package libsysinfo

import (
	. "launchpad.net/gocheck"
)

type LibSysInfoTestSuite struct{}

var (
	_ = Suite(&LibSysInfoTestSuite{})
)

func (s *LibSysInfoTestSuite) TestProcessHostName(c *C) {
	hostnames := []string{
		"wheezy64-puppet3",
		"wheezy64-puppet3.vagrantup.com",
	}

	expected := "wheezy64-puppet3"

	var obtained string
	var err error
	for _, h := range hostnames {
		obtained, err = processHostname(h)

		c.Assert(err, IsNil)
		c.Assert(obtained, Equals, expected)
	}
}

func (s *LibSysInfoTestSuite) TestProcessDomainName_Found(c *C) {
	h := "wheezy64-puppet3"
	d := "vagrantup.com"

	obtained, err := processDomainName(h + "." + d)

	c.Assert(err, IsNil)
	c.Assert(obtained, Equals, d)
}

func (s *LibSysInfoTestSuite) TestProcessDomainName_NotFound(c *C) {
	obtained, err := processDomainName("none")

	c.Assert(err, Equals, ErrDomainNameNotFound)
	c.Assert(obtained, Equals, "")
}

func (s *LibSysInfoTestSuite) TestProcessHostId(c *C) {
	id := "007f0101"

	obtained, err := processHostId(id + "\n")

	c.Assert(err, IsNil)
	c.Assert(obtained, Equals, id)
}

func (s *LibSysInfoTestSuite) TestProcessFileSystems(c *C) {
	fixture := `

nodev	mqueue
	ext3
	ext2
nodev	rpc_pipefs
nodev	nfs
nodev	nfs4
nodev	nfsd
nodev	vboxsf

	`
	expected := []string{"ext3", "ext2"}
	obtained := processFileSystems(fixture)
	c.Assert(obtained, DeepEquals, expected)
}

func (s *LibSysInfoTestSuite) TestProcessCpuInfos(c *C) {
	fixtures := `processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 42
model name	: Intel(R) Pentium(R) CPU G630T @ 2.30GHz
stepping	: 7
cpu MHz		: 2294.833
cache size	: 3072 KB
physical id	: 0
siblings	: 2
core id		: 0
cpu cores	: 2
apicid		: 0
initial apicid	: 0
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx rdtscp lm constant_tsc arch_perfmon pebs bts rep_good xtopology nonstop_tsc aperfmperf pni pclmulqdq dtes64 monitor ds_cpl vmx est tm2 ssse3 cx16 xtpr pdcm sse4_1 sse4_2 popcnt xsave lahf_lm arat tpr_shadow vnmi flexpriority ept vpid
bogomips	: 4589.66
clflush size	: 64
cache_alignment	: 64
address sizes	: 36 bits physical, 48 bits virtual
power management:

processor	: 1
vendor_id	: GenuineIntel
cpu family	: 6
model		: 42
model name	: Intel(R) Pentium(R) CPU G630T @ 2.30GHz
stepping	: 7
cpu MHz		: 2294.833
cache size	: 3072 KB
physical id	: 0
siblings	: 2
core id		: 1
cpu cores	: 2
apicid		: 2
initial apicid	: 2
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx rdtscp lm constant_tsc arch_perfmon pebs bts rep_good xtopology nonstop_tsc aperfmperf pni pclmulqdq dtes64 monitor ds_cpl vmx est tm2 ssse3 cx16 xtpr pdcm sse4_1 sse4_2 popcnt xsave lahf_lm arat tpr_shadow vnmi flexpriority ept vpid
bogomips	: 4589.37
clflush size	: 64
cache_alignment	: 64
address sizes	: 36 bits physical, 48 bits virtual
power management:

`
	obtained := processCpuInfos(fixtures)

	expected := []CpuInfo{
		CpuInfo{
			Processor:     "0",
			VendorId:      "GenuineIntel",
			CpuFamily:     "6",
			Model:         "42",
			ModelName:     "Intel(R) Pentium(R) CPU G630T @ 2.30GHz",
			Stepping:      "7",
			CPUMHz:        "2294.833",
			CacheSize:     "3072",
			CacheSizeUnit: "KB",
			PhysicalId:    "0",
			Siblings:      "2",
			CoreId:        "0",
			CpuCores:      "2",
			ApicId:        "0",
			InitialApicId: "0",
			Fpu:           "yes",
			FpuException:  "yes",
			CpuIdLevel:    "13",
			Wp:            "yes",
			Flags: []string{
				"fpu",
				"vme",
				"de",
				"pse",
				"tsc",
				"msr",
				"pae",
				"mce",
				"cx8",
				"apic",
				"sep",
				"mtrr",
				"pge",
				"mca",
				"cmov",
				"pat",
				"pse36",
				"clflush",
				"dts",
				"acpi",
				"mmx",
				"fxsr",
				"sse",
				"sse2",
				"ss",
				"ht",
				"tm",
				"pbe",
				"syscall",
				"nx",
				"rdtscp",
				"lm",
				"constant_tsc",
				"arch_perfmon",
				"pebs",
				"bts",
				"rep_good",
				"xtopology",
				"nonstop_tsc",
				"aperfmperf",
				"pni",
				"pclmulqdq",
				"dtes64",
				"monitor",
				"ds_cpl",
				"vmx",
				"est",
				"tm2",
				"ssse3",
				"cx16",
				"xtpr",
				"pdcm",
				"sse4_1",
				"sse4_2",
				"popcnt",
				"xsave",
				"lahf_lm",
				"arat",
				"tpr_shadow",
				"vnmi",
				"flexpriority",
				"ept",
				"vpid",
			},
			Bogomips:       "4589.66",
			ClflushSize:    "64",
			CacheAlignment: "64",
			AddressSizes:   "36 bits physical, 48 bits virtual",
		},
		CpuInfo{
			Processor:     "1",
			VendorId:      "GenuineIntel",
			CpuFamily:     "6",
			Model:         "42",
			ModelName:     "Intel(R) Pentium(R) CPU G630T @ 2.30GHz",
			Stepping:      "7",
			CPUMHz:        "2294.833",
			CacheSize:     "3072",
			CacheSizeUnit: "KB",
			PhysicalId:    "0",
			Siblings:      "2",
			CoreId:        "1",
			CpuCores:      "2",
			ApicId:        "2",
			InitialApicId: "2",
			Fpu:           "yes",
			FpuException:  "yes",
			CpuIdLevel:    "13",
			Wp:            "yes",
			Flags: []string{
				"fpu",
				"vme",
				"de",
				"pse",
				"tsc",
				"msr",
				"pae",
				"mce",
				"cx8",
				"apic",
				"sep",
				"mtrr",
				"pge",
				"mca",
				"cmov",
				"pat",
				"pse36",
				"clflush",
				"dts",
				"acpi",
				"mmx",
				"fxsr",
				"sse",
				"sse2",
				"ss",
				"ht",
				"tm",
				"pbe",
				"syscall",
				"nx",
				"rdtscp",
				"lm",
				"constant_tsc",
				"arch_perfmon",
				"pebs",
				"bts",
				"rep_good",
				"xtopology",
				"nonstop_tsc",
				"aperfmperf",
				"pni",
				"pclmulqdq",
				"dtes64",
				"monitor",
				"ds_cpl",
				"vmx",
				"est",
				"tm2",
				"ssse3",
				"cx16",
				"xtpr",
				"pdcm",
				"sse4_1",
				"sse4_2",
				"popcnt",
				"xsave",
				"lahf_lm",
				"arat",
				"tpr_shadow",
				"vnmi",
				"flexpriority",
				"ept",
				"vpid",
			},
			Bogomips:       "4589.37",
			ClflushSize:    "64",
			CacheAlignment: "64",
			AddressSizes:   "36 bits physical, 48 bits virtual",
		},
	}

	c.Assert(len(obtained), Equals, len(expected))
	c.Assert(obtained, DeepEquals, expected)
}

func (s *LibSysInfoTestSuite) TestProcessIfConfigOutput_Loopback(c *C) {
	fixtures := `
lo        Link encap:Local Loopback
          inet addr:127.0.0.1  Mask:255.0.0.0
          inet6 addr: ::1/128 Scope:Host
          UP LOOPBACK RUNNING  MTU:16436  Metric:1
          RX packets:8 errors:0 dropped:0 overruns:0 frame:0
          TX packets:8 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0
          RX bytes:1104 (1.0 KiB)  TX bytes:1104 (1.0 KiB)

`

	obtained := processIfconfigOutput("lo", fixtures)
	expected := NetworkInterface{
		Name:          "lo",
		V4Addr:        "127.0.0.1",
		V6Addr:        "::1/128",
		MacAddr:       "",
		BroadcastAddr: "",
		NetMask:       "255.0.0.0",
	}

	c.Assert(obtained, DeepEquals, expected)
}

func (s *LibSysInfoTestSuite) TestProcessIfConfigOutput_Eth0(c *C) {
	fixtures := `
eth0      Link encap:Ethernet  HWaddr 08:00:27:b3:27:23
          inet addr:10.0.2.15  Bcast:10.0.2.255  Mask:255.255.255.0
          inet6 addr: fe80::a00:27ff:feb3:2723/64 Scope:Link
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:125599 errors:0 dropped:0 overruns:0 frame:0
          TX packets:71401 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:92587916 (88.2 MiB)  TX bytes:6635432 (6.3 MiB)

`

	obtained := processIfconfigOutput("eth0", fixtures)
	expected := NetworkInterface{
		Name:          "eth0",
		V4Addr:        "10.0.2.15",
		V6Addr:        "fe80::a00:27ff:feb3:2723/64",
		MacAddr:       "08:00:27:b3:27:23",
		BroadcastAddr: "10.0.2.255",
		NetMask:       "255.255.255.0",
	}

	c.Assert(obtained, DeepEquals, expected)
}

func (s *LibSysInfoTestSuite) TestProcessMemInfos(c *C) {
	fixtures := `
MemTotal:         250856 kB
MemFree:          152536 kB
Buffers:            4872 kB
Cached:            61592 kB
SwapCached:            0 kB
Active:            44096 kB
Inactive:          35240 kB
Active(anon):      12928 kB
Inactive(anon):      164 kB
Active(file):      31168 kB
Inactive(file):    35076 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:        466940 kB
SwapFree:         466940 kB
Dirty:                 0 kB
Writeback:             0 kB
AnonPages:         12876 kB
Mapped:             6024 kB
Shmem:               220 kB
Slab:              11660 kB
SReclaimable:       4980 kB
SUnreclaim:         6680 kB
KernelStack:         552 kB
PageTables:         1784 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:      592368 kB
Committed_AS:      53608 kB
VmallocTotal:   34359738367 kB
VmallocUsed:       17448 kB
VmallocChunk:   34359719927 kB
HardwareCorrupted:     0 kB
AnonHugePages:         0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:       40896 kB
DirectMap2M:      221184 kB
`
	obtained := processMemInfos(fixtures)

	expected := Meminfos{
		MemTotal:   250856,
		MemFree:    152536,
		Buffers:    4872,
		Cached:     61592,
		SwapCached: 0,
		SwapTotal:  466940,
		SwapFree:   466940,
		UnitUsed:   "kb",
	}

	c.Assert(obtained, DeepEquals, expected)
}
