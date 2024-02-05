package app

import (
	"github.com/gomig/cache"
	"github.com/gomig/config"
	"github.com/gomig/crypto"
	"github.com/gomig/logger"

	// {{if eq .database "mysql"}}
	"github.com/jmoiron/sqlx" // {{end}}
)

func confOrPanic() config.Config {
	if Config() == nil {
		panic("Failed to find default config")
	}
	return Config()
}

func cacheOrPanic() cache.Cache {
	if Cache() == nil {
		panic("Failed to find default cache")
	}
	return Cache()
}

// CryptoResolver resolve crypto driver by name
func CryptoResolver(driver string) crypto.Crypto {
	return Crypto(driver)
}

// CacheResolver resolve cache driver by name
func CacheResolver(driver string) cache.Cache {
	return Cache(driver)
}

// DateFormatter get default app date formatter
func DateFormatter() logger.TimeFormatter {
	// {{if eq .locale "fa"}}
	//- return logger.JalaaliFormatter
	// {{else}}
	return logger.GregorianFormatter // {{end}}
}

// IsUnderMaintenance check if under maintenance mode
func IsUnderMaintenance() (bool, error) {
	return cacheOrPanic().Exists("maintenance")
}

// {{if eq .database "mysql"}}
// DatabaseResolver resolve database driver by name
func DatabaseResolver(driver string) *sqlx.DB {
	return Database(driver)
} // {{end}}
