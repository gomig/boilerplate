package commands

import (
	"__ns__/src/app"
	"fmt"

	"github.com/gomig/cache"
	"github.com/spf13/cobra"
)

// CleanupCommand clear expired cache records
var CleanupCommand = &cobra.Command{
	Use:   "cleanup",
	Short: "clear expired cache records",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cleaning expired caches...")
		if err := cache.CleanFileExpiration(app.StoragePath("cache")); err != nil {
			fmt.Println("failed: " + err.Error())
		}
		fmt.Println("done!")
	},
}
