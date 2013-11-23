package libsysinfo

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type CacheTestSuite struct{}

var (
	_ = Suite(&CacheTestSuite{})

	key = "foo"
	val = "bar"
)

func (s *CacheTestSuite) TestSet(c *C) {
	set := newCachedValues(1)

	set.Set(key, val)

	all := set.All()
	c.Assert(len(all), Equals, 1)
	c.Assert(all[key], DeepEquals, val)
}

func (s *CacheTestSuite) TestExists(c *C) {
	set := newCachedValues(1)

	v, e := set.Exists(key)
	c.Assert(v, Equals, "")
	c.Assert(e, Equals, false)

	set.Set(key, val)

	v, e = set.Exists(key)
	c.Assert(v, Equals, val)
	c.Assert(e, Equals, true)
}

func (s *CacheTestSuite) TestGet(c *C) {
	set := newCachedValues(1)

	set.Set(key, val)

	obtained := set.Get(key)
	c.Assert(obtained, Equals, val)
}

func (s *CacheTestSuite) TestDelete(c *C) {
	set := newCachedValues(1)

	set.Set(key, val)
	c.Assert(len(set.All()), Equals, 1)

	set.Delete(key)
	c.Assert(len(set.All()), Equals, 0)
}

func (s *CacheTestSuite) TestEmpty(c *C) {
	set := newCachedValues(1)

	set.Set(key, val)

	all := set.All()
	c.Assert(len(all), Equals, 1)
	c.Assert(all[key], DeepEquals, val)

	set.Empty()
	c.Assert(len(set.All()), Equals, 0)
}
