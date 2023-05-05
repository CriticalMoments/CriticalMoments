package datamodel

const BannerMaxLineCountSystemDefault = -1
const BannerMaxLineCountSystemUnlimited = 0

type BannerAction struct {
	Body              string
	ShowDismissButton bool
	MaxLineCount      int
	TapActionName     string
	Theme             string
}

func (ba BannerAction) Validate() bool {
	return ba.ValidateReturningUserReadableIssue() == ""
}

func (b BannerAction) ValidateReturningUserReadableIssue() string {
	if b.Body == "" {
		return "Banners must have body text"
	}
	if b.MaxLineCount != BannerMaxLineCountSystemDefault && b.MaxLineCount < 0 {
		// Technically -1 allowed, but that's an internal between cmcore and libraries
		// Not user facing or a value they should put in json or see in libraries
		return "Banner max line count must be a positive integer, or 0 for no limit"
	}

	return ""
}
