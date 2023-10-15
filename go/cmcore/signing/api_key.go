package signing

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ApiKey struct {
	// Raw strings making up the API Key.
	// Don't manipulate outside of Parse/New
	signedPortion string
	signature     string

	// Parsed Properties
	version int
	props   map[string]string
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
	return u.VerifyMessage([]byte(k.signedPortion), k.signature)
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
		return nil, errors.New("bundle ID required")
	}

	if u.privateKey == nil {
		return nil, errors.New("can not create new API keys without a private key")
	}

	key := ApiKey{
		version: currentApiKeyVersion,
		props: map[string]string{
			propKeyBundleId: bundleId,
		},
	}

	// save exact signedPortion string to it doesn't mutate
	sp := key.buildNewKeySignedPortion()
	key.signedPortion = sp

	sig, err := u.SignMessage([]byte(key.signedPortion))
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

	// Signed Portion is everything before last separator
	i := strings.LastIndex(s, apiSeparator)
	if i <= 0 {
		return nil, errors.New("invalid API Key")
	}
	sp := s[:i]

	k := ApiKey{
		version:       v,
		signature:     sig,
		props:         props,
		signedPortion: sp,
	}

	return &k, nil
}

func parseVersionNumber(s string) (int, error) {
	if strings.Index(s, apiPrefix) != 0 {
		return -1, errors.New("all API keys should start with CM")
	}
	i := strings.Index(s, apiSeparator)
	if i < len(apiPrefix)+1 {
		return -1, errors.New("no version number or invalid format")
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
		return nil, errors.New("invalid API key format")
	}

	values := map[string]string{}
	propsSection := s[propStart+1 : propEnd]
	for _, rawProp := range strings.Split(propsSection, apiSeparator) {
		propBytes, err := base64.StdEncoding.DecodeString(rawProp)
		if err != nil {
			return nil, errors.New("invalid API key format")
		}
		decodedProp := string(propBytes)
		pi := strings.Index(decodedProp, propertyKeyDelimiter)
		if pi <= 0 {
			return nil, errors.New("invalid API key format")
		}
		key := decodedProp[:pi]
		value := decodedProp[pi+1:]
		if key == "" || value == "" {
			return nil, errors.New("invalid API Key Prop")
		}
		values[key] = value
	}

	return values, nil
}

// Stringer interface
func (k *ApiKey) String() string {
	return fmt.Sprintf("%s%s%s", k.signedPortion, apiSeparator, k.signature)
}

func (k *ApiKey) buildNewKeySignedPortion() string {
	// Make prop order deterministic.
	// Doesn't tecnhically need to be, but will save confusion
	propKeys := make([]string, 0)
	for k := range k.props {
		propKeys = append(propKeys, k)
	}
	sort.Strings(propKeys)

	propsSections := []string{}
	for _, propKey := range propKeys {
		propVal := k.props[propKey]
		propString := fmt.Sprintf("%s%s%s", propKey, propertyKeyDelimiter, propVal)
		b64 := base64.StdEncoding.EncodeToString([]byte(propString))
		propsSections = append(propsSections, b64)
	}
	propsSection := strings.Join(propsSections, apiSeparator)

	return fmt.Sprintf("%s%d%s%s", apiPrefix, k.version, apiSeparator, propsSection)
}
