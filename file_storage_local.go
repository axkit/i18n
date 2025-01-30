package i18n

import (
	"os"
	"path/filepath"
	"strings"
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

// ExtractFilename returns the file name from the full path.
func (s *LocalFileStorage) ExtractFilename(fullname string) (string, error) {

	_, name := filepath.Split(fullname)

	return name, nil
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

// ParseFileName returns the language index and suffix from the filename.
//
// Example:
// ParseFileName("en.t18n") returns English, ""
// ParseFileName("en.grid.t18n") returns English, "grid"
func (s *LocalFileStorage) ParseFilename(filename string) (li LangIndex, suffix string) {
	from := strings.Index(filename, ".")
	to := strings.LastIndex(filename, ".")
	if from == -1 && to == -1 {
		// it also covers the case when filename has not "."
		return Unknown, ""
	}

	if from == to {
		return GetSetLangIndex(filename[0:from]), ""
	}

	return GetSetLangIndex(filename[0:from]), filename[from+1 : to]
}
