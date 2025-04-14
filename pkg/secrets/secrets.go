package secrets

type Secrets interface {
	GetCryptoKey(path string, keyName string) ([]byte, error)
}
