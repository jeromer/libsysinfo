// +build linux

package libsysinfo

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const ()

var (
	cacheKeys = map[string]string{
		"HOSTNAME_FULL":        "fullhostname",
		"HOSTNAME":             "hostname",
		"DOMAIN_NAME":          "domainname",
		"LSB_FULL":             "lsbfull",
		"LSB_DIST_CODE_NAME":   "lsbdistcodename",
		"LSB_DIST_DESCRIPTION": "lsbdistdescrption",
		"LSB_DIST_ID":          "lsbdistid",
		"LSB_DIST_RELEASE":     "lsbdistrelease",
		"HOST_ID":              "hostid",
		"FILE_SYSTEMS":         "filesystems",
	}

	simpleValuesCache     = newCachedValues(len(cacheKeys))
	fileSystemCache       []string
	cpuInfoCache          []CpuInfo
	networkInterfaceCache []NetworkInterface
	memInfoCache          Meminfos

	ErrDomainNameNotFound = &LibSysInfoErr{"Domain name not found"}
	ErrNoNetIfaceFound    = &LibSysInfoErr{"No network interface found"}
	ErrIfConfigNotFound   = &LibSysInfoErr{"No ifconfig command found"}
)

// ----

type CpuInfo struct {
	Processor      string
	VendorId       string
	CpuFamily      string
	Model          string
	ModelName      string
	Stepping       int
	CPUMHz         float64
	CacheSize      int
	CacheSizeUnit  string
	PhysicalId     string
	Siblings       int
	CoreId         string
	CpuCores       int
	ApicId         string
	InitialApicId  string
	Fpu            string
	FpuException   string
	CpuIdLevel     int
	Wp             string
	Flags          []string
	Bogomips       float64
	ClflushSize    int
	CacheAlignment int
	AddressSizes   string
}

type LsbReleaseInfo struct {
	Codename      string
	Description   string
	DistributorId string
	Release       string
}

type NetworkInterface struct {
	Name          string
	V4Addr        string
	V6Addr        string
	MacAddr       string
	BroadcastAddr string
	NetMask       string
}

type Meminfos struct {
	MemTotal   int
	MemFree    int
	Buffers    int
	Cached     int
	SwapCached int
	SwapTotal  int
	SwapFree   int

	// The unit used in /proc/meminfo, lowered. Most likely always "kb"
	UnitUsed string
}

// ----

func Hostname() (string, error) {
	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys["HOSTNAME"],
		Fetcher:     getFullHostname,
		Processor:   processHostname,
		CacheBucket: simpleValuesCache,
	}

	return llv.run()
}

func Domain() (string, error) {
	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys["DOMAIN_NAME"],
		Fetcher:     getFullHostname,
		Processor:   processDomainName,
		CacheBucket: simpleValuesCache,
	}

	return llv.run()
}

func Fqdn() (string, error) {
	fqdn := func(fullHostname string) (string, error) {
		return fullHostname, nil
	}

	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys["FQDN"],
		Fetcher:     getFullHostname,
		Processor:   fqdn,
		CacheBucket: simpleValuesCache,
	}

	return llv.run()
}

func LsbRelease() (LsbReleaseInfo, error) {
	var lsbr LsbReleaseInfo

	v, err := lsbReleaseItem("LSB_DIST_CODE_NAME", "Codename")
	if err != nil {
		return lsbr, err
	}
	lsbr.Codename = v

	v, err = lsbReleaseItem("LSB_DIST_DESCRIPTION", "Description")
	if err != nil {
		return lsbr, err
	}
	lsbr.Description = v

	v, err = lsbReleaseItem("LSB_DIST_ID", "Distributor ID")
	if err != nil {
		return lsbr, err
	}
	lsbr.DistributorId = v

	v, err = lsbReleaseItem("LSB_DIST_RELEASE", "Release")
	if err != nil {
		return lsbr, err
	}
	lsbr.Release = v

	return lsbr, nil
}

func HostId() (string, error) {
	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys["HOST_ID"],
		Fetcher:     getHostId,
		Processor:   processHostId,
		CacheBucket: simpleValuesCache,
	}

	return llv.run()
}

func FileSystems() ([]string, error) {
	if len(fileSystemCache) > 0 {
		return fileSystemCache, nil
	}

	buff, err := getFileSystems()
	if err != nil {
		return fileSystemCache, err
	}

	return processFileSystems(buff), nil
}

func CpuInfos() ([]CpuInfo, error) {
	if len(cpuInfoCache) > 0 {
		return cpuInfoCache, nil
	}

	buff, err := getCpuInfos()
	if err != nil {
		return []CpuInfo(nil), err
	}

	return processCpuInfos(buff), nil
}

