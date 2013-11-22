package libsysinfo

import (
	"os/exec"
	"strings"
)

const ()

var (
	cacheKeys = map[string]string{
		"HOSTNAME_FULL":      "fullhostname",
		"HOSTNAME":           "hostname",
		"DOMAIN_NAME":        "domainname",
		"LSB_FULL":           "lsbfull",
		"LSB_DIST_CODE_NAME": "lsbdistcodename",
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
	return lsbDist("LSB_DIST_CODE_NAME", processLsbRelease)
}

func lsbDist(k string, procFunc processor) (string, error) {
	llv := &lazyLoadedValue{
		CacheKey:    cacheKeys[k],
		Fetcher:     getLsbRelease,
		Processor:   procFunc,
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

func processLsbRelease(lsb string) (string, error) {
	var out string
	for _, line := range strings.Split(lsb, "\n") {
		if len(line) <= 0 {
			continue
		}

		// C => Codename
		if line[0] != 'C' {
			continue
		}

		out = strings.TrimSpace(strings.TrimLeft(line, "Codename:"))

		break
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
