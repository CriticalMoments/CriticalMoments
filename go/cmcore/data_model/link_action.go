package datamodel

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type LinkAction struct {
	UrlString          string
	UseEmbeddedBrowser bool
}

type jsonLinkAction struct {
	UrlString          string `json:"url"`
	UseEmbeddedBrowser *bool  `json:"useEmbeddedBrowser,omitempty"`
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

func (l *LinkAction) Valid() bool {
	return l.Check() == nil
}

func (l *LinkAction) Check() UserPresentableErrorInterface {
	if l.UrlString == "" {
		return NewUserPresentableError("Link actions must have a url")
	}
	url, err := url.Parse(l.UrlString)
	// We don't want to accept schemeless URLs ("/local/path")
	// We do accept "Opaque" URLs as iOS uses this ("app-settings:root=Sounds") https://pkg.go.dev/net/url#URL
	if err != nil {
		return NewUserPresentableErrorWSource(fmt.Sprintf("Link action url string is not a valid URL: \"%v\"", l.UrlString), err)
	} else if url == nil || url.Scheme == "" {
		return NewUserPresentableError(fmt.Sprintf("Link action url string is not a valid URL: \"%v\"", l.UrlString))
	}
	// Embedded browser option only available for scheme http(s)
	if l.UseEmbeddedBrowser {
		if url.Scheme != "https" && url.Scheme != "http" {
			return NewUserPresentableError("For a link action, useEmbeddedBrowser is set to true, but the link is not a http/https link. Only web links can be opened in the embedded browser. Either disable useEmbeddedBrowser, or change the link to a web link.")
		}
	}

	return nil
}

func (l *LinkAction) UnmarshalJSON(data []byte) error {
	var jl jsonLinkAction
	err := json.Unmarshal(data, &jl)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an action with type=link. Check the format, variable names, and types (eg float vs int).", err)
	}

	// Defaults
	useEmbeddedBrowser := false
	if jl.UseEmbeddedBrowser != nil {
		useEmbeddedBrowser = *jl.UseEmbeddedBrowser
	}

	l.UrlString = jl.UrlString
	l.UseEmbeddedBrowser = useEmbeddedBrowser

	return l.Check()
}

func (l *LinkAction) AllEmbeddedThemeNames() ([]string, error) {
	return []string{}, nil
}

func (l *LinkAction) AllEmbeddedActionNames() ([]string, error) {
	return []string{}, nil
}

func (l *LinkAction) AllEmbeddedConditions() ([]*Condition, error) {
	return []*Condition{}, nil
}

func (l *LinkAction) PerformAction(ab ActionBindings, actionName string) error {
	return ab.ShowLink(l)
}
