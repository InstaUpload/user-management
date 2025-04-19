package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetEnvInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetSuperUsers() []string {
	superUsers := GetEnvString("SUPERUSERS", "")
	superUsers = strings.TrimSpace(superUsers)
	if superUsers == "" {
		return []string{}
	} else {
		return strings.Split(superUsers, ",")
	}
}
