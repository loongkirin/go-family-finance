package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func BcryptHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func BcryptVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Encrypt(key, text string) (string, error) {
	return EncryptBytes([]byte(key), []byte(text))
}

func EncryptBytes(key, text []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCTR(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return EncodeBytesToBase64(ciphertext), nil
}

func Decrypt(key, b64s string) (string, error) {
	text, err := DecodeBase64ToBytes(b64s)
	if err != nil {
		return "", nil
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", nil
	}
	if len(text) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCTR(block, iv)
	cfb.XORKeyStream(text, text)
	return string(text), nil
}

func EncodeBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func EncodeBytesToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeBase64(b64s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(b64s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DecodeBase64ToBytes(b64s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64s)
}
