package auth

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"

	"github.com/CorefluxCommunity/zeusctl/config"
)

var (
	host     string
	caPath   string
	user     string
	password string
)

func NewAuthCommand() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with Vault",
		Run:   authenticate,
	}

	authCmd.Flags().StringVar(&host, "host", "", "Vault host URL")
	authCmd.Flags().StringVar(&caPath, "ca-path", "", "Path to CA certificate")
	authCmd.Flags().StringVar(&user, "user", "", "Username for Vault login")
	authCmd.Flags().StringVar(&password, "password", "", "Password for Vault login")

	return authCmd
}

func authenticate(cmd *cobra.Command, args []string) {
	apiConfig := &api.Config{
		Address: host,
	}

	if caPath != "" {
		apiConfig.ConfigureTLS(&api.TLSConfig{CAPath: caPath})
	}

	client, err := api.NewClient(apiConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Vault client: %s\n", err)
		os.Exit(1)
	}

	options := map[string]interface{}{
		"password": password,
	}
	path := fmt.Sprintf("auth/userpass/login/%s", user)

	secret, err := client.Logical().Write(path, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error logging in to Vault: %s\n", err)
		os.Exit(1)
	}

	// Save the configuration to a file
	err = config.SaveConfig(secret.Auth.ClientToken, host, caPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving configuration: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Logged in successfully. Configuration saved.")
}
