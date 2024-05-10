package cmd

import (
	"github.com/spf13/cobra"
  
  "gitlab.com/coreflux-cloud/cfctl.git/auth"
  "gitlab.com/coreflux-cloud/cfctl.git/contexts"
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

