package main

import (
	"fmt"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore"
)

func main() {
	// Get a greeting message and print it.
	message := cmcore.CmCorePing()
	fmt.Println(message)
}
