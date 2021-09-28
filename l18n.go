package l18n

import (
	"fmt"
	"strings"
)

type (
	L18n struct {
		scopes map[string]interface{}
	}
	Scoped struct {
		scopes map[string]interface{}
	}
)

func L(langs []string) *L18n {
	scopes := map[string]interface{}{}
	for _, l := range langs {
		scopes[l] = map[string]interface{}{}
	}
	return &L18n{scopes}
}
func (l *L18n) Lang(lang string) (scoped *Scoped, err error) {
	rawScopes, ok := l.scopes[lang]
	if !ok {
		return nil, fmt.Errorf("language \"%s\" is not registered", lang)
	}
	scopes, ok := rawScopes.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("wrong type of language \"%s\"scope", lang)
	}
	return &Scoped{scopes}, nil
}
func (l *L18n) Add(path []string, values map[string]interface{}) (err error) {
	var (
		finalScopeIdx = len(path) - 1
		key           = path[len(path)-1]
		scope         = l.scopes
		ok            bool
		rawScope      interface{}
	)
	for registeredLang := range l.scopes {
		_, ok = values[registeredLang]
		if !ok {
			return fmt.Errorf("translation for \"%s\" language is not provided", registeredLang)
		}
	}
	for lang, translation := range values {
		rawScope, ok = l.scopes[lang]
		if !ok {
			return fmt.Errorf("language \"%s\" is not registered", lang)
		}
		scope, ok = rawScope.(map[string]interface{})
		if !ok {
			return fmt.Errorf("")
		}
		for i, chunk := range path[:finalScopeIdx] {
			rawScope, ok = scope[chunk]
			if ok {
				scope, ok = rawScope.(map[string]interface{})
				if !ok {
					return fmt.Errorf(
						"wrong scope type (path: %s), expected: map[string]interface{}, got: %#v",
						strings.Join(
							path[:i+1],
							"/"),
						scope)
				}
			} else {
				s := map[string]interface{}{}
				scope[chunk] = s
				scope = s
			}
		}
		rawScope, ok = scope[key]
		if ok {
			_, ok = rawScope.(map[string]interface{})
			if ok {
				return fmt.Errorf("wrong path \"%s\": there is already exists a scope", strings.Join(path, "/"))
			}
			return fmt.Errorf("translation \"%s\" already exists", strings.Join(path, "/"))
		}
		scope[key] = translation
	}
	return
}

func (s *Scoped) Get(path []string, args map[string]interface{}) (v string, err error) {
	var (
		key       = path[len(path)-1]
		scopes    = s.scopes
		rawScopes interface{}
		ok        bool
	)
	for _, chunk := range path[:len(path)-1] {
		rawScopes, ok = scopes[chunk]
		if !ok {
			return "", fmt.Errorf("translation \"%s\" does not exist", strings.Join(path, "/"))
		}
		scopes, ok = rawScopes.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("translation \"%s\" does not exist", strings.Join(path, "/"))
		}
	}
	rawTranslation := scopes[key]
	switch translation := rawTranslation.(type) {
	case string:
		return translation, nil
	case func(map[string]interface{}) (string, error):
		return translation(args)
	}
	return "", fmt.Errorf("not found")
}
