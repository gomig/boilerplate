package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCommand get gomig version
var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "get gomig version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.3.0")
	},
}
