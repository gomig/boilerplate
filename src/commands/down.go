package commands

import (
	"fmt"
	"mekramy/__boiler/src/app"

	"github.com/spf13/cobra"
)

// DownCommand put app to maintenance mode
var DownCommand = &cobra.Command{
	Use:   "down",
	Short: "set app status to maintenance mode",
	Run: func(cmd *cobra.Command, args []string) {
		cache := app.Cache()
		if cache == nil {
			fmt.Println("failed: app cache driver not found!")
			return
		}
		cache.PutForever("maintenance", true)
		fmt.Println("app is under maintenance mode!")
	},
}
