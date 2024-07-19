package contexts

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"

	"github.com/CorefluxCommunity/zeusctl/config"
)

var configFilePath string

func NewContextCommand() *cobra.Command {
	contextCmd := &cobra.Command{
		Use:   "context [name]",
		Short: "Load a named context and export its secrets as environment variables",
		Args:  cobra.ExactArgs(1),
		Run:   loadContext,
	}

	contextCmd.Flags().StringVarP(&configFilePath, "config", "c", "", "Path to the contexts configuration file")

	return contextCmd
}

func loadContext(cmd *cobra.Command, args []string) {
	vaultConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %s\n", err)
		os.Exit(1)
	}

	clientConfig := &api.Config{
		Address: vaultConfig.Host,
	}
	if vaultConfig.CAPath != "" {
		clientConfig.ConfigureTLS(&api.TLSConfig{CAPath: vaultConfig.CAPath})
	}

	client, err := api.NewClient(clientConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Vault client: %s\n", err)
		os.Exit(1)
	}
	client.SetToken(vaultConfig.Token)

	contextName := args[0]
	contexts, err := LoadContexts(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading contexts: %s\n", err)
		os.Exit(1)
	}

	found := false
	for _, ctx := range contexts.Contexts {
		if ctx.Name == contextName {
			found = true
			for envName, secretInfo := range ctx.Secrets {
				secret, err := client.Logical().Read(secretInfo.Path)
				if err != nil {
					fmt.Printf("Error reading secret from %s: %s\n", secretInfo.Path, err)
					continue

				}

				if secret == nil || secret.Data == nil {
					fmt.Printf("No data found at path: %s\n", secretInfo.Path)
					continue
				}

				value, ok := secret.Data["data"].(map[string]interface{})[secretInfo.Key]
				if !ok {

					fmt.Printf("Secret key %s not found at path: %s\n", secretInfo.Key, secretInfo.Path)

					continue
				}

				var finalValue string
				if secretInfo.DecodeBase64 {
					decodedBytes, decodeErr := base64.StdEncoding.DecodeString(fmt.Sprintf("%v", value))
					if decodeErr != nil {
						fmt.Printf("Error decoding base64 value for %s: %s\n", envName, decodeErr)
						continue
					}
					finalValue = string(decodedBytes)
				} else {
					finalValue = fmt.Sprintf("%v", value)
				}
				if secretInfo.EncodeBase64 {
					finalValue = base64.StdEncoding.EncodeToString([]byte(finalValue))
				}

				// Output export command instead of setting directly
				fmt.Printf("export %s='%s'\n", envName, finalValue)
			}

			break
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Context '%s' not found\n", contextName)
	}
}
