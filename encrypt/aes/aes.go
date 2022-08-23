//@File     aes.go
//@Time     2022/08/23
//@Author   #Suyghur,

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// PKCS7CBCEncrypt AES加密
func PKCS7CBCEncrypt(origData, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	enc := make([]byte, len(origData))
	blockMode.CryptBlocks(enc, origData)
	return enc, nil
}

// PKCS7CBCDecrypt AES解密
func PKCS7CBCDecrypt(encData, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(encData))
	blockMode.CryptBlocks(origData, encData)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func GenerateIV(key []byte) []byte {
	length := len(key)
	iv := make([]byte, length)
	for i := 0; i < length; i++ {
		iv[i] = key[length-1-i]
	}
	return iv
}
