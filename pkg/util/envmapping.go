package util

import (
	"fmt"
	"os"
)

// MustMapEnv environment variable must be there
func MustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v
}
