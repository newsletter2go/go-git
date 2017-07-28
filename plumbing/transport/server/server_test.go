package server_test

import (
	"testing"

	"github.com/src-d/go-git-fixtures"
	"github.com/newsletter2go/go-git/plumbing/transport"
	"github.com/newsletter2go/go-git/plumbing/transport/client"
	"github.com/newsletter2go/go-git/plumbing/transport/server"
	"github.com/newsletter2go/go-git/storage/filesystem"
	"github.com/newsletter2go/go-git/storage/memory"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	fixtures.Suite
	loader       server.MapLoader
	client       transport.Transport
	clientBackup transport.Transport
	asClient     bool
}

func (s *BaseSuite) SetUpSuite(c *C) {
	s.Suite.SetUpSuite(c)
	s.loader = server.MapLoader{}
	if s.asClient {
		s.client = server.NewClient(s.loader)
	} else {
		s.client = server.NewServer(s.loader)
	}

	s.clientBackup = client.Protocols["file"]
	client.Protocols["file"] = s.client
}

func (s *BaseSuite) TearDownSuite(c *C) {
	if s.clientBackup == nil {
		delete(client.Protocols, "file")
	} else {
		client.Protocols["file"] = s.clientBackup
	}
}

func (s *BaseSuite) prepareRepositories(c *C, basic *transport.Endpoint,
	empty *transport.Endpoint, nonExistent *transport.Endpoint) {

	f := fixtures.Basic().One()
	fs := f.DotGit()
	path := fs.Root()
	ep, err := transport.NewEndpoint(path)
	c.Assert(err, IsNil)
	*basic = ep
	sto, err := filesystem.NewStorage(fs)
	c.Assert(err, IsNil)
	s.loader[ep.String()] = sto

	path = "/empty.git"
	ep, err = transport.NewEndpoint(path)
	c.Assert(err, IsNil)
	*empty = ep
	s.loader[ep.String()] = memory.NewStorage()

	path = "/non-existent.git"
	ep, err = transport.NewEndpoint(path)
	c.Assert(err, IsNil)
	*nonExistent = ep
}
