package i18n

import (
	"os"
	"path/filepath"
)

// LocalFileStorage implements the FileStorager interface for local files.
type LocalFileStorage struct {
	names []string
}

var _ FileStorager = (*LocalFileStorage)(nil)

func NewLocalFileStorage() *LocalFileStorage {
	return &LocalFileStorage{}
}

// ReadFile reads the file content specified by fullname.
func (s *LocalFileStorage) ReadFile(fullname string) ([]byte, error) {
	return os.ReadFile(fullname)
}

func (s *LocalFileStorage) RegisteredFilenames() []string {
	return s.names
}

// RegisterFiles registers files by mask in the directories specified by paths.
func (s *LocalFileStorage) RegisterFiles(mask string, paths ...string) error {

	for _, dir := range paths {

		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, de := range dirEntries {
			if de.IsDir() {
				continue
			}

			if len(mask) > 0 && mask != "*" {
				ok, err := filepath.Match(mask, de.Name())
				if err != nil {
					return err
				}
				if !ok {
					continue
				}
			}

			fi, err := de.Info()
			if err != nil {
				return err
			}

			s.names = append(s.names, filepath.Join(dir, fi.Name()))
		}
	}
	return nil
}
