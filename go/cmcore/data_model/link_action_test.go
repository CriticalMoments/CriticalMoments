package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLinkActionValidators(t *testing.T) {
	l := LinkAction{}
	if l.Validate() {
		t.Fatal("Links require a url")
	}
	l.UrlString = "This_isnt_a_url"
	if l.Validate() {
		t.Fatal("Links require a valid url and should not validate")
	}
	l.UrlString = "app-settings:root=Photos"
	if !l.Validate() {
		t.Fatal("Link vaidation failed for valid opaque url")
	}
	l.UrlString = "/Local/Urls/Dont/Count"
	if l.Validate() {
		t.Fatal("Links require a valid scheme and should not validate")
	}
	l.UrlString = "../Relative/Urls/Dont/Count"
	if l.Validate() {
		t.Fatal("Links require a valid scheme and should not validate")
	}
	l.UrlString = "https://scosman.net/asdf?t=5"
	if !l.Validate() {
		t.Fatal(l.ValidateReturningUserReadableIssue())
	}
	l.UrlString = "custom://any_scheme_is_okay/asdf.ext"
	if !l.Validate() {
		t.Fatal(l.ValidateReturningUserReadableIssue())
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
