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
