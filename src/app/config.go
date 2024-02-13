package app

import "github.com/gomig/config"

// SetupConfig driver
func SetupConfig() {
	// <%if eq .config "env"%>
	if c, err := config.NewEnvConfig(ConfigPath("config.env")); err == nil {
		_container.Register("--APP-CONFIG", c)
	} else {
		panic(err)
	} // <%else if eq .config "json"%>
	if c, err := config.NewJSONConfig(ConfigPath("config.json")); err == nil {
		_container.Register("--APP-CONFIG", c)
	} else {
		panic(err)
	} // <%else if eq .config "memory"%>
	if c, err := config.NewMemoryConfig(nil); err == nil {
		_container.Register("--APP-CONFIG", c)
	} else {
		panic(err)
	} // <%end%>
}

// Config get config manager
// leave name empty to resolve default
func Config(names ...string) config.Config {
	name := "--APP-CONFIG"
	if len(names) > 0 {
		name = names[0]
	}
	if dep, exists := _container.Resolve(name); exists {
		if res, ok := dep.(config.Config); ok {
			return res
		}
	}
	return nil
}
