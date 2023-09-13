package commands

import (
	"fmt"
	"mekramy/__boiler/src/app"

	"github.com/spf13/cobra"
)

// UpCommand exit app from maintenance mode
var UpCommand = &cobra.Command{
	Use:   "up",
	Short: "set app status to active mode",
	Run: func(cmd *cobra.Command, args []string) {
		cache := app.Cache()
		if cache == nil {
			fmt.Println("failed: app cache driver not found!")
			return
		}
		cache.Forget("maintenance")
		fmt.Println("app is active!")
	},
}
