package app

import (
	"github.com/gomig/translator"
)

// SetupTranslator driver
func SetupTranslator() {
	conf := confOrPanic()
	appLocale := conf.Cast("locale").StringSafe("// {{.locale}}")

	// {{if eq .translator "json"}}
	if t, err := translator.NewJSONTranslator(appLocale, ConfigPath("strings")); err == nil {
		_container.Register("--APP-TRANSLATOR", t)
	} else {
		panic("failed to build json translator driver")
	}
	// {{else}}
	if t := translator.NewMemoryTranslator(appLocale); t != nil {
		_container.Register("--APP-TRANSLATOR", t)
	} else {
		panic("failed to build translator driver")
	}
	// {{end}}
}

// Translator get translator driver
// leave name empty to resolve default
func Translator(names ...string) translator.Translator {
	name := "--APP-TRANSLATOR"
	if len(names) > 0 {
		name = names[0]
	}
	if dep, exists := _container.Resolve(name); exists {
		if res, ok := dep.(translator.Translator); ok {
			return res
		}
	}
	return nil
}
