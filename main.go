package main

import (
	"github.com/spf13/cobra"

	"github.com/CorefluxCommunity/zeusctl/cmd"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "zeusctl",
		Short: "zeusctl is a CLI tool for operating Hashicorp Vault",
	}

	cmd.InitVaultCommands(rootCmd)

	rootCmd.Execute()
}
