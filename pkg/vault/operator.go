package vault

import (
	"fmt"

	"github.com/CorefluxCommunity/zeusctl/pkg/crypto"
	"github.com/CorefluxCommunity/zeusctl/pkg/utils"
	"github.com/hashicorp/vault/api"
)

// connect to Vault server and execute unseal operation
func (vault *VaultClient) unseal(keys []string) (*api.SealStatusResponse, error) {
	resp, err := vault.ApiClient.Sys().SealStatus()
	if err != nil {
		return nil, err
	}

	if !resp.Initialized {
		return resp, fmt.Errorf("%s - Vault server is not initialized", vault.url.Host)
	}

	// if node is already unsealed, skip it
	if !resp.Sealed {
		utils.PrintSuccess(vault.url.Host + " - already unsealed, skipping unseal operation")
		return resp, nil
	}

	for _, key := range keys {
		resp, err = vault.ApiClient.Sys().Unseal(key)
		if err != nil {
			return nil, err
		}

		if !resp.Sealed {
			break
		}
	}

	utils.PrintInfo(fmt.Sprintf("%s - provided %d unseal key share(s) toward unseal progress", vault.url.Host, len(keys)))

	resp, err = vault.ApiClient.Sys().SealStatus()
	if err != nil {
		return nil, err
	}

	if !resp.Sealed {
		utils.PrintSuccess(fmt.Sprintf("%s - Vault unsealed", vault.url.Host))
	}

	return resp, nil
}

// connect to Vault server and execute unseal operation
func (vault *VaultClient) generateRoot(keys []string) (*api.GenerateRootStatusResponse, error) {
	resp, err := vault.ApiClient.Sys().GenerateRootStatus()
	if err != nil {
		return nil, err
	}

	// if node is already unsealed, skip it
	if !resp.Started {
		utils.PrintWarning(vault.url.Host + " - root token generation process has not been started")
		return resp, nil
	}

	nonce := resp.Nonce
	for _, key := range keys {
		resp, err = vault.ApiClient.Sys().GenerateRootUpdate(key, nonce)
		if err != nil {
			return nil, err
		}

		msg := fmt.Sprintf("%s - provided unseal key share, root token generation progress: %d of %d key shares",
			vault.url.Host, resp.Progress, resp.Required)
		utils.PrintInfo(msg)

		if resp.Complete {
			msg = fmt.Sprintf("%s - root token generation complete", vault.url.Host)
			utils.PrintSuccess(msg)

			return resp, nil
		}
	}

	return resp, nil
}

func printSealStatus(resp *api.SealStatusResponse) {
	status := "unsealed"
	if resp.Sealed {
		status = "sealed"
	} else {
		utils.PrintKV("Cluster name", resp.ClusterName)
		utils.PrintKV("Cluster ID", resp.ClusterID)
	}

	utils.PrintKV("Seal status", status)
	utils.PrintKV("Key threshold/shares", fmt.Sprintf("%d/%d", resp.T, resp.N))
	utils.PrintKV("Progress", fmt.Sprintf("%d/%d", resp.Progress, resp.T))
	utils.PrintKV("Version", resp.Version)
}

func printGenRootStatus(resp *api.GenerateRootStatusResponse) {
	status := "not started"
	if resp.Started {
		status = "started"

		if resp.Complete {
			status = "complete"
		}
	}

	utils.PrintKV("Root generation", status)

	if resp.Started {
		utils.PrintKV("Nonce", resp.Nonce)
		utils.PrintKV("Progress", fmt.Sprintf("%d/%d", resp.Progress, resp.Required))

		if resp.PGPFingerprint != "" {
			utils.PrintKV("PGP fingerprint", resp.PGPFingerprint)
		}
	}

	if resp.EncodedRootToken != "" {
		utils.PrintKV("Encoded root token", resp.EncodedRootToken)
	}
}

// Unseal will decrypt the provided unseal key(s) and unseal each of the
// provided Vault cluster nodes.
func Unseal(vaultAddrs []string, encryptedKeys []string) error {
	keys, err := crypto.DecryptUnsealKeys(encryptedKeys)
	if err != nil {
		return err
	}

	for i, addr := range vaultAddrs {
		vault, err := NewVaultClient(addr)
		if err != nil {
			return err
		}

		resp, err := vault.unseal(keys)
		if err != nil {
			return err
		}

		if i == len(vaultAddrs)-1 {
			fmt.Println()
			utils.PrintHeader("Vault Cluster Status")
			printSealStatus(resp)
		}
	}

	return nil
}

// GenerateRoot will decrypt the provided unseal key and enter the key share
// to progress the root generation attempt.
func GenerateRoot(vaultAddr string, encryptedKeys []string) error {
	keys, err := crypto.DecryptUnsealKeys(encryptedKeys)
	if err != nil {
		return err
	}

	vault, err := NewVaultClient(vaultAddr)
	if err != nil {
		return err
	}

	resp, err := vault.generateRoot(keys)
	if err != nil {
		return err
	}

	fmt.Println()
	utils.PrintHeader("Root Token Generation Status")
	printGenRootStatus(resp)

	return nil
}

// ListVaultStatus will output of the status the provided Vault address.
func ListVaultStatus(vaultAddr string) error {
	vault, err := NewVaultClient(vaultAddr)
	if err != nil {
		return err
	}

	resp, err := vault.ApiClient.Sys().SealStatus()
	if err != nil {
		return err
	}

	printSealStatus(resp)

	return nil
}
