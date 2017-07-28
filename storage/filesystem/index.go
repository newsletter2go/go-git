package filesystem

import (
	"os"

	"github.com/newsletter2go/go-git/plumbing/format/index"
	"github.com/newsletter2go/go-git/storage/filesystem/internal/dotgit"
	"github.com/newsletter2go/go-git/utils/ioutil"
)

type IndexStorage struct {
	dir *dotgit.DotGit
}

func (s *IndexStorage) SetIndex(idx *index.Index) error {
	f, err := s.dir.IndexWriter()
	if err != nil {
		return err
	}

	defer ioutil.CheckClose(f, &err)

	e := index.NewEncoder(f)
	err = e.Encode(idx)
	return err
}

func (s *IndexStorage) Index() (*index.Index, error) {
	idx := &index.Index{
		Version: 2,
	}

	f, err := s.dir.Index()
	if err != nil {
		if os.IsNotExist(err) {
			return idx, nil
		}

		return nil, err
	}

	defer ioutil.CheckClose(f, &err)

	d := index.NewDecoder(f)
	err = d.Decode(idx)
	return idx, err
}
