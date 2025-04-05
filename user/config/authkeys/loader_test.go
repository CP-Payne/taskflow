package authkeys_test

import (
	"testing"

	"github.com/CP-Payne/taskflow/user/config"
	"github.com/CP-Payne/taskflow/user/config/authkeys"
)

func TestAuthKeys_Load(t *testing.T) {
	cfg := config.New("../.env")

	if cfg.KeyPath == "" {
		t.Errorf("failed to load configuration")
	}

	authKeys := authkeys.NewAuthKeys()
	err := authKeys.Load(cfg.KeyPath)
	if err != nil {
		t.Errorf("Load, expected nil, received: %v", err)
	}

	if authKeys.PrivateKey() == nil {
		t.Error("Load, expected PrivateKey() to not be nil")
	}

	if authKeys.PublicKey() == nil {
		t.Error("Load, expected PublicKey() to not be nil")
	}
}
