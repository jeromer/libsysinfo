package libsysinfo

import (
	"os/exec"
	"strings"
)

const ()

var (
	cacheKeys = map[string]string{
		"HOSTNAME_FULL": "fullhostname",
		"HOSTNAME":      "hostname",
		"DOMAIN_NAME":   "domainname",
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
