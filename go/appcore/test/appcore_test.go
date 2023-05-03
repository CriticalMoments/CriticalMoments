package testing

import (
	"testing"

	"github.com/CriticalMoments/CriticalMoments/go/appcore"
)

func TestPing(t *testing.T) {
	pingResponse := appcore.GoPing()
	if pingResponse != "AppcorePong->PongCmCore" {
		t.Fatalf("appcore ping failure: %v", pingResponse)
	}
}
