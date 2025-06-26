package utils

import "os"

func GetEnvOr(key, defaultValue string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}

	return defaultValue
}
