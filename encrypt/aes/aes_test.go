//@File     aes_test.go
//@Time     2022/08/23
//@Author   #Suyghur,

package aes

import (
	"encoding/base64"
	"testing"
)

func TestGenerateIV(t *testing.T) {
	key := []byte("abcdef0123456789")
	iv := GenerateIV(key)
	t.Logf("秘钥 : %s", key)
	t.Logf("iv : %s", iv)
}

func TestPKCS7CBCEncrypt(t *testing.T) {
	text := "hello world"
	aesKey := []byte("abcdef0123456789")
	iv := GenerateIV(aesKey)
	t.Logf("明文 : %s", text)
	t.Logf("秘钥 : %s", aesKey)
	t.Logf("iv : %s", iv)
	enc, err := PKCS7CBCEncrypt([]byte(text), aesKey, iv)
	if err != nil {
		t.Error(err)
	}
	t.Logf("密文 : %s", base64.StdEncoding.EncodeToString(enc))
}

func TestPKCS7CBCDecrypt(t *testing.T) {
	text := "RMpvM76TJoUE9qXphAPJiw=="
	aesKey := []byte("abcdef0123456789")
	iv := GenerateIV(aesKey)
	enc, _ := base64.StdEncoding.DecodeString(text)
	t.Logf("明文 : %s", text)
	t.Logf("秘钥 : %s", aesKey)
	t.Logf("iv : %s", iv)
	raw, err := PKCS7CBCDecrypt(enc, aesKey, iv)
	if err != nil {
		t.Error(err)
	}
	t.Logf("明文 : %s", raw)
}
