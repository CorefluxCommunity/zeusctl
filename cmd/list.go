package cmd

import (
	"fmt"

	"github.com/CorefluxCommunity/zeusctl/pkg/utils"
	"github.com/CorefluxCommunity/zeusctl/pkg/yubikey"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listClustersSubCmd)
	listCmd.AddCommand(listYubiKeysSubCmd)

	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List connected YubiKeys and configured Vault clusters",
	Long:  `List connected YubiKeys and configured Vault clusters.`,
}

var listClustersSubCmd = &cobra.Command{
	Use:   "clusters",
	Short: "List Vault clusters",
	Long:  `List Vault clusters in Zeus configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		i := 0
		for name, cluster := range config.Clusters {
			keys, err := cluster[0].keyring()
			if err != nil {
				utils.PrintFatal(err.Error(), 1)
			}

			utils.PrintHeader(name)
			utils.PrintKVSlice("Server(s)", cluster[0].Servers)

			uniqKeys := utils.Unique(keys)
			if len(keys) != len(uniqKeys) {
				dupCount := len(keys) - len(uniqKeys)
				utils.PrintKV("Key(s)", fmt.Sprintf("%d (%d duplicates)", len(keys), dupCount))
			} else {
				utils.PrintKV("Key(s)", fmt.Sprintf("%d", len(keys)))
			}

			if i < len(config.Clusters)-1 {
				fmt.Println()
			}

			i++
		}
	},
}

var listYubiKeysSubCmd = &cobra.Command{
	Use:   "yubikeys",
	Short: "List connected YubiKeys",
	Long:  `List overview of connected YubiKeys.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := yubikey.ListYubiKeys()
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}
