package env

import (
	"os"
)

// LoadEnv loads the key Environment Variable into dest if exists.
func LoadEnv(dest *string, key string) {
	v := os.Getenv(key)
	if len(v) != 0 {
		*dest = v
	}
}
