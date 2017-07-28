package packfile_test

import (
	"bytes"
	"math/rand"

	"github.com/newsletter2go/go-git/plumbing"
	. "github.com/newsletter2go/go-git/plumbing/format/packfile"
	"github.com/newsletter2go/go-git/plumbing/storer"
	"github.com/newsletter2go/go-git/storage/filesystem"
	"github.com/newsletter2go/go-git/storage/memory"

	"github.com/src-d/go-git-fixtures"
	. "gopkg.in/check.v1"
)

type EncoderAdvancedSuite struct {
	fixtures.Suite
}

var _ = Suite(&EncoderAdvancedSuite{})

func (s *EncoderAdvancedSuite) TestEncodeDecode(c *C) {
	fixs := fixtures.Basic().ByTag("packfile").ByTag(".git")
	fixs = append(fixs, fixtures.ByURL("https://github.com/src-d/go-git.git").
		ByTag("packfile").ByTag(".git").One())
	fixs.Test(c, func(f *fixtures.Fixture) {
		storage, err := filesystem.NewStorage(f.DotGit())
		c.Assert(err, IsNil)
		s.testEncodeDecode(c, storage)
	})

}

func (s *EncoderAdvancedSuite) testEncodeDecode(c *C, storage storer.Storer) {

	objIter, err := storage.IterEncodedObjects(plumbing.AnyObject)
	c.Assert(err, IsNil)

	expectedObjects := map[plumbing.Hash]bool{}
	var hashes []plumbing.Hash
	err = objIter.ForEach(func(o plumbing.EncodedObject) error {
		expectedObjects[o.Hash()] = true
		hashes = append(hashes, o.Hash())
		return err

	})
	c.Assert(err, IsNil)

	// Shuffle hashes to avoid delta selector getting order right just because
	// the initial order is correct.
	auxHashes := make([]plumbing.Hash, len(hashes))
	for i, j := range rand.Perm(len(hashes)) {
		auxHashes[j] = hashes[i]
	}
	hashes = auxHashes

	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, storage, false)
	_, err = enc.Encode(hashes)
	c.Assert(err, IsNil)

	scanner := NewScanner(buf)
	storage = memory.NewStorage()
	d, err := NewDecoder(scanner, storage)
	c.Assert(err, IsNil)
	_, err = d.Decode()
	c.Assert(err, IsNil)

	objIter, err = storage.IterEncodedObjects(plumbing.AnyObject)
	c.Assert(err, IsNil)
	obtainedObjects := map[plumbing.Hash]bool{}
	err = objIter.ForEach(func(o plumbing.EncodedObject) error {
		obtainedObjects[o.Hash()] = true
		return nil
	})
	c.Assert(err, IsNil)
	c.Assert(obtainedObjects, DeepEquals, expectedObjects)

	for h := range obtainedObjects {
		if !expectedObjects[h] {
			c.Errorf("obtained unexpected object: %s", h)
		}
	}

	for h := range expectedObjects {
		if !obtainedObjects[h] {
			c.Errorf("missing object: %s", h)
		}
	}
}
