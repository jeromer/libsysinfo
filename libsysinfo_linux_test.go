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
