//@File     rsa.go
//@Time     2022/08/23
//@Author   #Suyghur,

package rsa

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"io/ioutil"
)

type Security struct {
	//公钥字符串
	pubStr string
	//私钥字符串
	priStr string
	//公钥
	pubKey *rsa.PublicKey
	//私钥
	priKey *rsa.PrivateKey
}

var RSA = &Security{}

// SetPublicKey 设置公钥
func (sec *Security) SetPublicKey(pubStr string) (err error) {
	sec.pubStr = pubStr
	sec.pubKey, err = sec.GetPublicKey()
	return err
}

// SetPrivateKey 设置私钥
func (sec *Security) SetPrivateKey(priStr string) (err error) {
	sec.priStr = priStr
	sec.priKey, err = sec.GetPrivateKey()
	return err
}

func (sec *Security) GetPublicKey() (*rsa.PublicKey, error) {
	return getPubKey([]byte(sec.pubStr))
}

func (sec *Security) GetPrivateKey() (*rsa.PrivateKey, error) {
	return getPriKey([]byte(sec.priStr))
}

// PubKeyEncrypt 公钥加密
func (sec *Security) PubKeyEncrypt(input []byte) ([]byte, error) {
	if sec.pubKey == nil {
		return []byte(""), errors.New(`Please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(sec.pubKey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// PubKeyDecrypt 公钥解密
func (sec *Security) PubKeyDecrypt(input []byte) ([]byte, error) {
	if sec.pubKey == nil {
		return []byte(""), errors.New(`Please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(sec.pubKey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// PriKeyEncrypt 私钥加密
func (sec *Security) PriKeyEncrypt(input []byte) ([]byte, error) {
	if sec.priKey == nil {
		return []byte(""), errors.New(`Plesase set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(sec.priKey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// PriKeyDecrypt 私钥解密
func (sec *Security) PriKeyDecrypt(input []byte) ([]byte, error) {
	if sec.priKey == nil {
		return []byte(""), errors.New(`Please set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(sec.priKey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}
