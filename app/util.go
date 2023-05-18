package app

import (
	"os"
	"strconv"
)

func getTitle() string {
	return envStr("LPW_TITLE", "Change password in 21vek domain")
}

func getPattern() string {
	return envStr("LPW_PATTERN", ".{8,}")
}

func getPatternInfo() string {
	return envStr("LPW_PATTERN_INFO", "The password must be at least 8 characters and contain uppercase and lowercase letters, numbers, special characters (3 out of 4 categories).")
}

func envStr(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultValue
}

func envInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			return defaultValue
		}
		return i
	}
	return defaultValue
}

func envBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return defaultValue
		}
		return b
	}
	return defaultValue
}
