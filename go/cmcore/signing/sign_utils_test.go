package signing

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"os"
	"sync"
	"testing"
)

func TestKeyMarshaling(t *testing.T) {
	// random key for tests
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	privkSer, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := privateKey.PublicKey
	pubkSer, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		t.Fatal(err)
	}

	base64PrivKey := base64.StdEncoding.EncodeToString(privkSer)
	base64PubKey := base64.StdEncoding.EncodeToString(pubkSer)

	// Useful for generating new keys
	//fmt.Printf("Private Key:\n%s\n\n", base64PrivKey)
	//fmt.Printf("Public Key:\n%s\n\n", base64PubKey)

	// Test private key unmarshal
	u, err := NewSignUtilWithSerializedPrivateKey(base64PrivKey)
	if err != nil {
		t.Fatal(err)
	}
	if u.privateKey == nil || u.publicKey == nil {
		t.Fatal("Failed to parse private key")
	}
	if !u.privateKey.Equal(privateKey) || !u.publicKey.Equal(&privateKey.PublicKey) {
		t.Fatal("Incorrect key parsed")
	}

	// Test public key unmarshal
	pu, err := NewSignUtilWithSerializedPublicKey(base64PubKey)
	if err != nil {
		t.Fatal(err)
	}
	if pu.privateKey != nil || pu.publicKey == nil {
		t.Fatal("Failed to parse public key")
	}
	if !pu.publicKey.Equal(&privateKey.PublicKey) {
		t.Fatal("Incorrect key parsed")
	}
}

func TestBuiltinKey(t *testing.T) {
	msg := []byte("Hello world")
	// Signature from our real private key
	signature := "MEYCIQCbS7vF3X2t31nhh/BRNt/VrX2BIK2t1OGEINd4hMtdoAIhAJy6xAyj1OWNt+Y492Pq3tKntVYImJTMYPqbj9FNkpxH"

	verified, err := SharedSignUtil().VerifyMessage(msg, signature)
	if err != nil {
		t.Fatal(err)
	}
	if !verified {
		t.Fatal("Message not verified")
	}
}

var testPrivKey = "MGgCAQEEHOEUmigOOoZ+STQ1jkYuXRN2hXLbxLKTvKdlXEygBwYFK4EEACGhPAM6AASDljuXqf/dic4vnAfRtqFsl/fQANciY+xACkgOOE9MGgvu+XIfTOqsqagLJ6ZUedbZus5FUa4awQ=="
var testPubKey = "ME4wEAYHKoZIzj0CAQYFK4EEACEDOgAEg5Y7l6n/3YnOL5wH0bahbJf30ADXImPsQApIDjhPTBoL7vlyH0zqrKmoCyemVHnW2brORVGuGsE="

func TestSignAndVerifyMessage(t *testing.T) {

	psu, err := NewSignUtilWithSerializedPrivateKey(testPrivKey)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("Hello world")
	sig, err := psu.SignMessage(msg)
	if err != nil {
		t.Fatal(err)
	}

	su, err := NewSignUtilWithSerializedPublicKey(testPubKey)
	if err != nil {
		t.Fatal(err)
	}
	valid, err := su.VerifyMessage(msg, sig)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("Validaton failed")
	}
}

func TestLoadPrivKeyFromEnv(t *testing.T) {
	origEnv := os.Getenv(privateKeyEnvVarName)
	// defer cleanup critical for cleanup -- once we clear the env var, let singleton re-run
	defer func() {
		t.Setenv(privateKeyEnvVarName, origEnv)
		sharedSignUtil = nil
		privateKeyOnce = sync.Once{}
	}()

	t.Setenv(privateKeyEnvVarName, "")
	// allow our singleton to repopulate
	sharedSignUtil = nil
	privateKeyOnce = sync.Once{}
	if SharedSignUtil().privateKey != nil {
		t.Fatal("PrivKey populated when it shouldn't be from env")
	}
	t.Setenv(privateKeyEnvVarName, testPrivKey)
	// allow our singleton to repopulate
	sharedSignUtil = nil
	privateKeyOnce = sync.Once{}
	if SharedSignUtil().privateKey == nil {
		t.Fatal("Failed to load PrivKey from env")
	}

}
