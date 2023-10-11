package cmcore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

var loadEnvOnce sync.Once

func LoadEnv() {
	loadEnvOnce.Do(func() {
		envPath, err := filepath.Abs("../../.env")
		if err != nil {
			fmt.Printf("Error loading env: %v", err)
			return
		}
		if _, err := os.Stat(envPath); errors.Is(err, os.ErrNotExist) {
			// No .env is fine
			return
		}
		err = godotenv.Load(envPath)
		if err != nil {
			fmt.Println("Error loading .env file")
		}
	})
}
