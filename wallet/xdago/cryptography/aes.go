package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AesEncrypt(data, key, iv []byte) ([]byte, error) {
	aesBlockEncryptor, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	content := PKCS5Padding(data, aesBlockEncryptor.BlockSize())
	encrypted := make([]byte, len(content))

	aesEncryptor := cipher.NewCBCEncrypter(aesBlockEncryptor, iv)
	aesEncryptor.CryptBlocks(encrypted, content)

	return encrypted, nil
}

func AesDecrypt(encrypted, key, iv []byte) ([]byte, error) {

	decrypted := make([]byte, len(encrypted))
	var aesBlockDecrypt cipher.Block
	aesBlockDecrypt, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesDecrypt := cipher.NewCBCDecrypter(aesBlockDecrypt, iv)
	aesDecrypt.CryptBlocks(decrypted, encrypted)
	content, err := PKCS5Trimming(decrypted)
	if err != nil {
		return nil, err
	}
	return content, nil

}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS5Trimming(encrypt []byte) ([]byte, error) {
	padding := encrypt[len(encrypt)-1]
	if len(encrypt)-int(padding) > 0 {
		return encrypt[:len(encrypt)-int(padding)], nil
	}
	return nil, errors.New("aes decrypt pkcs5 trimming error")
}
