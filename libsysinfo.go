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
	}

	globalCache = NewCachedValues(len(cacheKeys))
)

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

func LsbDistCodeName() (string, error) {
	return lsbDist("LSB_DIST_CODE_NAME", "Codename")
}

func LsbDistDescription() (string, error) {
	return lsbDist("LSB_DIST_DESCRIPTION", "Description")
}

func lsbDist(k string, lsbItem string) (string, error) {
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
		return "", nil
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
