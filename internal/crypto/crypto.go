// Package crypto описание работы шифрованием
package crypto

// Encryptor шифровальшик
type Encryptor interface {
	Encrypt([]byte) ([]byte, error)
}

type EncryptFunc func([]byte) ([]byte, error)

func (f EncryptFunc) Encrypt(b []byte) ([]byte, error) {
	return f(b)
}

var _ Encryptor = EncryptFunc(nil)

// Decryptor расшифровщик
type Decryptor interface {
	Decrypt([]byte) ([]byte, error)
}

type DecryptFunc func([]byte) ([]byte, error)

func (f DecryptFunc) Decrypt(b []byte) ([]byte, error) {
	return f(b)
}

type EncryptDecryptor interface {
	Encryptor
	Decryptor
}
