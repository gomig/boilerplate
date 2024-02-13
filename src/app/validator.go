package app

import (
	"github.com/gomig/validator"
	// <%if eq .locale "fa"%>
	"github.com/gomig/validator/validations" // <%end%>
)

// SetupValidator driver
func SetupValidator() {
	conf := confOrPanic()
	appLocale := conf.Cast("locale").StringSafe("<% .locale %>")
	if v := validator.NewValidator(Translator(), appLocale); v != nil {
		// <%if eq .locale "fa"%>
		validations.RegisterExtraValidations(v) // <%end%>
		_container.Register("--APP-VALIDATOR", v)
	} else {
		panic("failed to build validator driver")
	}
}

// Validator get validator driver
// leave name empty to resolve default
func Validator(names ...string) validator.Validator {
	name := "--APP-VALIDATOR"
	if len(names) > 0 {
		name = names[0]
	}
	if dep, exists := _container.Resolve(name); exists {
		if res, ok := dep.(validator.Validator); ok {
			return res
		}
	}
	return nil
}
