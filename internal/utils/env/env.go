package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return num
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}

	return duration
}
