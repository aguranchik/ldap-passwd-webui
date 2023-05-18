package app

import (
	"os"
	"strconv"
)

func getTitle() string {
	return envStr("LPW_TITLE", "Смена пароля в домене 21век")
}

func getPattern() string {
	return envStr("LPW_PATTERN", ".{8,}")
}

func getPatternInfo() string {
	return envStr("LPW_PATTERN_INFO", "Пароль должен состоять не менне чем из 8 символов и содержать 3 из 4 категорий символов: большие и маленькие буквы, цифры и спецсимволы (#$%!).")
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
