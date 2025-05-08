package i18n

import (
	"path/filepath"
	"strings"
)

type DefaultFilenameParser struct{}

var _ FilenameParser = (*DefaultFilenameParser)(nil)

// ExtractFilename returns the file name from the full path.
func (DefaultFilenameParser) ExtractFilename(fullname string) (string, error) {

	_, name := filepath.Split(fullname)

	return name, nil
}

// ParseFileName returns the language index and suffix from the filename.
//
// Example:
// ParseFileName("en.t18n") returns English, ""
// ParseFileName("en.grid.t18n") returns English, "grid"
func (DefaultFilenameParser) ParseFilename(filename string) (li Language, suffix string) {
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
