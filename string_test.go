package i18n

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	resetLangState()

	en := GetSetLangIndex("en")
	fr := GetSetLangIndex("fr")
	es := GetSetLangIndex("es")

	data := []byte(`{"en":"Hello","fr":"Bonjour"}`)
	n, err := ToString(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("InLang", func(t *testing.T) {

		if got := n.InLang(es, WithDefault("Hola")); got != "Hola" {
			t.Errorf("expected '%s', got '%s'", "Hola", got)
		}

		if got := n.InLang(en); got != "Hello" {
			t.Errorf("expected 'Hello', got '%s'", got)
		}

		if got := n.InLang(fr); got != "Bonjour" {
			t.Errorf("expected 'Bonjour', got '%s'", got)
		}

		if got := n.InLang(Unknown); got != UnknownCode {
			t.Errorf("expected '%s', got '%s'", UnknownCode, got)
		}

		if got := n.InLang(100); got != UnknownCode {
			t.Errorf("expected '%s', got '%s'", UnknownCode, got)
		}

		if got := n.InLang(100, WithDefault("Hi")); got != UnknownCode {
			t.Errorf("expected '%s', got '%s'", UnknownCode, got)
		}

		if got := n.InLang(GetSetLangIndex("en-US")); got != "Hello" {
			t.Errorf("expected '%s', got '%s'", "Hello", got)
		}

	})

	t.Run("Bytes", func(t *testing.T) {
		expected := string(data)
		if got := string(n.Bytes()); got != expected {
			t.Errorf("expected '%s', got '%s'", expected, got)
		}
	})

	t.Run("Value", func(t *testing.T) {
		expected := string(data)
		if got, err := n.Value(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if string(got.([]byte)) != expected {
			t.Errorf("expected '%s', got '%s'", expected, got)
		}
	})

	t.Run("Scan", func(t *testing.T) {
		n := String{}
		if err := n.Scan(data); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := n.InLang(fr); got != "Bonjour" {
			t.Errorf("expected 'Bonjour', got '%s'", got)
		}
	})

	t.Run("StringValidator", func(t *testing.T) {
		validator := StringValidator()
		if !validator(data) {
			t.Errorf("expected valid JSON to pass validation")
		}
		invalidData := []byte(`{"en":"Hello","unknown":1}`)
		if validator(invalidData) {
			t.Errorf("expected invalid JSON to fail validation")
		}
	})

	t.Run("ToString", func(t *testing.T) {
	})

}

func ExampleFileLocalStorage() {

	// Init part

	_ = GetSetLangIndex("en-US")
	en := GetSetLangIndex("en")
	de := GetSetLangIndex("de")

	fs := NewLocalFileStorage()
	if err := fs.RegisterFiles("*.t18n", "./testdata"); err != nil {
		fmt.Printf("file registration error: %v", err)
		return
	}

	tc := NewContainer(
		WithPrimaryLanguage(en),
		WithStorage(fs),
		WithBrackets("%"),
	)

	if err := tc.ReadRegisteredFiles(); err != nil {
		fmt.Printf("file registration error: %v", err)
		return
	}

	// API Request processing part
	languageCodeFromHTTPRequest := "en-GB"
	requestLi := Index(languageCodeFromHTTPRequest)

	// get translation for the language code from the request
	tr := tc.Lang(requestLi)

	trn := tc.Namespace("customer1", Index("en-US"))

	fmt.Println(tr.Value("%Lift%"))          // from ./testdata/en-GB.t18n
	fmt.Println(tr.Value("%Save%"))          // from ./testdata/en.t18n
	fmt.Println(tc.Lang(de).Value("%Save%")) // from ./testdata/de.t18n
	fmt.Println(trn.Value("%Lift%"))         // from ./testdata/en-US.customer1.t18n
	// Output:
	// Elevator
	// Save
	// Speichern
	// Hoist
}
