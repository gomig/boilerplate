package commands

import (
	"__ns__/src/app"
	"fmt"
	"path/filepath"

	"github.com/gomig/utils"
	"github.com/spf13/cobra"
)

// ClearCommand clear app logs
var ClearCommand = &cobra.Command{
	Use:   "clear [Directory name or all for clear anything]",
	Short: "clear logs directory",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "all" {
			dirs, err := utils.GetSubDirectory(app.LogPath())
			if err != nil {
				fmt.Printf("failed: %s\n", err.Error())
				return
			}
			for _, dir := range dirs {
				if err := utils.ClearDirectory(filepath.Join(app.LogPath(), dir)); err != nil {
					fmt.Printf("failed: %s\n", err.Error())
					return
				}
			}
		} else {
			if isDir, err := utils.IsDirectory(filepath.Join(app.LogPath(), args[0])); err != nil {
				fmt.Printf("failed: %s\n", err.Error())
				return
			} else if !isDir {
				fmt.Printf("failed: %s log directory not found\n", args[0])
				return
			}

			if err := utils.ClearDirectory(filepath.Join(app.LogPath(), args[0])); err != nil {
				fmt.Printf("failed: %s\n", err.Error())
				return
			}
		}
		fmt.Printf("cleared!\n")
	},
}
