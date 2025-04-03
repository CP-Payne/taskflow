package authkeys

import (
	"fmt"
	"os"
	"path/filepath"
)

type authKeys struct {
	private []byte
	public  []byte
}

func NewAuthKeys() *authKeys {
	return &authKeys{}
}

func (a *authKeys) Load(dirPath string) error {
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
