package cryptogram

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt(data, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nil, nonce, data, nil)

	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func Decrypt(encrypted string, key []byte) (string, error) {
	encData, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce, ciphertext := encData[:12], encData[12:]

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func Example() {
	key := []byte("myverystrongpasswordo32bitlength")

	data := "Это секретное сообщение!"

	// Encryption
	encrypted, err := Encrypt([]byte(data), key)
	if err != nil {
		return
	}

	// Decryption
	_, err = Decrypt(encrypted, key)
	if err != nil {
		return
	}
}
