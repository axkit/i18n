package i18n

import (
	"encoding/json"
	"errors"
	"strings"
)

// TranslationRequest is a request for a translation container in a specific language.
type TranslationRequest struct {
	lang      LangIndex
	namespace string
	tc        *TranslationContainer
}

// Lang returns a TranslationRequest in a specific language and default namespace.
func (tc *TranslationContainer) Lang(li LangIndex) TranslationRequest {
	return TranslationRequest{
		lang: li,
		tc:   tc,
	}
}

// Namespace returns a TranslationRequest in a specific namespace.
func (tr TranslationRequest) Namespace(namespace string) TranslationRequest {
	tr.namespace = namespace
	return tr
}

func (tr TranslationRequest) value(key string) Item {

	for {
		if res, ok := tr.item(key); ok {
			return res
		}

		if tr.lang = NextLangIndex(tr.lang); tr.lang == Unknown {
			break
		}
	}

	switch tr.tc.cfg.strategy {
	case ReturnResourceCode:
		return Item{Key: key, Value: key}
	case ReturnEmptyString:
		return Item{Key: key, Value: ""}
	}

	return Item{Key: key, Value: NotFoundMarker}
}

// Value returns a translation value for a specific key.
func (tr TranslationRequest) Value(key string) string {
	return tr.value(key).Value
}

// Hint returns a hint for a specific key.
func (tr TranslationRequest) Hint(key string) string {
	return tr.value(key).Hint
}

func (tr TranslationRequest) item(id string) (Item, bool) {

	if len(id) > 2 &&
		tr.tc.cfg.bracketSymbol != "" &&
		strings.HasPrefix(id, tr.tc.cfg.bracketSymbol) &&
		strings.HasSuffix(id, tr.tc.cfg.bracketSymbol) {
		id = id[1 : len(id)-1]
	}

	rsi, ok := tr.tc.translations[key{lang: tr.lang, namespace: tr.namespace}]
	if ok {
		var idx int
		if idx, ok = rsi.index[id]; ok {
			return rsi.items[idx], true
		}
	}

	if tr.tc.cfg.primaryLanguage == Unknown {
		return Item{}, false
	}

	rsi, ok = tr.tc.translations[key{lang: tr.tc.cfg.primaryLanguage, namespace: tr.namespace}]
	if ok {
		var idx int
		if idx, ok = rsi.index[id]; ok {
			return rsi.items[idx], true
		}
	}
	return Item{}, false
}

func (tr TranslationRequest) ValueWithDefault(id string, notFoundValue string) string {
	res, ok := tr.item(id)
	if !ok {
		return notFoundValue
	}
	return res.Value
}

// JSON returns translation in JSON format.
func (tr TranslationRequest) JSON() ([]byte, error) {
	kv := make(map[string]ResponseItem)
	set, ok := tr.tc.translations[key{lang: tr.lang}]
	if !ok {
		if tr.tc.cfg.primaryLanguage != Unknown {
			set, ok = tr.tc.translations[key{lang: tr.tc.cfg.primaryLanguage}]
		}
	}
	if !ok {
		return nil, errors.New("no translation found")
	}

	for _, item := range set.items {
		kv[tr.tc.genKey(item.Key)] = ResponseItem{
			Value: item.Value,
			Hint:  item.Hint,
		}
	}

	if tr.tc.cfg.primaryLanguage != Unknown && tr.tc.cfg.primaryLanguage != tr.lang {
		if set, ok = tr.tc.translations[key{lang: tr.tc.cfg.primaryLanguage}]; ok {
			for _, item := range set.items {
				k := tr.tc.genKey(item.Key)
				if _, ok := kv[k]; ok {
					continue
				}
				kv[k] = ResponseItem{
					Value: item.Value,
					Hint:  item.Hint,
				}
			}
		}
	}

	buf, err := json.Marshal(kv)
	return buf, err
}
