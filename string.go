package i18n

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// String holds decoded names. Index of the slice calculates by LangIndex.
// In the database it is stored as jsonb: {"en":"Name","cz":"Jméno","sr":"Име"}
type String []string

var (
	// NoValue represents a value that is found.
	NoValue = "$no_value$"
	// NoFoundIndex is used if language code is not found.
	NoFoundIndex = Unknown
)

const emptyString = "#$%!@!@!@!@!"

type StringOption func() string

func WithDefault(s string) StringOption {
	return func() string {
		return s
	}
}

// InLang returns string in language identified by code index.
func (n String) InLang(li LangIndex, opts ...StringOption) string {

	if len(n) == 0 || li < 0 {
		return UnknownCode
	}

	for {
		if int(li) >= len(n) {
			li = NextLangIndex(li)
		} else {
			break
		}
	}

	if li == Unknown {
		if NoFoundIndex != Unknown {
			li = NoFoundIndex
		} else {
			return UnknownCode
		}
	}

	for {
		res := n[li]
		if res != emptyString {
			return res
		}

		li = NextLangIndex(li)
		if li != Unknown {
			continue
		}

		if len(opts) > 0 {
			return opts[0]()
		}

		if NoFoundIndex != Unknown {
			li = NoFoundIndex
			continue
		}
		break
	}

	return NoValue
}

// Bytes returns jsonb representation of the Name.
func (n String) Bytes() []byte {
	var buffer bytes.Buffer

	buffer.WriteString("{")
	sep := ""
	for li, val := range n {
		if val == emptyString {
			continue
		}
		buffer.WriteString(sep)
		buffer.WriteString(`"`)
		buffer.WriteString(Code(LangIndex(li)))
		buffer.WriteString(`":"`)
		buffer.WriteString(val)
		buffer.WriteString(`"`)
		sep = ","
	}
	buffer.WriteString("}")
	return buffer.Bytes()
}

// Value implements interface sql.Valuer
func (n String) Value() (driver.Value, error) {
	return n.MarshalJSON()
}

// Scan implements database/sql Scanner interface.
func (n *String) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Name.Scan: expected []byte, got %T (%q)", value, value)
	}

	var err error
	*n, err = ToString(v)
	return err
}

// ToString decodes jsonb like `{"en":"Name","cz":"Jméno","sr":"Име"}` into String type.
func ToString(b []byte) (String, error) {

	var parsed map[string]string
	err := json.Unmarshal(b, &parsed)
	if err != nil {
		return String{}, err
	}

	n := make(String, Len())
	for i := 0; i < len(n); i++ {
		n[i] = emptyString
	}

	for code, val := range parsed {
		li := int(GetSetLangIndex(code))
		if li >= len(n) {
			x := li - len(n)
			if x < 0 {
				x *= -1
			}
			x++
			for i := 0; i < x; i++ {
				n = append(n, emptyString)
			}
		}
		n[li] = val
	}

	return n, nil
}

// MarshalJSON implements json.Marshaler interface.
func (n String) MarshalJSON() ([]byte, error) {
	if len(n) == 0 {
		return []byte("null"), nil
	}

	return n.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (n *String) UnmarshalJSON(buf []byte) error {
	name, err := ToString(buf)
	if err != nil {
		return err
	}

	*n = name
	return nil
}

// StringValidator returns function that validates jsonb data for String type.
func StringValidator() func([]byte) bool {

	codes := LangCodes()
	validCodes := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		validCodes[code] = struct{}{}
	}

	return func(b []byte) bool {

		if !json.Valid(b) {
			return false
		}

		data := make(map[string]any)
		err := json.Unmarshal(b, &data)
		if err != nil {
			return false
		}

		// check if all keys are valid language codes.
		// check if all values are strings.
		for key, value := range data {
			// check if key is one of language codes.
			if _, ok := validCodes[key]; !ok {
				return false
			}

			if _, ok := value.(string); !ok {
				return false
			}
		}

		return true
	}
}
