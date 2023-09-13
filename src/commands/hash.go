package commands

import (
	"fmt"

	"github.com/gomig/crypto"
	"github.com/spf13/cobra"
)

// HashCommand get hash cli command
func HashCommand(resolver func(driver string) crypto.Crypto, defDriver string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "hash [String to hash] [Algorithm]"
	cmd.Short = "hash string using registered cryptography driver"
	cmd.Long = "available algo is: MD4, MD5, SHA1, SHA256, SHA256224, SHA384, SHA512, SHA512224, SHA512256, SHA3224, SHA3256, SHA3384, SHA3512, KECCAK256, KECCAK512"
	cmd.Args = cobra.MinimumNArgs(2)
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var err error
		driver, err := cmd.Flags().GetString("driver")
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		cryptoDriver := resolver(driver)
		if cryptoDriver == nil {
			fmt.Printf("failed: %s crypto driver not found\n", driver)
			return
		}

		var algo crypto.HashAlgo = 0
		if !algo.Parse(args[1]) {
			fmt.Println("invalid algorithm. see full command description")
			return
		}

		res, err := cryptoDriver.Hash(args[0], algo)
		if err != nil {
			fmt.Printf("failed! %s\n", err.Error())
			return
		}

		fmt.Printf("%s\n", res)
	}
	cmd.Flags().StringP("driver", "d", defDriver, "crypto driver name")
	return cmd
}
