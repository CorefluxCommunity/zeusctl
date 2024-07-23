package cmd

import (
	"github.com/spf13/cobra"

	"github.com/CorefluxCommunity/zeusctl/pkg/utils"
	"github.com/CorefluxCommunity/zeusctl/pkg/vault"
)

func init() {
	generateRootServerSubCmd.Flags().StringVarP(&vaultGenerateRootNonce, "nonce", "n", "", "nonce for root token generation")
	generateRootClusterSubCmd.Flags().StringVarP(&vaultGenerateRootNonce, "nonce", "n", "", "nonce for root token generation")

	generateRootCmd.AddCommand(generateRootServerSubCmd)
	generateRootCmd.AddCommand(generateRootClusterSubCmd)

	rootCmd.AddCommand(generateRootCmd)
}

var generateRootCmd = &cobra.Command{
	Use:   "generate-root",
	Short: "Generate Vault root token",
	Long:  `Decrypt the unseal key and generate root token for Vault cluster.`,
}

var generateRootServerSubCmd = &cobra.Command{
	Use:   "server <vault address> <unseal key path> -n <nonce>",
	Short: "Generate root token for Vault cluster",
	Long:  `Decrypt the unseal key and generate Vault root token.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vaultAddr := args[0]
		keyPath := args[1]

		keys, err := utils.ReadKeyFile(keyPath)
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}

		if err := vault.GenerateRoot(vaultAddr, keys); err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}

var generateRootClusterSubCmd = &cobra.Command{
	Use:   "cluster <cluster name> -n <nonce>",
	Short: "Generate root token for Vault cluster",
	Long:  `Decrypt the unseal key and generate Vault root token.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]

		vaultAddr, err := getVaultAddress(clusterName)
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}

		cluster, err := getVaultClusterConfig(clusterName)
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}

		keys, err := cluster.keyring()
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}

		if err := vault.GenerateRoot(vaultAddr, utils.Unique(keys)); err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}