func NetworkInterfaces() ([]NetworkInterface, error) {
	// XXX : switch to a cgo/iotctl based implementation
	// XXX : parsing ifconfig's result is a PITA

	if len(networkInterfaceCache) > 0 {
		return networkInterfaceCache, nil
	}

	var ifaces []NetworkInterface
	devices, err := findNetworkDevices()
	if err != nil {
		return []NetworkInterface{}, err
	}
	if len(devices) <= 0 {
		return ifaces, ErrNoNetIfaceFound
	}

	ifConfig, err := findIfconfig()
	if err != nil {
		return ifaces, err
	}

	for _, d := range devices {
		out, err := exec.Command(ifConfig, d).Output()
		if err != nil {
			return ifaces, err
		}

		ifaces = append(ifaces, processIfconfigOutput(d, string(out)))
	}

	return ifaces, nil
}

func MemInfos() (Meminfos, error) {
	if memInfoCache.MemTotal > 0 {
		return memInfoCache, nil
	}

	buff, err := getMemInfos()
	if err != nil {
		return Meminfos{}, err
	}

	return processMemInfos(buff), nil
}

// ----

func findNetworkDevices() ([]string, error) {
	var devs []string

	f, err := os.Open("/sys/class/net/")
	if err != nil {
		return devs, err
	}
	defer f.Close()

	allNames := -1
	names, err := f.Readdirnames(allNames)
	if err != nil {
		return devs, err
	}

	if len(names) <= 0 {
		return names, ErrNoNetIfaceFound
	}

	return names, nil
}

func findIfconfig() (string, error) {
	possiblePaths := []string{
		"/sbin/ifconfig",
		"/bin/ifconfig",
		"/usr/sbin/ifconfig",
	}

	var f *os.File
	var err error

	for _, path := range possiblePaths {
		f, err = os.Open(path)
		if os.IsNotExist(err) {
			continue
		}
		defer f.Close()

		return path, nil
	}

	return "", ErrIfConfigNotFound
}

func lsbReleaseItem(k string, lsbItem string) (string, error) {
	proc := func(lsb string) (string, error) {
		return processLsbItem(lsb, lsbItem)
	}

	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys[k],
		Fetcher:     getLsbRelease,
		Processor:   proc,
		CacheBucket: simpleValuesCache,
	}

	return llv.run()
}

func processDomainName(fullHostname string) (string, error) {
	pos := strings.Index(fullHostname, ".")
	if pos == -1 {
		return "", ErrDomainNameNotFound
	}

	return fullHostname[pos+1:], nil
}

func processHostname(fullHostname string) (string, error) {
	pos := strings.Index(fullHostname, ".")
	if pos == -1 {
		return fullHostname, nil
	}

	return fullHostname[:pos], nil
}

func processLsbItem(lsb string, item string) (string, error) {
	var out string
	var tmp string

	for _, line := range strings.Split(lsb, "\n") {
		if len(line) <= 0 {
			continue
		}

		tmp = line[0:len(item)]
		if tmp == item {
			out = strings.TrimSpace(strings.TrimLeft(line, item+":"))
			break
		}
	}

	return strings.ToLower(out), nil
}

func processHostId(id string) (string, error) {
	return strings.Trim(id, "\n"), nil
}

func processFileSystems(buff string) []string {
	var tmp string
	var fileSystems []string
	var isNodev bool

	for _, line := range strings.Split(buff, "\n") {
		isNodev = len(line) <= 0 || line[0] == 'n'
		if isNodev {
			continue
		}

		tmp = strings.TrimSpace(line)
		if len(tmp) <= 0 {
			continue
		}

		fileSystems = append(fileSystems, tmp)
	}

	return fileSystems
}

func processCpuInfos(buff string) []CpuInfo {
	var parts []string
	var k, v string
	var cpuInfos []CpuInfo
	var tmp CpuInfo

	lines := strings.Split(buff, "\n")
	lineCount := len(lines)

	for i, line := range lines {
		if line == "" {
			// extra empty lines means end of file
			if i+1 == lineCount {
				break
			}

			cpuInfos = append(cpuInfos, tmp)
			tmp = CpuInfo{}
			continue
		}

		parts = strings.Split(line, ":")
		if len(parts) == 2 {
			k = strings.ToLower(strings.TrimSpace(parts[0]))
			v = strings.TrimSpace(parts[1])
			if v == "" {
				continue
			}

			switch k {
			case "processor":
				tmp.Processor = v
			case "vendor_id":
				tmp.VendorId = v
			case "cpu family":
				tmp.CpuFamily = v
			case "model":
				tmp.Model = v
			case "model name":
				tmp.ModelName = v
			case "stepping":
				tmp.Stepping = atoi(v)
			case "cpu mhz":
				tmp.CPUMHz = atof64(v)
			case "cache size":
				cacheSize := strings.Split(v, " ")
				tmp.CacheSize = atoi(cacheSize[0])
				tmp.CacheSizeUnit = cacheSize[1]
			case "physical id":
				tmp.PhysicalId = v
			case "siblings":
				tmp.Siblings = atoi(v)
			case "core id":
				tmp.CoreId = v
			case "cpu cores":
				tmp.CpuCores = atoi(v)
			case "apicid":
				tmp.ApicId = v
			case "initial apicid":
				tmp.InitialApicId = v
			case "fpu":
				tmp.Fpu = v
			case "fpu_exception":
				tmp.FpuException = v
			case "cpuid level":
				tmp.CpuIdLevel = atoi(v)
			case "wp":
				tmp.Wp = v
			case "flags":
				tmp.Flags = strings.Split(v, " ")
			case "bogomips":
				tmp.Bogomips = atof64(v)
			case "clflush size":
				tmp.ClflushSize = atoi(v)
			case "cache_alignment":
				tmp.CacheAlignment = atoi(v)
			case "address sizes":
				tmp.AddressSizes = v
			default:
				continue
			}
		}
	}

	return cpuInfos
}

