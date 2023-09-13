package commands

import (
	"fmt"
	"mekramy/__boiler/src/app"

	"github.com/spf13/cobra"
)

// ServeCommand serve web app
var ServeCommand = &cobra.Command{
	Use:   "serve",
	Short: "start web server [dev|prod]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if args[0] == "prod" {
				app.GoProd()
				fmt.Println("app mode changed to production")
			} else if args[0] == "dev" {
				app.GoDev()
				fmt.Println("app mode changed to development")
			}
		}
		app.Server().Listen(fmt.Sprintf(":%d", app.Config().Cast("web.port").IntSafe(8888)))
	},
}
