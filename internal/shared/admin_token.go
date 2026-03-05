package shared

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

func tokenCipher() (cipher.AEAD, error) {
	key := sha256.Sum256([]byte(Config.Key))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func EncryptAdminToken(plaintext string) (string, error) {
	aead, err := tokenCipher()
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nil, nonce, []byte(plaintext), nil)
	combined := append(nonce, ciphertext...)
	return base64.RawURLEncoding.EncodeToString(combined), nil
}

func DecryptAdminToken(token string) (string, error) {
	aead, err := tokenCipher()
	if err != nil {
		return "", err
	}

	combined, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	if len(combined) < nonceSize {
		return "", errors.New("token too short")
	}

	nonce := combined[:nonceSize]
	ciphertext := combined[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
