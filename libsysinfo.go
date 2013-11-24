// +build linux

package libsysinfo

import (
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
	}

	globalCache = newCachedValues(len(cacheKeys))

	ErrDomainNameNotFound = &LibSysInfoErr{"Domain name not found"}
)

type LibSysInfoErr struct {
	Msg string
}

func (e LibSysInfoErr) Error() string {
	return e.Msg
}

func Hostname() (string, error) {
	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys["HOSTNAME"],
		Fetcher:     getFullHostname,
		Processor:   processHostname,
		CacheBucket: globalCache,
	}

	return llv.run()
}

func Domain() (string, error) {
	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys["DOMAIN_NAME"],
		Fetcher:     getFullHostname,
		Processor:   processDomainName,
		CacheBucket: globalCache,
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
		CacheBucket: globalCache,
	}

	return llv.run()
}

func LsbRelease() (map[string]string, error) {
	items := make(map[string]string, 4)

	v, err := lsbReleaseItem("LSB_DIST_CODE_NAME", "Codename")
	if err != nil {
		return items, err
	}
	items["codename"] = v

	v, err = lsbReleaseItem("LSB_DIST_DESCRIPTION", "Description")
	if err != nil {
		return items, err
	}
	items["description"] = v

	v, err = lsbReleaseItem("LSB_DIST_ID", "Distributor ID")
	if err != nil {
		return items, err
	}
	items["distributorid"] = v

	v, err = lsbReleaseItem("LSB_DIST_RELEASE", "Release")
	if err != nil {
		return items, err
	}
	items["release"] = v

	return items, nil
}

func HostId() (string, error) {
	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys["HOST_ID"],
		Fetcher:     getHostId,
		Processor:   processHostId,
		CacheBucket: globalCache,
	}

	return llv.run()
}

func lsbReleaseItem(k string, lsbItem string) (string, error) {
	proc := func(lsb string) (string, error) {
		return processLsbItem(lsb, lsbItem)
	}

	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys[k],
		Fetcher:     getLsbRelease,
		Processor:   proc,
		CacheBucket: globalCache,
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

func getFullHostname() (string, error) {
	cacheKey := cacheKeys["HOSTNAME_FULL"]

	hf, exists := globalCache.Exists(cacheKey)
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

	globalCache.Set(cacheKey, hf)

	return hf, nil
}

func getLsbRelease() (string, error) {
	cacheKey := cacheKeys["LSB_FULL"]

	lsb, exists := globalCache.Exists(cacheKey)
	if exists {
		return lsb, nil
	}

	out, err := exec.Command("lsb_release", "-a").Output()
	if err != nil {
		return string(out), err
	}

	lsb = string(out)

	globalCache.Set(cacheKey, lsb)

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
