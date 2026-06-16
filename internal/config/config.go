package config

import "os"

var (
	DBPath  = getEnv("ECHOID_DB_PATH", "fingerprints.db")
	TempDir = getEnv("ECHOID_TEMP_DIR", "temp")
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
