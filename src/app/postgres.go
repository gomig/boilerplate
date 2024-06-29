package app

import (
	"github.com/gomig/database/v2"
	"github.com/jmoiron/sqlx"
)

// SetupPostgres driver
func SetupPostgres() {
	conf := Config()
	host := conf.Cast("database.host").StringSafe("localhost")
	username := conf.Cast("database.username").StringSafe("root")
	password := conf.Cast("database.password").StringSafe("")
	port := conf.Cast("database.port").StringSafe("")
	db := conf.Cast("database.name").StringSafe("<% .name %>")

	if db, err := database.NewPostgresConnector(host, port, username, password, db); err == nil {
		_container.Register("--APP-DB", db)
	} else {
		panic("failed to init postgres database: " + err.Error())
	}
}

// Postgres get database driver
// leave name empty to resolve default
func Postgres(names ...string) *sqlx.DB {
	name := "--APP-DB"
	if len(names) > 0 {
		name = names[0]
	}
	if dep, exists := _container.Resolve(name); exists {
		if res, ok := dep.(*sqlx.DB); ok {
			return res
		}
	}
	return nil
}
