// Package aes реализация шифрования на aes
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type AesCrypt struct {
	gcm cipher.AEAD
}

func New(key []byte) (*AesCrypt, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AesCrypt{aesgcm}, nil
}

// Encrypt - шифрование aes
func (a *AesCrypt) Encrypt(data []byte) ([]byte, error) {

	nonce := make([]byte, a.gcm.NonceSize()) //12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	res := a.gcm.Seal(nil, nonce, data, nil)
	res = append(nonce, res...)

	return res, nil
}

// Decrypt - расшифровка
func (a *AesCrypt) Decrypt(ciphertext []byte) ([]byte, error) {

	nonceSize := a.gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	res, err := a.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}
