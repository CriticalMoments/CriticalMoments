package signing

import (
	"testing"
)

const testBundleId = "io.criticalmoments.demo"

func TestApiKeyCreation(t *testing.T) {
	psu, err := NewSignUtilWithSerializedPrivateKey(testPrivKey)
	if err != nil {
		t.Fatal(err)
	}

	k, err := NewSignedApiKeyWithSigner("", psu)
	if err == nil || k != nil {
		t.Fatal("Allowed empty bundle ID")
	}

	k, err = NewSignedApiKeyWithSigner(testBundleId, psu)
	if err != nil {
		t.Fatal(err)
	}
	if issue := testKeyIssue(k, psu); issue != "" {
		t.Fatal(issue)
	}

	p, err := ParseApiKey(k.String())
	if err != nil {
		t.Fatal(err)
	}
	if issue := testKeyIssue(p, psu); issue != "" {
		t.Fatal(issue)
	}
}

func testKeyIssue(k *ApiKey, psu *SignUtil) string {
	v, err := k.ValidWithSigner(psu)
	if err != nil {
		return err.Error()
	}
	if !v {
		return "key fails validation"
	}
	if k.BundleId() != testBundleId {
		return "Failed to save bid"
	}
	if k.version != 1 {
		return "Failed to save version"
	}
	if k.signature == "" {
		return "Failed to sign"
	}
	if k.String() == "" {
		return "Failed to stringer"
	}
	return ""
}

func TestApiKeyParsing(t *testing.T) {
	_, err := ParseApiKey("invalid")
	if err == nil {
		t.Fatal("invalid key parsed")
	}
	_, err = ParseApiKey("CM-invalid-sdf")
	if err == nil {
		t.Fatal("invalid key parsed")
	}
	_, err = ParseApiKey("CM12aa-invalid-sdf")
	if err == nil {
		t.Fatal("invalid key with non-digit version parsed")
	}

	// Valid Key Format but invaid signature
	// Never remove this. Future libraries need to continue to support older formats
	k, err := ParseApiKey("CM99-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-SIGNATURE")
	if err != nil {
		t.Fatal(err)
	}
	if k.Version() != 99 {
		t.Fatal("Failed to parse version number")
	}
	if k.signature != "SIGNATURE" {
		t.Fatal("Failed to parse signature")
	}
	if len(k.props) != 1 || k.BundleId() != testBundleId {
		t.Fatal("failed to parse bundle ID")
	}
	v, err := k.Valid()
	if err == nil || v {
		t.Fatal("Signature validation passed on incorrect signature")
	}

	psu, err := NewSignUtilWithSerializedPrivateKey(testPrivKey)
	if err != nil {
		t.Fatal(err)
	}

	// Valid key (with test signature)
	// Never remove this. Future libraries need to continue to support older formats
	k, err = ParseApiKey("CM1-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-MD0CHQDzGgtQb5nmadLEmE4OFSg3LHiCBkRFTgjQhz6ZAhx4OTuVnYWmdrFx5D3ysT2d+6QfzIsGzBsXrI+T")
	if err != nil {
		t.Fatal(err)
	}
	if issue := testKeyIssue(k, psu); issue != "" {
		t.Fatal(issue)
	}

	// normal
	// unsigned
	// missing bundle ID (but valid props)
	// unknown version number (future)
	// valiod format but missing b: key
	// No verion number / no prefix / no dashes / one dash etc
	// future ones: higher version number / more properties unrecognized
}

// TODO test real key with built in public Key
