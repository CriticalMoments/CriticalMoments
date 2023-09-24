package signing

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"sync"
)

type SignUtil struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewSignUtilWithSerializedPrivateKey(privateKeyString string) (*SignUtil, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyString)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS1PrivateKey([]byte(privateKeyBytes))
	if err != nil {
		return nil, err
	}
	u := SignUtil{
		publicKey:  &key.PublicKey,
		privateKey: key,
	}
	return &u, nil
}

func NewSignUtilWithSerializedPublicKey(publicKeyString string) (*SignUtil, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyString)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	u := SignUtil{
		publicKey: key,
	}
	return &u, nil
}

const cmPublicKey = `MIIBCgKCAQEA5Is8dkh83sAp4kfSxV9DEzxNF4VYwUUpRQ6E+uKsE44UbB6oDGRW7+xxvStGNrHT7+KiKVAEp8793iLqthOUfOcF8GEN+9sKSovZM7Fvdv3ZR7YB8U/sLv1gUo9Mi8xYHov7VuseWA1+XChlfCNv58hNtip/9Qz8Y6ViifiEA5KCSbo4wjUa7ULbWYdG3/PQouvDVb2OKY5+T0oxDRGzHkkq9GRqjxC5FuqLo/wWgUJnGrylCqvAmC5i0s7Cr4uH6bNINl8PuGIwWwl352sOZVCpEDJ2+j4ilp/iwgw+EHj/4nr+u5lLtPLQK1vbVnTGZCz+1+2CQAkbRUAJglS6ywIDAQAB`

var sharedSignUtil *SignUtil
var privateKeyOnce sync.Once

func SharedSignUtil() *SignUtil {
	privateKeyOnce.Do(func() {
		envPrivKey := os.Getenv("PRIVATE_CM_SIGN_KEY")
		if envPrivKey != "" {
			privateSignUtil, err := NewSignUtilWithSerializedPrivateKey(envPrivKey)
			if err != nil {
				fmt.Println("WARNING: a PRIVATE_CM_SIGN_KEY env var was was, but wasn't parseable. Signing and validation will fail.")
			} else {
				sharedSignUtil = privateSignUtil
			}
		}
		if sharedSignUtil == nil {
			publicSignKey, err := NewSignUtilWithSerializedPublicKey(cmPublicKey)
			if err != nil {
				panic("Hardcoded key is not parsable. We have serious issue.")
			}
			sharedSignUtil = publicSignKey
		}
	})
	return sharedSignUtil
}

func msgHash(msg []byte) ([]byte, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		return nil, err
	}
	return msgHash.Sum(nil), nil
}

func (u *SignUtil) SignMessage(msg []byte) (string, error) {
	if u.privateKey == nil {
		return "", errors.New("Can't sign a message without a private key")
	}

	msgHashRaw, err := msgHash(msg)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPSS(rand.Reader, u.privateKey, crypto.SHA256, msgHashRaw, nil)
	if err != nil {
		return "", err
	}

	stringSignature := base64.StdEncoding.EncodeToString(signature)
	return stringSignature, nil
}

func (u *SignUtil) VerifyMessage(msg []byte, signatureString string) error {
	signature, err := base64.StdEncoding.DecodeString(signatureString)
	if err != nil {
		return err
	}

	msgHashRaw, err := msgHash(msg)
	if err != nil {
		return err
	}

	return rsa.VerifyPSS(u.publicKey, crypto.SHA256, msgHashRaw, signature, nil)
}
