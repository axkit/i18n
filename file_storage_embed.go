package i18n

import (
	"embed"
)

// EmbedFileStorage implements the FileStorager interface for local files.
type EmbedFileStorage struct {
	fs *embed.FS
}

var _ FileStorager = (*EmbedFileStorage)(nil)

func NewEmbedFileStorage(fs *embed.FS) *EmbedFileStorage {
	return &EmbedFileStorage{
		fs: fs,
	}
}

// ReadFile reads the file content specified by name.
func (s *EmbedFileStorage) ReadFile(name string) ([]byte, error) {
	return s.fs.ReadFile(name)
}

// RegisteredFilenames returns the list of files located in the root
// directory of the embedded filesystem.
func (s *EmbedFileStorage) RegisteredFilenames() []string {
	var res []string
	de, _ := s.fs.ReadDir(".")

	for i := range de {
		if de[i].IsDir() {
			continue
		}
		res = append(res, de[i].Name())
	}
	return res
}
