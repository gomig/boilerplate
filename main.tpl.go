package main

import (
	"os"

	"mekramy/__boiler/src/app"
	"mekramy/__boiler/src/commands"
	"mekramy/__boiler/src/config"

	// {{if eq .web "y"}}
	"mekramy/__boiler/src/http"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gomig/http/middlewares"
	"github.com/gomig/logger"

	// {{end}}
	// {{if eq .database "mysql"}}
	"github.com/gomig/database/migration"
	// {{end}}
)

func main() {
	app.SetupConfig()
	config.Configure(app.Config())
	app.SetupCache()
	app.SetupCrypto()
	app.SetupLogger()
	app.SetupTranslator()
	config.ConfigureMessages(app.Translator())
	app.SetupValidator()
	// {{if eq .database "mysql"}}
	app.SetupDatabase()
	defer app.Database().Close() // {{end}}
	// {{if eq .database "mongo"}}
	app.SetupMongoDB()
	ctx, cancel := app.MongoOperationCtx()
	defer cancel()
	defer app.MongoClient().Disconnect(ctx) // {{end}}
	// {{if eq .web "y"}}
	app.SetupWeb(http.OnError)
	app.Server().Use(recover.New())
	if app.Config().Cast("web.log").BoolSafe(false) {
		appName := app.Config().Cast("name").StringSafe("// {{.name}}")
		_logger := logger.NewLogger("2006-01-02 15:04:05", app.DateFormatter())
		_logger.AddWriter("main", logger.NewFileLogger(app.LogPath("access"), appName, "2006-01-02", app.DateFormatter()))
		if !app.IsProd() {
			_logger.AddWriter("dev", os.Stdout)
		}
		app.Server().Use(middlewares.AccessLogger(_logger))
	}
	app.Server().Static("/", app.PublicPath())
	http.RegisterGlobalMiddlewares(app.Server())
	http.RegisterRoutes(app.Server())
	// {{end}}

	// Register commands and run app
	app.CLI().AddCommand(commands.HashCommand(app.CryptoResolver, "--APP-CRYPTO"))
	app.CLI().AddCommand(commands.ClearCommand)
	app.CLI().AddCommand(commands.DownCommand)
	app.CLI().AddCommand(commands.UpCommand)
	app.CLI().AddCommand(commands.VersionCommand)
	// {{ if ne .cache "redis" }}
	app.CLI().AddCommand(commands.CleanupCommand) // {{ end }}
	// {{if eq .database "mysql"}}
	app.CLI().AddCommand(migration.MigrationCommand(app.DatabaseResolver, "--APP-DB", "./database/migrations", "./database/seeds")) // {{end}}
	// {{if eq .web "y"}}
	app.CLI().AddCommand(commands.ServeCommand) // {{end}}
	app.Run()
}
