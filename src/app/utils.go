package app

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gomig/cache"
	"github.com/gomig/config"
	"github.com/gomig/crypto"
	"github.com/gomig/logger"
	"github.com/microcosm-cc/bluemonday"
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
	// <%if eq .locale "fa"%>
	//- return logger.JalaaliFormatter
	// <%else%>
	return logger.GregorianFormatter // <%end%>
}

// IsUnderMaintenance check if under maintenance mode
func IsUnderMaintenance() (bool, error) {
	return cacheOrPanic().Exists("maintenance")
}

// ValueOf get value of pointer or return fallback if value is nil
func ValueOf[T any](value *T, fallback T) T {
	if value == nil {
		return fallback
	} else {
		return *value
	}
}

// PointerOf get pointer of value
func PointerOf[T any](value T) *T {
	return &value
}

// NullableOf return nil if value is empty
func NullableOf[T comparable](v T) *T {
	var empty T
	if v == empty {
		return nil
	}
	return &v
}

// ClearText clear string from all html tags
func ClearText(data string, trim bool) string {
	res := bluemonday.StrictPolicy().Sanitize(data)
	if trim {
		return strings.TrimSpace(res)
	}
	return res
}

// SecureText secure string from script-like html tags
func SecureText(data string, trim bool) string {
	res := bluemonday.UGCPolicy().Sanitize(data)
	if trim {
		return strings.TrimSpace(res)
	}
	return res
}

// FuzzyPhrase generate fuzzy search phrase and numeric value for search queries
func FuzzyPhrase(search string, fallbackNum int64) (string, int64) {
	query := "%" + strings.ReplaceAll(regexp.QuoteMeta(search), " ", "%") + "%"
	if num, _ := strconv.ParseInt(strings.TrimPrefix(search, "#"), 10, 64); num != 0 {
		return query, num
	} else {
		return query, fallbackNum
	}
}
