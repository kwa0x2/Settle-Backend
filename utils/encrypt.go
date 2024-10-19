package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// PKCS7 Doldurma Fonksiyonu
func pkcs7Padding(data []byte, blockSize int) []byte {
	padSize := blockSize - len(data)%blockSize
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(data, pad...)
}

// Encrypt şifreleme fonksiyonu
func Encrypt(plainText string, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes long for AES-256")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// PlainText'i byte dizisine çevir ve doldurma işlemi yap
	paddedData := pkcs7Padding([]byte(plainText), aes.BlockSize)
	cipherText := make([]byte, len(paddedData))
	mode := NewECBEncrypter(block)
	mode.CryptBlocks(cipherText, paddedData)
	return base64.RawURLEncoding.EncodeToString(cipherText), nil
}

// ECB Şifreleme Yapılandırması
type ecb struct {
	b         cipher.Block
	blockSize int
}

func NewECBEncrypter(b cipher.Block) *ecb {
	return &ecb{b: b, blockSize: b.BlockSize()}
}

func (e *ecb) BlockSize() int { return e.blockSize }

func (e *ecb) CryptBlocks(dst, src []byte) {
	if len(src)%e.blockSize != 0 {
		panic("input not full blocks")
	}
	if len(dst) < len(src) {
		panic("output smaller than input")
	}
	for len(src) > 0 {
		e.b.Encrypt(dst[:e.blockSize], src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}
