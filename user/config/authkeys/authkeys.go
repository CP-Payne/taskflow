package authkeys

type AuthKeys interface {
	Load(dirPath string) error
	PrivateKey() []byte
	PublicKey() []byte
}
