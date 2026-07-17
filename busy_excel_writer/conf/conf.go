package conf

import "os"

// String returns the environment variable value if set, or the defaultVal.
func String(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
