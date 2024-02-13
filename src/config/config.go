package config

import (
	// <%if eq .database "mongo"%>
	"fmt" // <% end %>

	"github.com/gomig/config"
)

// Configure register/override app config
func Configure(config config.Config) {
	config.Set("name", "<% .name %>")
	config.Set("locale", "<% .locale %>")
	config.Set("key", "<% .key %>")
	config.Set("mc_key", "<% key64 %>")
	// <%if eq .database "mongo"%>
	rs := config.Cast("database.replica_set").StringSafe("")
	if rs != "" {
		config.Set("mongo.conStr", fmt.Sprintf("mongodb://127.0.0.1:%v/?directConnection=true&replicaSet=%s", config.Get("database.port"), rs))
	} else {
		config.Set("mongo.conStr", fmt.Sprintf("mongodb://127.0.0.1:%v/?directConnection=true", config.Get("database.port")))
	} // <%else if eq .database "mysql"%>
	config.Set("database.host", "")
	config.Set("database.username", "root")
	config.Set("database.password", "root")
	// <%end%>
}
