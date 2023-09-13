package app

import (
	"github.com/gomig/crypto"
)

// SetupCrypto driver
func SetupCrypto() {
	conf := confOrPanic()
	if k, err := conf.Cast("key").String(); err == nil {
		if k == "" {
			panic("app key cannot be empty")
		}
		if c := crypto.NewCryptography(k); c != nil {
			_container.Register("--APP-CRYPTO", c)
		} else {
			panic("failed to build crypto driver")
		}
	} else {
		panic(err)
	}
}

// Crypto get crypto driver
// leave name empty to resolve default
func Crypto(names ...string) crypto.Crypto {
	name := "--APP-CRYPTO"
	if len(names) > 0 {
		name = names[0]
	}
	if dep, exists := _container.Resolve(name); exists {
		if res, ok := dep.(crypto.Crypto); ok {
			return res
		}
	}
	return nil
}
