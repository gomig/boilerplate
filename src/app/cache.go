package app

import (
	"github.com/gomig/cache"

	// {{if eq .cache "redis"}}
	"github.com/go-redis/redis/v8" // {{end}}
)

// SetupCache driver
func SetupCache() {
	// {{if eq .cache "redis"}}
	conf := confOrPanic()
	appName := conf.Cast("name").StringSafe("{{ .name }}")
	host := conf.Cast("redis.host").StringSafe("localhost:6379")
	db := conf.Cast("redis.host").IntSafe(0)
	if c := cache.NewRedisCache(appName, redis.Options{Addr: host, DB: db}); c != nil {
		_container.Register("--APP-CACHE", c)
	} else {
		panic("failed to build redis cache driver")
	} // {{else if eq .cache "file"}}
	if c := cache.NewFileCache(appName, StoragePath("cache")); c != nil {
		_container.Register("--APP-CACHE", c)
	} else {
		panic("failed to build file cache driver")
	} // {{end}}
}

// Cache get cache manager
// leave name empty to resolve default
func Cache(names ...string) cache.Cache {
	name := "--APP-CACHE"
	if len(names) > 0 {
		name = names[0]
	}
	if dep, exists := _container.Resolve(name); exists {
		if res, ok := dep.(cache.Cache); ok {
			return res
		}
	}
	return nil
}
