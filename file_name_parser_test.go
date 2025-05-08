package i18n_test

import (
	"testing"

	"github.com/axkit/i18n"
)

func TestDefaultFilenameParser_ExtractFilename(t *testing.T) {
	parser := i18n.DefaultFilenameParser{}

	fullname := "/path/to/file.txt"
	expected := "file.txt"
	name, err := parser.ExtractFilename(fullname)
	if err != nil {
		t.Fatalf("ExtractFilename failed: %v", err)
	}
	if name != expected {
		t.Errorf("expected %s, got %s", expected, name)
	}
}

func TestDefaultFilenameParser_ParseFilename(t *testing.T) {

	parser := i18n.DefaultFilenameParser{}

	tests := []struct {
		filename       string
		expectedLang   i18n.Language
		expectedSuffix string
	}{
		{"en.t18n", i18n.Parse("en"), ""},
		{"en-US.customer1.t18n", i18n.Parse("en-US"), "customer1"},
		{"en-GB.t18n", i18n.Parse("en-GB"), ""},
		{"de.customer2.t18n", i18n.Parse("de"), "customer2"},
		{"invalidname", i18n.Unknown, ""},
	}

	for _, tt := range tests {
		lang, suffix := parser.ParseFilename(tt.filename)
		if lang != tt.expectedLang || suffix != tt.expectedSuffix {
			t.Errorf("for %s expected (%v, %s), got (%v, %s)", tt.filename, tt.expectedLang, tt.expectedSuffix, lang, suffix)
		}
	}
}
