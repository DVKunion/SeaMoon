package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"time"

	"github.com/DVKunion/SeaMoon/pkg/system/errors"
)

// PKCS7Padding pads the plaintext to be a multiple of the block size
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7UnPadding unpads the plaintext
func PKCS7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}
	unpadding := int(data[length-1])
	return data[:(length - unpadding)], nil
}

// AESEncrypt encrypts the given plaintext using AES with the given key and returns the ciphertext
func AESEncrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)

	ciphertext := make([]byte, blockSize+len(plaintext))
	iv := ciphertext[:blockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt decrypts the given ciphertext using AES with the given key and returns the plaintext
func AESDecrypt(ciphertext string, key []byte) ([]byte, error) {
	cipherData, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(cipherData) < blockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := cipherData[:blockSize]
	cipherData = cipherData[blockSize:]

	if len(cipherData)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherData, cipherData)

	return PKCS7UnPadding(cipherData)
}

// GenerateKey generates a key using HMAC-SHA256 based on the provided time
func GenerateKey(t time.Time) []byte {
	timeBytes := []byte(t.Format(time.DateTime))
	h := hmac.New(sha256.New, []byte(""))
	h.Write(timeBytes)
	return h.Sum(nil)
}
