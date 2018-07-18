package sinapay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"log"
	"sync"
)

var (
	// needRSAKey 需要加密的key
	needRSAKey = map[string]bool{"real_name": true, "cert_no": true, "verify_entity": true, "bank_account_no": true, "account_name": true,
		"phone_no": true, "validity_period": true, "verification_value": true, "telephone": true, "email": true, "organization_no": true,
		"legal_person": true, "legal_person_phone": true, "agent_name": true, "license_no": true, "agent_mobile": true, "trade_related_no": true}
	privateKey, publicKey           []byte
	privateKeyBlock, publicKeyBlock *pem.Block
	once                            sync.Once
)

//配置
const (
	charset      = "UTF-8"
	base64format = "UrlSafeNoPadding"
	keyType      = "PKCS8"
	algorithm    = crypto.SHA1
)

var (
	rsaKeys *RsaKeys
)

// RsaKeys 新浪支付证书
type RsaKeys struct {
	publicKey        *rsa.PublicKey
	privateKey       *rsa.PrivateKey
	sinapaypublicKey *rsa.PublicKey
}

// GetRSA 获取对象
func GetRSA() *RsaKeys {
	once.Do(newRsaKeys)
	return rsaKeys
}

// newRsaKeys 初始化对象
func newRsaKeys() {
	publicKey, err := ioutil.ReadFile(publicPEM)
	if err != nil {
		log.Panicln(publicPEM, " is cannot open")
	}
	privateKey, err := ioutil.ReadFile(privatePEM)
	if err != nil {
		log.Panicln(privatePEM, " is cannot open")
	}
	sinapaypublicKey, err := ioutil.ReadFile(sinaPayPublicPEM)
	if err != nil {
		log.Panicln(sinaPayPublicPEM, " is cannot open")
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		log.Panicln("public key error")
		return
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Panicln(err)
		return
	}
	pub := pubInterface.(*rsa.PublicKey)
	block, _ = pem.Decode(privateKey)
	if block == nil {
		log.Panicln("private key error")
		return
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Panicln(err)
		return
	}
	pri, ok := priv.(*rsa.PrivateKey)
	block, _ = pem.Decode(sinapaypublicKey)
	if block == nil {
		log.Panicln("sina pay public key error")
		return
	}
	sinapaypubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Panicln(err)
		return
	}
	sinapaypub := sinapaypubInterface.(*rsa.PublicKey)
	if ok {
		rsaKeys = &RsaKeys{
			publicKey:        pub,
			privateKey:       pri,
			sinapaypublicKey: sinapaypub,
		}
		return
	}
	log.Panicln("private key not supported")
	return
}

// PublicEncrypt 使用新浪公钥加密
func (r *RsaKeys) PublicEncrypt(data string) (string, error) {
	partLen := r.publicKey.N.BitLen()/8 - 11
	chunks := split([]byte(data), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, r.sinapaypublicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(bytes)
	}

	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

// PrivateDecrypt 使用客户私钥解密
func (r *RsaKeys) PrivateDecrypt(encrypted string) (string, error) {
	partLen := r.publicKey.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(encrypted)
	chunks := split([]byte(raw), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}

	return buffer.String(), err
}

// Sign 数据加签
func (r *RsaKeys) Sign(data string) (string, error) {
	h := algorithm.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	//
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, algorithm, hashed)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sign), err
}

// Verify 数据验签 使用客户公钥在测试环境通过
func (r *RsaKeys) Verify(data string, sign string) error {
	h := algorithm.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	decodedSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		log.Println(err)
		return err
	}
	return rsa.VerifyPKCS1v15(r.publicKey, algorithm, hashed, decodedSign)
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
