package i18n

import (
	"strings"
	"sync"
)

// Language is index of the language in the slice of language codes.
type Language int

// Unknown holds Language value for unknown language.
const Unknown Language = -1

// UnknownLanguageCode is used if Language is unknown.
var UnknownLanguageCode string = "?"

var (
	// Global variables protected by mutex.
	mux      sync.RWMutex //
	codes    []string     // slice of language codes
	nextCode []Language   // index of the next language code in the hierarchy
)

// String implements fmt.Stringer interface.
func (li Language) String() string {
	return code(li)
}

// code returns language code by language index.
// Returns UnknownCode variable if language index is invalid.
func code(li Language) string {
	mux.RLock()
	defer mux.RUnlock()
	if int(li) >= len(codes) || li < 0 {
		return UnknownLanguageCode
	}
	return codes[li]
}

// Lookup returns Language by language code.
// Returns Unknown if language code is not found.
func Lookup(code string) Language {
	mux.RLock()
	defer mux.RUnlock()
	return getLanguage(code)
}

// Parse parses language code and returns Language.
// If language code is not found, it adds it.
//
// If language code is complex, like zh-Hans-CN, than adds zh-Hans-CN, zh-Hans, zh codes
// and returns index of the first added Language.
func Parse(code string) Language {

	mux.RLock()
	res := getLanguage(code)
	mux.RUnlock()
	if res != Unknown {
		return res
	}

	// following part of the code is rearly used, only when new language is added.
	mux.Lock()
	defer mux.Unlock()

	var parts []string //zh-Hans-CN -> [zh-Hans-CN, zh-Hans, zh]

	for {
		parts = append(parts, code)
		pos := strings.LastIndex(code, "-")
		if pos == -1 {
			break
		}
		code = code[:pos]
	}

	for i := len(parts) - 1; i >= 0; i-- {
		next := Unknown
		if i < len(parts)-1 {
			next = getLanguage(parts[i+1])
		}
		li, added := getSetLanguage(parts[i])
		if added {
			nextCode[li] = next
			res = li
		}
	}

	return res
}

// getLanguage returns Language by code.
// Returns Unknown if code is not found.
func getLanguage(code string) Language {
	cnt := len(codes)
	for i := 0; i < cnt; i++ {
		if codes[i] == code {
			return Language(i)
		}
	}
	return Unknown
}

// getSetLanguage returns Language by language code.
// If language code is not found, than adds it.
func getSetLanguage(code string) (li Language, added bool) {

	if li := getLanguage(code); li != Unknown {
		return li, false
	}

	codes = append(codes, code)
	nextCode = append(nextCode, Unknown)
	return Language(len(codes) - 1), true
}

// LanguageCodes returns copy of the slice of language codes.
func LanguageCodes() []string {
	mux.RLock()
	defer mux.RUnlock()
	res := make([]string, len(codes))
	copy(res, codes)
	return res
}

// LanguageCount returns last Language.
// It is used to iterate over all registered languages.
func LanguageCount() Language {
	mux.RLock()
	defer mux.RUnlock()
	return Language(len(codes) - 1)
}

// NextLanguage returns index of the next language code in the hierarchy.
// Returns Unknown if language index is invalid or is the last one (Unknown).
func NextLanguage(li Language) Language {
	mux.RLock()
	defer mux.RUnlock()
	if int(li) >= len(nextCode) || li < 0 {
		return Unknown
	}
	return nextCode[li]
}

// Len returns number of supported languages.
func Len() int {
	mux.RLock()
	defer mux.RUnlock()
	return len(codes)
}
