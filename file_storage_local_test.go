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
