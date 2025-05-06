package utils

import (
	"os"
	"strconv"
	"strings"
	"time"
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

func TimeToString(t *time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.DateTime)
}

func StringToTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
