package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLinkActionValidators(t *testing.T) {
	l := LinkAction{}
	if l.Valid() {
		t.Fatal("Links require a url")
	}
	l.UrlString = "This_isnt_a_url"
	if l.Valid() {
		t.Fatal("Links require a valid url and should not validate")
	}
	l.UrlString = "app-settings:root=Photos"
	if !l.Valid() {
		t.Fatal("Link vaidation failed for valid opaque url")
	}
	l.UrlString = "/Local/Urls/Dont/Count"
	if l.Valid() {
		t.Fatal("Links require a valid scheme and should not validate")
	}
	l.UrlString = "../Relative/Urls/Dont/Count"
	if l.Valid() {
		t.Fatal("Links require a valid scheme and should not validate")
	}
	l.UrlString = "https://scosman.net/asdf?t=5"
	if !l.Valid() {
		t.Fatal(l.Check())
	}
	l.UrlString = "custom://any_scheme_is_okay/asdf.ext"
	if !l.Valid() {
		t.Fatal(l.Check())
	}
}

func TestLinkActionValidateEmbedded(t *testing.T) {
	l := LinkAction{
		UrlString: "app-settings:main",
	}
	if !l.Valid() {
		t.Fatal("Valid link failed to validate")
	}
	l.UseEmbeddedBrowser = true
	if l.Valid() {
		t.Fatal("Open embedded browser with no web url should fail")
	}
	l.UrlString = "https://scosman.net"
	if !l.Valid() {
		t.Fatal("Valid link failed to validate")
	}
	l.UrlString = "http://scosman.net"
	if !l.Valid() {
		t.Fatal("Valid link failed to validate")
	}
}

func TestJsonParsingValidLink(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/link/validLink.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}
	if ac.ActionType != ActionTypeEnumLink || ac.LinkAction.UrlString != "https://scosman.net" {
		t.Fatal("Failed to parse valid link action")
	}
	if ac.LinkAction.UseEmbeddedBrowser {
		t.Fatal("Default for embedded browser option should be false")
	}
}

func TestJsonParsingValidEmbeddedLink(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/link/validLinkEmbedded.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}
	if ac.ActionType != ActionTypeEnumLink || ac.LinkAction.UrlString != "https://scosman.net" {
		t.Fatal("Failed to parse valid link action")
	}
	if !ac.LinkAction.UseEmbeddedBrowser {
		t.Fatal("embedded browser option failed to parse")
	}
}
func TestJsonParsingInvalidLink(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/link/invalidLink.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err == nil || ac.ActionType == ActionTypeEnumLink {
		t.Fatal("Invalid links should not parse")
	}
}

func TestJsonParsingInvalidMissingUrlLink(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/link/invalidMissingLinkUrl.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err == nil || ac.ActionType == ActionTypeEnumLink {
		t.Fatal("Invalid links should not parse")
	}
}
