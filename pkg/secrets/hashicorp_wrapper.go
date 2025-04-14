package secrets

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

type VaultSecretManager struct {
	client *vault.Client
}

func NewVaultSecretManager(addr, roleID, secretID string) (*VaultSecretManager, error) {
	config := vault.DefaultConfig()
	config.Address = addr

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	// Assigning Roles
	appRoleAuth, err := auth.NewAppRoleAuth(
		roleID,
		&auth.SecretID{FromString: secretID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create AppRole auth: %w", err)
	}

	// Login using App Role
	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to login with AppRole credentials: %v", err)
	}
	if authInfo == nil || authInfo.Auth == nil || authInfo.Auth.ClientToken == "" {
		return nil, fmt.Errorf("authentication failed: No token received")
	}
	// AuthInfo token will be used automatically by the client for subsequent requests

	return &VaultSecretManager{
		client: client,
	}, nil
}

func (v *VaultSecretManager) GetCryptoKey(path string, keyName string) ([]byte, error) {
	secret, err := v.client.KVv2("secret").Get(context.Background(), path) // Path is 'secret/data/...' for API, but client handles 'data' part
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT public key from Vault: %w", err)
	}

	keyData, ok := secret.Data[keyName].(string)
	if !ok {
		return nil, fmt.Errorf("private key not found or not a string in secret '%s'", path)
	}

	return []byte(keyData), nil
}
