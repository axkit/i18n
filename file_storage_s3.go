package i18n

import (
	"net/url"
	"path"
	"strings"
)

// AmazonS3FileStorage implements the FileStorager interface for Amazon S3 files.
type AmazonS3FileStorage struct{}

func NewAmazonS3FileStorage() *AmazonS3FileStorage {
	return &AmazonS3FileStorage{}
}

// ReadFile reads the file content specified by fullname.
func (s *AmazonS3FileStorage) ReadFile(fullname string) ([]byte, error) {
	// TODO
	return nil, nil
}

// ExtractFilename returns the file name from the full path.
func (s *AmazonS3FileStorage) ExtractFilename(fullname string) (string, error) {

	u, err := url.Parse(fullname)
	if err != nil {
		return "", err
	}

	return path.Base(u.Path), nil
}

// Filenames returns files from the directory specified by dir and mask.
func (s *AmazonS3FileStorage) Filenames(mask string, paths ...string) ([]string, error) {

	var names []string
	// TODO

	return names, nil
}

// ParseFileName returns the language index and suffix from the filename.
//
// Example:
// ParseFileName("en.t18n") returns English, ""
// ParseFileName("en.grid.t18n") returns English, "grid"
func (lfr *AmazonS3FileStorage) ParseFileName(filename string) (li Language, suffix string) {
	from := strings.Index(filename, ".")
	to := strings.LastIndex(filename, ".")
	if from == -1 && to == -1 {
		// it also covers the case when filename has not "."
		return Unknown, ""
	}

	if from == to {
		return Parse(filename[0:from]), ""
	}

	return Parse(filename[0:from]), filename[from+1 : to]
}
