package cmd

import (
	"github.com/spf13/cobra"

	"github.com/CorefluxCommunity/zeusctl/pkg/utils"
	"github.com/CorefluxCommunity/zeusctl/pkg/vault"
)

func init() {
	unsealCmd.AddCommand(unsealServerSubCmd)
	unsealCmd.AddCommand(unsealClusterSubCmd)

	rootCmd.AddCommand(unsealCmd)
}

var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal Vault by server or cluster",
	Long:  `Decrypt PGP-encrypted unseal key and unseal Vault.`,
}

var unsealServerSubCmd = &cobra.Command{
	Use:   "server <vault address> <unseal key path>",
	Short: "Unseal Vault server",
	Long:  `Decrypt unseal key and unseal single Vault server.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vaultAddr := args[0]
		keyPath := args[1]

		keys, err := utils.ReadKeyFile(keyPath)
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}

		if err := vault.Unseal([]string{vaultAddr}, keys); err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}

var unsealClusterSubCmd = &cobra.Command{
	Use:   "cluster <cluster name>",
	Short: "Unseal Vault cluster",
	Long:  `Decrypt unseal key and unseal Vault cluster.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]

		cluster, err := getVaultClusterConfig(clusterName)
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}

		keys, err := cluster.keyring()
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}

		if len(cluster.Servers) == 0 {
			utils.PrintFatal("no Vault servers in configuration", 1)
		}

		if err := vault.Unseal(cluster.Servers, utils.Unique(keys)); err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}
