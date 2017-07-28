package filesystem

import (
	"github.com/newsletter2go/go-git/storage"
	"github.com/newsletter2go/go-git/storage/filesystem/internal/dotgit"
)

type ModuleStorage struct {
	dir *dotgit.DotGit
}

func (s *ModuleStorage) Module(name string) (storage.Storer, error) {
	fs, err := s.dir.Module(name)
	if err != nil {
		return nil, err
	}

	return NewStorage(fs)
}
