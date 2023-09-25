package signing

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"sync"
)

const privateKeyEnvVarName = "PRIVATE_CM_EC_KEY"

type SignUtil struct {
	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

func NewSignUtilWithSerializedPrivateKey(privateKeyString string) (*SignUtil, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyString)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParseECPrivateKey(privateKeyBytes)
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
	keyI, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	key, ok := keyI.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("incorrect key type")
	}
	u := SignUtil{
		publicKey: key,
	}
	return &u, nil
}

const cmPublicKey = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAESz8KF+TKa1t02O+nx+tKqfT5Nx5GIb6UDjpCtFiQ6Pz5nbmAl5fDfgDjAcTl9Fh2CWSL9KjNanUEMxlYoLELWg=="

var sharedSignUtil *SignUtil
var privateKeyOnce sync.Once

func SharedSignUtil() *SignUtil {
	privateKeyOnce.Do(func() {
		envPrivKey := os.Getenv(privateKeyEnvVarName)
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

	signature, err := ecdsa.SignASN1(rand.Reader, u.privateKey, msgHashRaw)
	if err != nil {
		return "", err
	}

	stringSignature := base64.StdEncoding.EncodeToString(signature)
	return stringSignature, nil
}

func (u *SignUtil) VerifyMessage(msg []byte, signatureString string) (bool, error) {
	signature, err := base64.StdEncoding.DecodeString(signatureString)
	if err != nil {
		return false, err
	}

	msgHashRaw, err := msgHash(msg)
	if err != nil {
		return false, err
	}

	return ecdsa.VerifyASN1(u.publicKey, msgHashRaw, signature), nil
}
