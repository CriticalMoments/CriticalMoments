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

	// Valid/typical V1 key (with test signature)
	// Never remove this. Future libraries need to continue to support older formats
	k, err = ParseApiKey("CM1-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-MD0CHQDzGgtQb5nmadLEmE4OFSg3LHiCBkRFTgjQhz6ZAhx4OTuVnYWmdrFx5D3ysT2d+6QfzIsGzBsXrI+T")
	if err != nil {
		t.Fatal(err)
	}
	if issue := testKeyIssue(k, psu); issue != "" {
		t.Fatal(issue)
	}

	// unsigned key
	k, err = ParseApiKey("CM1-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==")
	if err == nil {
		t.Fatal(err)
	}

	// future API version should be okay
	k, err = ParseApiKey("CM999-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-MDwCHCy2eURhq3oIpK5VY99l5/dZF30vwiLm31kar68CHCQD1ftBhPhxcNl6sys+obIooH5K/r/2dOzMoqc=")
	if err != nil || k.version != 999 || k.BundleId() != testBundleId {
		t.Fatal("API Key with future version not parsed")
	}

	// missing bundle ID (but another valid prop)
	k, err = ParseApiKey("CM1-eDppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-MD0CHQDSQjQuYPtfG9xRn1KvVQOI3zMRJu/YCZ1XoaLVAhwEzgxjn8ysier97gZjW0+JR9g9yGbiSVPNxcUY")
	if err == nil {
		t.Fatal("API Key missing bundle ID passed")
	}

	// no props
	k, err = ParseApiKey("CM1-Yjo=-MD0CHQC/sbrxCs5VI/NL86juc1SJpyJkZrhuCAOXn/ObAhwocvTYZF84qu/0rBu0+y1fJFkOfsEpM1YiDmQj")
	if err == nil {
		t.Fatal("API Key missing all props")
	}

	// future proof: includes properties this library version doesn't recognize
	k, err = ParseApiKey("CM1-aGVsbG86d29ybGQ=-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-MDwCHCcUOxCOKtnH8OfXYEJSS/Wt6ieDUe1FzJK+EDkCHG6g5F1rV+5n+dqfnoPYvbkLCwtRYRtDCw+cUJc=")
	if err != nil || k.version != 1 || k.BundleId() != testBundleId || k.props["hello"] != "world" {
		t.Fatal("API Key missing all props")
	}

	// empty key
	k, err = ParseApiKey("")
	if err == nil {
		t.Fatal("Empty key passes")
	}
}

func TestKeySignedWithPublicKey(t *testing.T) {
	// Real API key signed with built in Public Key
	k, err := ParseApiKey("CM1-aGVsbG86d29ybGQ=-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-MEUCIQCUfx6xlmQ0kdYkuw3SMFFI6WXrCWKWwetXBrXXG2hjAwIgWBPIMrdM1ET0HbpnXlnpj/f+VXtjRTqNNz9L/AOt4GY=")
	if err != nil {
		t.Fatal(err)
	}
	if k.BundleId() != testBundleId || k.version != 1 {
		t.Fatal("Failed to parse public key api key")
	}
	if v, err := k.Valid(); err != nil || !v {
		t.Fatal("API Key with public key failed signature")
	}
}
