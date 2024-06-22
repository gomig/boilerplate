package app

import (
	"github.com/gomig/database"
	"github.com/jmoiron/sqlx"
)

// SetupMySQL driver
func SetupMySQL() {
	conf := Config()
	host := conf.Cast("database.host").StringSafe("localhost")
	username := conf.Cast("database.username").StringSafe("root")
	password := conf.Cast("database.password").StringSafe("")
	db := conf.Cast("database.name").StringSafe("<% .name %>")

	if db, err := database.NewMySQLConnector(host, username, password, db); err == nil {
		_container.Register("--APP-DB", db)
	} else {
		panic("failed to init mysql database: " + err.Error())
	}
}

// MySQL get database driver
// leave name empty to resolve default
func MySQL(names ...string) *sqlx.DB {
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
