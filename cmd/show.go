package cmd

import (
	"fmt"

	"github.com/CorefluxCommunity/zeusctl/pkg/utils"
	"github.com/CorefluxCommunity/zeusctl/pkg/vault"
	"github.com/CorefluxCommunity/zeusctl/pkg/yubikey"
	"github.com/spf13/cobra"
)

func init() {
	showCmd.AddCommand(showClusterSubCmd)
	showCmd.AddCommand(showYubiKeySubCmd)

	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of YubiKeys and Vault clusters",
	Long:  `Show details of YubiKeys and Vault clusters.`,
}

var showClusterSubCmd = &cobra.Command{
	Use:   "cluster <cluster name>",
	Short: "Show Vault cluster status",
	Long:  `Show overview and unseal status of the specified Vault cluster.`,
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

		utils.PrintHeader("Vault Cluster Status")
		utils.PrintKVSlice("Server(s)", cluster.Servers)

		uniqKeys := utils.Unique(keys)
		if len(keys) != len(uniqKeys) {
			dupCount := len(keys) - len(uniqKeys)
			utils.PrintKV("Key(s)", fmt.Sprintf("%d (%d duplicates)", len(keys), dupCount))
		} else {
			utils.PrintKV("Key(s)", fmt.Sprintf("%d", len(keys)))
		}

		if err := vault.ListVaultStatus(cluster.Servers[0]); err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}

var showYubiKeySubCmd = &cobra.Command{
	Use:   "yubikey <serial number>",
	Short: "Show YubiKey details",
	Long:  `Show YubiKey details returned from OpenPGP application data objects.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sn := args[0]

		err := yubikey.ShowYubiKey(sn)
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}
