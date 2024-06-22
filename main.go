package main

import (
	"__ns__/src/app"
	"__ns__/src/commands"
	"__ns__/src/config"
	"os"

	// <%if eq .web "y"%>
	"__ns__/src/http"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gomig/http/middlewares"
	"github.com/gomig/logger" // <%end%>

	// <%if oneOf .database "mysql|postgres" %>
	"github.com/gomig/database/migration" // <%end%>
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

	// <%if eq .database "mysql"%>
	// Config MySQL
	app.SetupMySQL()
	defer app.MySQL().Close()
	app.CLI().AddCommand(migration.MigrationCommand(app.DatabaseResolver, "--APP-DB", "./database/migrations", "./database/seeds"))
	// <% else if eq .database "postgres"%>
	// Config Postgres
	app.SetupPostgres()
	defer app.Postgres().Close()
	app.CLI().AddCommand(migration.MigrationCommand(app.DatabaseResolver, "--APP-DB", "./database/migrations", "./database/seeds"))
	// <% else if eq .database "mongo"%>
	// Config Mongo
	app.SetupMongoDB()
	ctx, cancel := app.MongoOperationCtx()
	defer cancel()
	defer app.MongoClient().Disconnect(ctx)
	// <%end%>

	// <%if eq .web "y"%>
	// Config Web
	app.SetupWeb(http.OnError)
	app.Server().Use(recover.New())
	if app.Config().Cast("web.log").BoolSafe(false) {
		appName := app.Config().Cast("name").StringSafe("<% .name %>")
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
	app.CLI().AddCommand(commands.ServeCommand)
	// <%end%>

	// <% if eq .cache "file" %>
	// Setup cache
	app.CLI().AddCommand(commands.CleanupCommand) // <% end %>

	// Register base commands and run app
	app.CLI().AddCommand(commands.HashCommand(app.CryptoResolver, "--APP-CRYPTO"))
	app.CLI().AddCommand(commands.ClearCommand)
	app.CLI().AddCommand(commands.DownCommand)
	app.CLI().AddCommand(commands.UpCommand)
	app.CLI().AddCommand(commands.VersionCommand)
	app.Run()
}
