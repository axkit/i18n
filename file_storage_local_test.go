package i18n_test

import (
	"os"
	"testing"

	"github.com/axkit/i18n"
)

func TestLocalFileStorage_ReadFile(t *testing.T) {
	storage := i18n.NewLocalFileStorage()
	testFile := "./testdata/en.t18n"

	content, err := storage.ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	osContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(content) != string(osContent) {
		t.Errorf("expected %s, got %s", osContent, content)
	}
}

func TestLocalFileStorage_ExtractFilename(t *testing.T) {
	storage := i18n.NewLocalFileStorage()

	fullname := "/path/to/file.txt"
	expected := "file.txt"
	name, err := storage.ExtractFilename(fullname)
	if err != nil {
		t.Fatalf("ExtractFilename failed: %v", err)
	}
	if name != expected {
		t.Errorf("expected %s, got %s", expected, name)
	}
}

func TestLocalFileStorage_RegisterFiles(t *testing.T) {

	if err := os.Mkdir("./testdata/subdir", 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	defer os.RemoveAll("./testdata/subdir")

	storage := i18n.NewLocalFileStorage()

	err := storage.RegisterFiles("*.t18n", "./testdata")
	if err != nil {
		t.Fatalf("RegisterFiles failed: %v", err)
	}

	files := storage.RegisteredFilenames()
	if len(files) == 0 {
		t.Errorf("expected some .t18n files to be registered, but got none")
	}

	if err := storage.RegisterFiles("*.t18n", "./testdata/invalid-subdir"); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestLocalFileStorage_ParseFilename(t *testing.T) {

	storage := i18n.NewLocalFileStorage()

	tests := []struct {
		filename       string
		expectedLang   i18n.LangIndex
		expectedSuffix string
	}{
		{"en.t18n", i18n.GetSetLangIndex("en"), ""},
		{"en-US.customer1.t18n", i18n.GetSetLangIndex("en-US"), "customer1"},
		{"en-GB.t18n", i18n.GetSetLangIndex("en-GB"), ""},
		{"de.customer2.t18n", i18n.GetSetLangIndex("de"), "customer2"},
		{"invalidname", i18n.Unknown, ""},
	}

	for _, tt := range tests {
		lang, suffix := storage.ParseFilename(tt.filename)
		if lang != tt.expectedLang || suffix != tt.expectedSuffix {
			t.Errorf("for %s expected (%v, %s), got (%v, %s)", tt.filename, tt.expectedLang, tt.expectedSuffix, lang, suffix)
		}
	}
}
