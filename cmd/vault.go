package cmd

import (
	"github.com/spf13/cobra"

	"github.com/CorefluxCommunity/zeusctl/auth"
	"github.com/CorefluxCommunity/zeusctl/contexts"
)

func InitVaultCommands(rootCmd *cobra.Command) {
	vaultCmd := &cobra.Command{
		Use:   "vault",
		Short: "Commands related to Vault",
	}

	vaultCmd.AddCommand(auth.NewAuthCommand())
	vaultCmd.AddCommand(contexts.NewContextCommand())

	rootCmd.AddCommand(vaultCmd)
}
