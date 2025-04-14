package authkeys

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CP-Payne/taskflow/pkg/secrets"
)

type AuthKeys interface {
	LoadFromPath(dirPath string) error
	LoadFromSecretsManager(secretsManager secrets.Secrets, keyPath, keyName string, keyType KeyType) error
	PrivateKey() []byte
	PublicKey() []byte
}

type authKeys struct {
	private []byte
	public  []byte
}

type KeyType int

const (
	Private KeyType = iota
	Public
)

func NewAuthKeys() *authKeys {
	return &authKeys{}
}

func (a *authKeys) LoadFromPath(dirPath string) error {
	privKeyPath := filepath.Join(dirPath, "id_rsa.pem")
	pubKeyPath := filepath.Join(dirPath, "id_rsa.pem.pub")

	fmt.Printf("Attempting to load private key from %s\n", privKeyPath)
	prvKey, err := os.ReadFile(privKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key (%s): %w", privKeyPath, err)
	}

	fmt.Printf("Attempting to load public key from %s\n", pubKeyPath)
	pubKey, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key (%s): %w", pubKeyPath, err)
	}

	a.private = prvKey
	a.public = pubKey

	fmt.Println("Successfullyu loaded keys.")

	return nil
}

func (a *authKeys) LoadFromSecretsManager(secretsManager secrets.Secrets, keyPath, keyName string, keyType KeyType) error {
	key, err := secretsManager.GetCryptoKey(keyPath, keyName)
	if err != nil {
		return err
	}
	if keyType == Private {
		a.private = key
	} else {
		a.public = key
	}

	return nil
}

func (a *authKeys) PrivateKey() []byte {
	if a.private == nil {
		return nil
	}
	keyCopy := make([]byte, len(a.private))
	copy(keyCopy, a.private)
	return keyCopy
}

func (a *authKeys) PublicKey() []byte {
	if a.public == nil {
		return nil
	}
	keyCopy := make([]byte, len(a.public))
	copy(keyCopy, a.public)
	return keyCopy
}
