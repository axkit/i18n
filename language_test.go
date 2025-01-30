package i18n

import (
	"strings"
	"testing"
)

func TestCode(t *testing.T) {
	resetLangState()
	codes = []string{"en", "fr", "de"}
	if got := Code(1); got != "fr" {
		t.Errorf("expected 'fr', got '%s'", got)
	}
	if got := Code(10); got != UnknownCode {
		t.Errorf("expected '%s', got '%s'", UnknownCode, got)
	}

}

func TestIndex(t *testing.T) {
	resetLangState()
	codes = []string{"en", "fr", "de"}
	if got := Index("fr"); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
	if got := Index("es"); got != Unknown {
		t.Errorf("expected Unknown, got %d", got)
	}

	if en := Index("en"); en.String() != "en" {
		t.Errorf("expected 'en', got '%s'", en)
	}

}

func TestLangCodes(t *testing.T) {
	resetLangState()
	codes = []string{"en", "fr", "de"}
	if got := LangCodes(); !equalSlices(got, []string{"en", "fr", "de"}) {
		t.Errorf("expected ['en', 'fr', 'de'], got %v", got)
	}
}

func TestNextLangIndex(t *testing.T) {
	resetLangState()
	codes = []string{"en", "fr"}
	nextCode = []LangIndex{1, Unknown}
	if got := NextLangIndex(0); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
	if got := NextLangIndex(10); got != Unknown {
		t.Errorf("expected Unknown, got %d", got)
	}
}

func TestToName(t *testing.T) {
	resetLangState()
	codes = []string{"en", "fr"}
	data := []byte(`{"en":"Hello","fr":"Bonjour"}`)
	n, err := ToString(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := n.InLang(1); got != "Bonjour" {
		t.Errorf("expected 'Bonjour', got '%s'", got)
	}
}

func TestGetSetLangIndex(t *testing.T) {

	tests := []struct {
		input         []string
		expectedCodes []string
		expectedNext  []LangIndex
	}{
		{
			[]string{"zh-Hans-CN", "en", "zh-Hant"},
			[]string{"zh", "zh-Hans", "zh-Hans-CN", "en", "zh-Hant"},
			[]LangIndex{-1, 0, 1, -1, 0},
		},
		{
			[]string{"en", "zh-Hans-CN", "zh-Hant"},
			[]string{"en", "zh", "zh-Hans", "zh-Hans-CN", "zh-Hant"},
			[]LangIndex{-1, -1, 1, 2, 1},
		},
		{
			[]string{"zh-Hant", "zh-Hans-CN", "en-US", "en", "en-GB"},
			[]string{"zh", "zh-Hant", "zh-Hans", "zh-Hans-CN", "en", "en-US", "en-GB"},
			[]LangIndex{-1, 0, 0, 2, -1, 4, 4},
		},
	}

	for _, tt := range tests {
		resetLangState()
		t.Run(strings.Join(tt.input, ","), func(t *testing.T) {
			for _, code := range tt.input {
				_ = GetSetLangIndex(code)
			}
			if gotCodes := LangCodes(); !equalSlices(gotCodes, tt.expectedCodes) {
				t.Errorf("expected codes %v, got %v", tt.expectedCodes, gotCodes)
			}
			if gotNext := getNextCodesSnapshot(); !equalLanguageIndices(gotNext, tt.expectedNext) {
				t.Errorf("expected nextCode %v, got %v", tt.expectedNext, gotNext)
			}
		})
	}
}

// resetLangState resets the global state of the language package.
func resetLangState() {
	mux.Lock()
	defer mux.Unlock()
	codes = nil
	nextCode = nil
	NoFoundIndex = Unknown
}

func getNextCodesSnapshot() []LangIndex {
	mux.RLock()
	defer mux.RUnlock()
	res := make([]LangIndex, len(nextCode))
	copy(res, nextCode)
	return res
}

// equalSlices is a helper function to compare two string slices.
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// equalLanguageIndices is a helper function to compare two LangIndex slices.
func equalLanguageIndices(a, b []LangIndex) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
