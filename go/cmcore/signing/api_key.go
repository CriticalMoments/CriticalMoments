package signing

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ApiKey struct {
	raw       string // Only populated if parsed
	version   int
	props     map[string]string
	signature string
}

func (k *ApiKey) Version() int {
	return k.version
}

func (k *ApiKey) BundleId() string {
	return k.props[propKeyBundleId]
}

func (k *ApiKey) Valid() (bool, error) {
	return k.ValidWithSigner(SharedSignUtil())
}

func (k *ApiKey) ValidWithSigner(u *SignUtil) (bool, error) {
	return u.VerifyMessage([]byte(k.signedPortion()), k.signature)
}

/*
 API Key Format
  - Top level separator "-"
	- 3 sections: API Key Version, props, signature (- seperated)
	- API version: CM1 for now. Can evolve later, but should be back compatible so V1 can parse any future version.
	- Props: a list. only "b" (ios bundle ID) is required for now. Ignore others when parsing (but do sign). Divided by dashes again, base64 encoded.
	- Signature: a signature of the entire key, up until here
*/

const currentApiKeyVersion = 1
const apiPrefix = "CM"
const apiSeparator = "-"
const propertyKeyDelimiter = ":"
const propKeyBundleId = "b"

func NewSignedApiKey(bundleId string) (*ApiKey, error) {
	return NewSignedApiKeyWithSigner(bundleId, SharedSignUtil())
}

func NewSignedApiKeyWithSigner(bundleId string, u *SignUtil) (*ApiKey, error) {
	if bundleId == "" {
		return nil, errors.New("Bundle ID required")
	}

	if u.privateKey == nil {
		return nil, errors.New("Can't create new API keys without a private key")
	}

	key := ApiKey{
		version: currentApiKeyVersion,
		props: map[string]string{
			propKeyBundleId: bundleId,
		},
	}

	p := key.signedPortion()
	sig, err := u.SignMessage([]byte(p))
	if err != nil {
		return nil, err
	}
	key.signature = sig

	return &key, nil
}

func ParseApiKey(s string) (*ApiKey, error) {
	v, err := parseVersionNumber(s)
	if err != nil {
		return nil, err
	}
	sig, err := parseSignature(s)
	if err != nil {
		return nil, err
	}
	props, err := parseProperties(s)
	if err != nil {
		return nil, err
	}
	if props[propKeyBundleId] == "" {
		return nil, errors.New("API key must have bundle ID property")
	}

	k := ApiKey{
		version:   v,
		signature: sig,
		props:     props,
	}

	return &k, nil
}

func parseVersionNumber(s string) (int, error) {
	if strings.Index(s, apiPrefix) != 0 {
		return -1, errors.New("All API keys should start with CM")
	}
	i := strings.Index(s, apiSeparator)
	if i < len(apiPrefix)+1 {
		return -1, errors.New("No version number or invalid format")
	}
	versionString := s[len(apiPrefix):i]
	v, err := strconv.Atoi(versionString)
	if err != nil {
		return -1, err
	}
	return v, nil
}

func parseSignature(s string) (string, error) {
	i := strings.LastIndex(s, apiSeparator)
	if i <= 0 {
		return "", errors.New("invalid key -- no signature marker")
	}
	signature := s[i+1:]
	return signature, nil
}

func parseProperties(s string) (map[string]string, error) {
	propStart := strings.Index(s, apiSeparator)
	propEnd := strings.LastIndex(s, apiSeparator)
	if propStart <= 0 || propEnd <= propStart {
		return nil, errors.New("Invalid API key format")
	}

	values := map[string]string{}
	propsSection := s[propStart+1 : propEnd]
	for _, rawProp := range strings.Split(propsSection, apiSeparator) {
		propBytes, err := base64.StdEncoding.DecodeString(rawProp)
		if err != nil {
			return nil, errors.New("Invalid API key format")
		}
		decodedProp := string(propBytes)
		pi := strings.Index(decodedProp, propertyKeyDelimiter)
		if pi <= 0 {
			return nil, errors.New("Invalid API key format")
		}
		key := decodedProp[:pi]
		value := decodedProp[pi+1:]
		if key == "" || value == "" {
			return nil, errors.New("Invalid API Key Prop")
		}
		values[key] = value
	}

	return values, nil
}

// Stringer interface
func (k *ApiKey) String() string {
	if k.raw != "" {
		return k.raw
	}

	return fmt.Sprintf("%s%s%s", k.signedPortion(), apiSeparator, k.signature)
}

func (k *ApiKey) signedPortion() string {
	if k.raw != "" {
		// Use the real string from the raw key if this was parsed.
		// Signature will include it for back compat, so can't skip this
		i := strings.LastIndex(k.raw, apiSeparator)
		if i <= 0 {
			return ""
		}
		return k.raw[:i-1]
	}

	// No raw string. This is new key. Build it.
	propsSections := []string{}
	for propKey, propVal := range k.props {
		propString := fmt.Sprintf("%s%s%s", propKey, propertyKeyDelimiter, propVal)
		b64 := base64.StdEncoding.EncodeToString([]byte(propString))
		propsSections = append(propsSections, b64)
	}
	propsSection := strings.Join(propsSections, apiSeparator)

	return fmt.Sprintf("%s%d%s%s", apiPrefix, k.version, apiSeparator, propsSection)
}
