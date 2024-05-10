package contexts

import (
  "fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"

	"gitlab.com/coreflux-cloud/cfctl.git/config" 
)

func NewContextCommand() *cobra.Command {
  contextCmd := &cobra.Command{
	  Use: "context [name]",
		Short: "Load a named context and export its secrets as environment variables",
		Args: cobra.ExactArgs(1),
		Run: loadContext,
	}

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
  contexts, err := LoadContexts()
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
        
        // Output export command instead of setting directly
        fmt.Printf("export %s='%v'\n", envName, value)
      }
      
      break
    }
  }

  if !found {
    fmt.Fprintf(os.Stderr, "Context '%s' not found\n", contextName)
  }
}

