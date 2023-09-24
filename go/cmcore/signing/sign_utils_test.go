package signing

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestKeyMarshaling(t *testing.T) {
	// random key for tests
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	privkSer := x509.MarshalPKCS1PrivateKey(privateKey)
	pubkSer := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

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
	privKeySerial := `MIIEowIBAAKCAQEA5Is8dkh83sAp4kfSxV9DEzxNF4VYwUUpRQ6E+uKsE44UbB6oDGRW7+xxvStGNrHT7+KiKVAEp8793iLqthOUfOcF8GEN+9sKSovZM7Fvdv3ZR7YB8U/sLv1gUo9Mi8xYHov7VuseWA1+XChlfCNv58hNtip/9Qz8Y6ViifiEA5KCSbo4wjUa7ULbWYdG3/PQouvDVb2OKY5+T0oxDRGzHkkq9GRqjxC5FuqLo/wWgUJnGrylCqvAmC5i0s7Cr4uH6bNINl8PuGIwWwl352sOZVCpEDJ2+j4ilp/iwgw+EHj/4nr+u5lLtPLQK1vbVnTGZCz+1+2CQAkbRUAJglS6ywIDAQABAoIBAEhUPIVeuY4xmM/RVUY7uNmsmuVXwVghUEdXqgRQmo7xx0rUhPCvDMiPtwtcV7NVojJoMlQKy/5jxvp3aHrJRZQl9T43KRrNHruq+MmgXRt2iT5lvsWlOqVAcSyPx3Ty7ex09s1ySb8qPhRigIPCH1dmkBmX57khK/tJSx9JNFaAfXprBd38E+KGCa3B7RB1u6HQswN2TzWo6r1caq4kkaKGR9XDNz1SI4rkrtRqV11F+nk8OJaSbSF/rQdDyBhvhfrRgIFQ8fgOAyMaPemzOew6adzWP26xJuRKW59hyvekAgUyl6jyQQGYp5bQRMTQHg34YIcmkj4P+0vR/jK/OAECgYEA8yG02YDsQl4thhtX25QjBnfWZvBdhdJ1ORARU/YMlD8LEJKhYPCV8KLfdGNXvqDmimS1UFUCpXttZD0NtdOBZTgPjdI9/Qm0Jw6OhAazdZXP1oRwxzPgdK/BeRMRU5NdtlpaOxeP/ZMmtqBOM2N0y+2QhPAq1GLz2AsXkT/3/AECgYEA8KPfstof+xlF76DCSYJo7LvCsQ13Zr/8RcY2W4i0n82B1drLNVrfvFUVxUu2FRttQZKu0tR0l2RU4/y4ImUV1/eK7Xesqll3TnDvmnSki2BHaliCIXbY1wWyomOcR7haFCjM
N+1UEMc+N9a3QyU2IKA5unnPeUN2h/EC/BhH5ssCgYBbGUoWJURhKcCM+znUQJFPHx/qui2QsubRVr/nYc4czfJrZ0WoePz1iVGI3qBGASvgtxNo4jF3p+O5J1c3xeQ59ON/FEO9yCEEcWPc/FXJvTR/AGjxevKjRieMIiTf19vJM9mTQqTlMnnS/AXRI3bj4kPAS+0AX4NWc/GErx9QAQKBgEYsE2SNVPwdH5bEM0PKYpx+GEUXHzV4ULFsHpfMopdjDzR0jANwD4RU73dMH7nB+LdBdfeG+sTW/iZJoMxu29LRndKnrlMyqabXKhfJYd4+4jRxwOjPRmZVhAT0tTL44FO2ne7FJ1mJMGyKEYDkDgevkYX+VXEQKjV0I6Gt1vHHAoGBAMSQUDRIr6UQ44qnPOjRZYRizxTln/lGXHj8BFTgVaT2f/w8nlCQE/TDnGhiqdzi+msjmonOjvn00ItoBj/ll6WTSnD0kX+eNkh3XqduOBabFfgQIMgpitq5Lkh83pkhN+2l7nTKTCEn8jHO9qMkh46sqS4BRQqEwKoyZhFcKiCH`

	psu, err := NewSignUtilWithSerializedPrivateKey(privKeySerial)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("Hello world")
	signature, err := psu.SignMessage(msg)
	if err != nil {
		t.Fatal(signature)
	}
	fmt.Printf("Signatrure: %v\n", signature)

	err = SharedSignUtil().VerifyMessage(msg, signature)
	if err != nil {
		t.Fatal(err)
	}

	// TODO actually test. Need a message and stored signature
}

func TestSignAndVerifyMessage(t *testing.T) {
}
