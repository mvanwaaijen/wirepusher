package wirepusher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// bin2hex convenience function to make it similar to c# wirepusher example
func bin2hex(data []byte) string {
	return hex.EncodeToString(data)
}

// hex2bin convenience function to make it similar to c# wirepusher example
func hex2bin(data string) ([]byte, error) {
	return hex.DecodeString(data)
}

// randomIV returns a randomly generated iv of 16 bytes which is needed for AES encryption
func randomIV() ([]byte, error) {
	riv := make([]byte, 16)
	_, err := rand.Read(riv)
	if err != nil {
		return nil, errors.Wrap(err, "generation of random iv failed")
	}
	return riv, nil
}

// encrypt the text with aes using the iv and key
func encrypt(text string, iv []byte, key []byte) (string, error) {
	if len(key) != AES_KEY_SIZE {
		return "", fmt.Errorf("invalid key length. expected: %d, found: %d", AES_KEY_SIZE, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "creating new aes-128 cipher failed")
	}

	textToEncrypt, err := pkcs7pad([]byte(text), len(key))
	if err != nil {
		return "", errors.Wrap(err, "padding failed")
	}

	encryptedText := make([]byte, len(textToEncrypt))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptedText, textToEncrypt)

	return clean(base64.StdEncoding.EncodeToString(encryptedText)), nil
}

// pkcs7pad add pkcs7 padding
func pkcs7pad(data []byte, blockSize int) ([]byte, error) {
	if blockSize < 0 || blockSize > 256 {
		return nil, fmt.Errorf("pkcs7: Invalid block size %d", blockSize)
	} else {
		padLen := blockSize - len(data)%blockSize
		padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
		return append(data, padding...), nil
	}
}

// clean the encrypted and base64 encoded from the codes '+', '/' and '=' and replace them with '-', '.' and '_'
func clean(text string) string {
	replacer := strings.NewReplacer("+", "-", "/", ".", "=", "_")
	return replacer.Replace(text)
}

// hash the password with sha1 algorithm
func hash(password string) ([]byte, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return nil, errors.Wrap(err, "error hashing password")
	}
	return hash.Sum(nil), nil
}
