package main

import (
	"os"

	"github.com/CP-Payne/taskflow/pkg/authkeys"
	"github.com/CP-Payne/taskflow/pkg/secrets"
	"github.com/CP-Payne/taskflow/user/config"
	"github.com/CP-Payne/taskflow/user/internal/auth"
	"github.com/CP-Payne/taskflow/user/internal/repository/memory"
	"github.com/CP-Payne/taskflow/user/internal/server"
	"github.com/CP-Payne/taskflow/user/internal/service"
	"go.uber.org/zap"
)

func main() {
	_ = config.New("user/.env")

	logger := zap.Must(zap.NewDevelopment()).Sugar()
	defer logger.Sync()

	vaultAddr := os.Getenv("VAULT_ADDR")
	roleID := os.Getenv("APPROLE_ROLE_ID")
	secretID := os.Getenv("APPROLE_SECRET_ID")

	secretsManager, err := secrets.NewVaultSecretManager(vaultAddr, roleID, secretID)
	if err != nil {
		logger.Fatalf("failed to create secrets manager: %v", err)
	}

	keyPath := os.Getenv("VAULT_KEY_PATH")
	keyName := os.Getenv("VAULT_KEY_NAME")

	authKeys := authkeys.NewAuthKeys()
	err = authKeys.LoadFromSecretsManager(secretsManager, keyPath, keyName, authkeys.Private)
	if err != nil {
		logger.Fatalf("failed to fetch auth key: %v", err)
	}

	// logger.Debugf("Fetched Private key from Vault: \n%s\n", authKeys.PrivateKey())
	authenticator := auth.NewJWTAuthenticator(authKeys)

	repo := memory.NewInMemory()
	srv := service.New(repo, authenticator, logger)
	server.StartGRPCServer("localhost:3033", srv, logger)
}
