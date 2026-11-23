package utils

import (
	"os"
	"strconv"
)

func Getenv[T ~string | ~int | ~int64 | ~bool | ~float64](name string, def T) T {
	value, found := os.LookupEnv(name)
	if !found {
		return def
	}

	switch any(def).(type) {
	case string:
		return any(value).(T)
	case int:
		v, err := strconv.Atoi(value)
		if err == nil {
			return any(v).(T)
		}
	case int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return any(v).(T)
		}
	case bool:
		v, err := strconv.ParseBool(value)
		if err == nil {
			return any(v).(T)
		}
	case float64:
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return any(v).(T)
		}
	}

	return def
}
