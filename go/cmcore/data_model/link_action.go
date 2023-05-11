package datamodel

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type LinkAction struct {
	UrlString string
}

type jsonLinkAction struct {
	UrlString string `json:"url"`
}

func unpackLinkFromJson(rawJson json.RawMessage, ac *ActionContainer) (ActionTypeInterface, error) {
	var link LinkAction
	err := json.Unmarshal(rawJson, &link)
	if err != nil {
		return nil, err
	}
	ac.LinkAction = &link
	return &link, nil
}

func (l *LinkAction) Validate() bool {
	return l.ValidateReturningUserReadableIssue() == ""
}

func (l *LinkAction) ValidateReturningUserReadableIssue() string {
	if l.UrlString == "" {
		return "Link actions must have a url"
	}
	url, err := url.Parse(l.UrlString)
	// We don't want to accept schemeless URLs ("/local/path")
	// We don't want to accept opaque URLs ("scheme:path") https://pkg.go.dev/net/url#URL
	if err != nil || url == nil || url.Scheme == "" || url.Opaque != "" {
		return fmt.Sprintf("Link action url string is not a valid URL: \"%v\"", l.UrlString)
	}

	return ""
}

func (l *LinkAction) UnmarshalJSON(data []byte) error {
	var jl jsonLinkAction
	err := json.Unmarshal(data, &jl)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an action with type=link. Check the format, variable names, and types (eg float vs int).", err)
	}

	l.UrlString = jl.UrlString

	if validationIssue := l.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (l *LinkAction) AllEmbeddedThemeNames() ([]string, error) {
	return []string{}, nil
}

func (l *LinkAction) AllEmbeddedActionNames() ([]string, error) {
	return []string{}, nil
}

func (l *LinkAction) PerformAction(ab ActionBindings) error {
	return ab.ShowLink(l)
}
