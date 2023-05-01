package appcore

import (
	"github.com/CriticalMoments/CriticalMoments/go/cmcore"
)

func GoPing() string {
	return "AppcorePong->" + cmcore.CmCorePing()
}
