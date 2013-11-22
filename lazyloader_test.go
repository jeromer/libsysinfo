package libsysinfo

import (
	. "launchpad.net/gocheck"
)

type LazyLoaderTestSuite struct {
	cv *CachedValues
}

var (
	_ = Suite(&LazyLoaderTestSuite{})
	k = "foo"
	v = "bar"
)

func (s *LazyLoaderTestSuite) SetUpTest(c *C) {
	s.cv = NewCachedValues(10)
}

func (s *LazyLoaderTestSuite) TearDownTest(c *C) {
	s.cv.Empty()
}

func (s *LazyLoaderTestSuite) TestLazyLoad_CacheExists(c *C) {
	s.cv.Set(k, v)

	llv := &lazyLoadedValue{
		CacheKey:    k,
		Fetcher:     fetcherPanic,
		Processor:   processorPanic,
		CacheBucket: s.cv,
	}

	v, err := llv.run()
	c.Assert(err, IsNil)
	c.Assert(v, Equals, val)
}

func (s *LazyLoaderTestSuite) TestLazyLoad_CacheDoesNotExists(c *C) {
	_, e := s.cv.Exists(k)
	c.Assert(e, Equals, false)

	llv := &lazyLoadedValue{
		CacheKey:    k,
		Fetcher:     fetcherDummy,
		Processor:   processorDummy,
		CacheBucket: s.cv,
	}

	v, err := llv.run()
	c.Assert(err, IsNil)

	_, e = s.cv.Exists(k)
	c.Assert(e, Equals, true)

	c.Assert(v, Equals, "fetchedprocessed")
}

func fetcherPanic() (string, error) {
	panicIfCalled()
	return "", nil
}

func fetcherDummy() (string, error) {
	return "fetched", nil
}

func processorPanic(s string) (string, error) {
	panicIfCalled()
	return "", nil
}

func processorDummy(in string) (string, error) {
	return in + "processed", nil
}

func panicIfCalled() {
	panic("Should not be called")
}
