package i18n

import (
	"sort"
)

var NotFoundMarker = "\u2638"

// TranslationRequestStrategy defines a strategy for handling missing translations.
type TranslationRequestStrategy int8

const (
	ReturnResourceCode TranslationRequestStrategy = iota
	ReturnNotFoundVariable
	ReturnEmptyString
)

type key struct {
	lang Language
	// namespace holds suffix of the localization file. Format: en.{namespace}.t18n
	// empty for default resource file, e.g.: en.t18n or en-US.t18n
	namespace string
}

type file struct {
	key
	name     string
	fullName string // path + file name
}

// Item represents an item in a .t18n file.
type Item struct {
	Key   string
	Value string
	Hint  string
}

// ResponseItem represents a row to be returned to the client.
type ResponseItem struct {
	Value string `json:"v"`
	Hint  string `json:"h,omitempty"`
}

// Set holds a set of items.
type Set struct {
	items []Item
	index map[string]int // key -> index in items
}

// TranslationContainer is a store of all translated resource items.
type TranslationContainer struct {
	cfg          containerConfig
	translations map[key]Set
	files        []file
	customDirs   []string
}

// FileStorager is an interface wrapping the methods for reading files.
//
// FileNamesByMask returns filenames by a mask in directories.
// ExtractFilename extracts a filename from a full path.
// ParseFilename parses a filename and returns a language index and a custom suffix.
// ReadFile reads a file and returns its content.
type FileStorager interface {
	//Filenames(mask string, paths ...string) ([]string, error)
	RegisteredFilenames() []string
	ExtractFilename(fullname string) (string, error)
	ParseFilename(filename string) (Language, string)
	ReadFile(filename string) ([]byte, error)
}

// FileContentParser is an interface wrapping the ParseFileContent method.
//
// ParseFileContent receives a content of a file and returns a slice of items.
type FileContentParser interface {
	ParseFileContent(data []byte) ([]Item, error)
}

// containerConfig defines options for Container.
type containerConfig struct {

	// Default language is used if translation is not found in the requested language.
	primaryLanguage Language

	// strategy defines a strategy for handling missing translations.
	strategy TranslationRequestStrategy

	// suffixPriority maps text and it's priority
	suffixPriority map[string]int

	// bracketSymbol is used to wrap resource key in the client.
	bracketSymbol string

	storage FileStorager

	parser FileContentParser
}

type ContainerOption func(o *containerConfig)

// WithPrimaryLanguage assigns a primary language.
func WithPrimaryLanguage(li Language) ContainerOption {
	return func(o *containerConfig) {
		o.primaryLanguage = li
	}
}

// WithFileSuffixes assigns suffixes of translation files in the order of applying priority.
// The first suffix has the highest priority.
func WithFileSuffixes(suffix ...string) ContainerOption {
	return func(o *containerConfig) {
		for _, s := range suffix {
			if _, ok := o.suffixPriority[s]; !ok {
				o.suffixPriority[s] = len(o.suffixPriority) + 1
			}
		}
	}
}

// WithBrackets assigns wrapping symbol used by the client.
// Example: if bracketSymbol is "%", then the client will
// receive ["%hello%": "ahoj", "%bye%": "cau", "%Hello%": "Ahoj"]
func WithBrackets(bracketSymbol string) ContainerOption {
	return func(o *containerConfig) {
		o.bracketSymbol = bracketSymbol
	}
}

func WithStorage(storage FileStorager) ContainerOption {
	return func(o *containerConfig) {
		o.storage = storage
	}
}

func WithCustomFileParser(parser FileContentParser) ContainerOption {
	return func(o *containerConfig) {
		o.parser = parser
	}
}

func WithStrategy(strategy TranslationRequestStrategy) ContainerOption {
	return func(o *containerConfig) {
		o.strategy = strategy
	}
}

// NewContainer creates a new localization container.
func NewContainer(opts ...func(o *containerConfig)) *TranslationContainer {
	tc := TranslationContainer{
		translations: make(map[key]Set),
		//lang:         GetSetLangIndex("?"),
		cfg: containerConfig{
			primaryLanguage: Unknown,
			suffixPriority:  map[string]int{"": 0}, // file without suffix has the lowest priority.
			storage:         &LocalFileStorage{},
			parser:          &DefaultParser{},
			strategy:        ReturnResourceCode,
		},
	}

	for _, opt := range opts {
		opt(&tc.cfg)
	}
	return &tc
}

// AddFiles registers full .t18n file names in the container.
func (tc *TranslationContainer) addFiles(fullFileNames ...string) error {

	for _, ffn := range fullFileNames {
		name, err := tc.cfg.storage.ExtractFilename(ffn)
		if err != nil {
			return err
		}

		pfi := file{name: name, fullName: ffn}
		pfi.lang, pfi.namespace = tc.cfg.storage.ParseFilename(name)
		tc.files = append(tc.files, pfi)
	}
	return nil
}

func (tc *TranslationContainer) sortFilesBySuffixPriority() {
	sort.Slice(tc.files, func(i, j int) bool {
		if tc.files[i].lang == tc.files[j].lang {
			return tc.cfg.suffixPriority[tc.files[i].namespace] < tc.cfg.suffixPriority[tc.files[j].namespace]
		}
		return tc.files[i].lang < tc.files[j].lang
	})
}

// ReadRegisteredFiles reads content of all registered files, parses and stores content
// in the container.
func (tc *TranslationContainer) ReadRegisteredFiles() error {

	tc.addFiles(tc.cfg.storage.RegisteredFilenames()...)

	tc.sortFilesBySuffixPriority()

	for _, f := range tc.files {
		items, err := tc.loadFile(f.fullName)
		if err != nil {
			return err
		}

		key := key{
			lang:      f.lang,
			namespace: f.namespace,
		}

		if ti, ok := tc.translations[key]; ok {
			// replace
			for j := range items {
				if idx, ok := ti.index[items[j].Key]; ok {
					ti.items[idx] = items[j]
				} else {
					ti.items = append(ti.items, items[j])
					ti.index[items[j].Key] = len(ti.items) - 1
				}
			}
			// important to assign back, because ti is a copy,
			// and ti.items can refer to another address.
			tc.translations[key] = ti
		} else {
			x := Set{index: make(map[string]int)}
			x.items = items
			for i, item := range items {
				x.index[item.Key] = i
			}
			tc.translations[key] = x
		}
	}
	return nil
}

func (tc *TranslationContainer) loadFile(filename string) ([]Item, error) {

	data, err := tc.cfg.storage.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return tc.cfg.parser.ParseFileContent(data)
}

// genKey generates resource key for JSON response.
func (c *TranslationContainer) genKey(id string) string {
	if c.cfg.bracketSymbol == "" {
		return id
	}
	return c.cfg.bracketSymbol + id + c.cfg.bracketSymbol
}