func processIfconfigOutput(device string, out string) NetworkInterface {
	// XXX : this is so horrible ...
	const hwaddr = "HWaddr"
	var nif NetworkInterface

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if len(line) <= 0 {
			continue
		}
		isFirstLine := (line[0] != ' ')

		if isFirstLine && strings.Contains(line, hwaddr) {
			parts := strings.Split(line, hwaddr)
			nif.MacAddr = strings.TrimSpace(parts[len(parts)-1])
			continue
		}

		line = strings.TrimSpace(line)

		// i => inet or inet6
		if len(line) <= 0 || line[0] != 'i' {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) <= 0 {
			continue
		}

		if fields[0] == "inet" {
			for _, f := range fields[1:] {
				if f[0] == 'a' {
					nif.V4Addr = strings.TrimPrefix(f, "addr:")
					continue
				}

				if f[0] == 'B' {
					nif.BroadcastAddr = strings.TrimPrefix(f, "Bcast:")
					continue
				}

				if f[0] == 'M' {
					nif.NetMask = strings.TrimPrefix(f, "Mask:")
					continue
				}
			}

			continue
		}

		if fields[0] == "inet6" {
			if len(fields) < 3 {
				continue
			}

			nif.V6Addr = fields[2]
		}
	}

	nif.Name = device

	return nif
}

func processMemInfos(buff string) Meminfos {
	var parts []string
	var k, v, u string
	var mi Meminfos

	lines := strings.Split(buff, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts = strings.Fields(line)
		if len(parts) != 3 {
			continue
		}

		k = strings.ToLower(strings.TrimSpace(strings.Trim(parts[0], ":")))
		v = strings.TrimSpace(parts[1])
		if v == "" {
			continue
		}
		if mi.UnitUsed == "" {
			u = strings.TrimSpace(parts[2])
			mi.UnitUsed = strings.ToLower(u)
		}

		switch k {
		case "memtotal":
			mi.MemTotal = atoi(v)
		case "memfree":
			mi.MemFree = atoi(v)
		case "buffers":
			mi.Buffers = atoi(v)
		case "cached":
			mi.Cached = atoi(v)
		case "swapcached":
			mi.SwapCached = atoi(v)
		case "swaptotal":
			mi.SwapTotal = atoi(v)
		case "swapfree":
			mi.SwapFree = atoi(v)
		default:
			continue
		}
	}

	return mi
}

// ----

func getFullHostname() (string, error) {
	cacheKey := cacheKeys["HOSTNAME_FULL"]

	hf, exists := simpleValuesCache.Exists(cacheKey)
	if exists {
		return hf, nil
	}

	out, err := exec.Command("hostname", "-f").Output()
	if err != nil {
		return string(out), err
	}

	hf = string(out)

	// removing \n
	pos := strings.Index(hf, "\n")
	if pos > -1 {
		hf = hf[:pos]
	}

	simpleValuesCache.Set(cacheKey, hf)

	return hf, nil
}

func getLsbRelease() (string, error) {
	cacheKey := cacheKeys["LSB_FULL"]

	lsb, exists := simpleValuesCache.Exists(cacheKey)
	if exists {
		return lsb, nil
	}

	out, err := exec.Command("lsb_release", "-a").Output()
	if err != nil {
		return string(out), err
	}

	lsb = string(out)

	simpleValuesCache.Set(cacheKey, lsb)

	return lsb, nil
}

func getHostId() (string, error) {
	// XXX : hostid will not be reused, no need to cache the full output
	out, err := exec.Command("hostid").Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func getFileSystems() (string, error) {
	buff, err := ioutil.ReadFile("/proc/filesystems")
	if err != nil {
		return "", nil
	}

	return string(buff), nil
}

func getCpuInfos() (string, error) {
	buff, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	return string(buff), err
}

func getMemInfos() (string, error) {
	buff, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return "", err
	}

	return string(buff), err
}

// ----

type LibSysInfoErr struct {
	Msg string
}

func (e LibSysInfoErr) Error() string {
	return e.Msg
}
