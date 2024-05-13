package main

import (
	"github.com/spf13/cobra"

	"gitlab.com/coreflux-cloud/cfctl.git/cmd"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cfctl",
		Short: "cfctl is a CLI tool for managing Vault configurations",
	}

	cmd.InitVaultCommands(rootCmd)

	rootCmd.Execute()
}
