package i18n

import (
	"strings"
	"sync"
)

// LangIndex is index of the language in the slice of language codes.
type LangIndex int

// Unknown specifies that language code is not found.
const Unknown LangIndex = -1

// UnknownCode is used if language code is not found.
var UnknownCode string = "?"

var (
	mux      sync.RWMutex //
	codes    []string     // slice of language codes
	nextCode []LangIndex  // index of the next language code in the hierarchy
)

// String implements fmt.Stringer interface.
func (li LangIndex) String() string {
	return Code(li)
}

// Code returns language code by language index.
// Returns UnknownCode variable if language index is invalid.
func Code(li LangIndex) string {
	mux.RLock()
	defer mux.RUnlock()
	if int(li) >= len(codes) || li < 0 {
		return UnknownCode
	}
	return codes[li]
}

// Index returns index of the language by language code.
// Returns Unknown if language code is not found.
func Index(code string) LangIndex {
	mux.RLock()
	defer mux.RUnlock()
	return getLangIndex(code)
}

// GetSetLangIndex parses language code and returns index of the language.
// If language code is not found, it adds it.
// If language code is complex, like zh-Hans-CN, than adds zh-Hans-CN, zh-Hans, zh.
// Returns index of the first added language code.
func GetSetLangIndex(code string) LangIndex {

	mux.RLock()
	res := getLangIndex(code)
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
			next = getLangIndex(parts[i+1])
		}
		li, added := getSetLangIndex(parts[i])
		if added {
			nextCode[li] = next
			res = li
		}
	}

	return res
}

// getLangIndex returns index of the language code.
// Returns Unknown if language code is not found.
func getLangIndex(code string) (li LangIndex) {
	cnt := len(codes)
	for i := 0; i < cnt; i++ {
		if codes[i] == code {
			return LangIndex(i)
		}
	}
	return Unknown
}

// getSetLangIndex returns index of the language code.
// If language code is not found, it adds it.
func getSetLangIndex(code string) (li LangIndex, added bool) {

	if li := getLangIndex(code); li != Unknown {
		return li, false
	}

	codes = append(codes, code)
	nextCode = append(nextCode, Unknown)
	return LangIndex(len(codes) - 1), true
}

// LangCodes returns copy of the slice of language codes.
func LangCodes() []string {
	mux.RLock()
	defer mux.RUnlock()
	res := make([]string, len(codes))
	copy(res, codes)
	return res
}

// NextLangIndex returns index of the next language code in the hierarchy.
// Returns Unknown if language index is invalid or is the last one (Unknown).
func NextLangIndex(li LangIndex) LangIndex {
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
