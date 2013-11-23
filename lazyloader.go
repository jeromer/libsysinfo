package libsysinfo

type fetcher func() (string, error)
type processor func(string) (string, error)

type lazyLoadedValue struct {
	CacheKey    string
	Fetcher     fetcher
	Processor   processor
	CacheBucket *cachedValues
}

func (llv *lazyLoadedValue) run() (string, error) {
	v, e := llv.CacheBucket.Exists(llv.CacheKey)
	if e {
		return v, nil
	}

	fetched, err := llv.Fetcher()
	if err != nil {
		return "", err
	}

	result, err := llv.Processor(fetched)
	if err != nil {
		return "", err
	}

	llv.CacheBucket.Set(llv.CacheKey, result)

	return result, nil
}
