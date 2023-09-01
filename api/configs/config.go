package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		return fmt.Errorf("loading %s: %w", file, err)
	}
	return nil
}

func GetStringEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("env does not exist: %s", key)
	}
	return val, nil
}

func GetIntEnv(key string) (int64, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("env does not exist: %s", key)
	}
	parsedVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s to int: %w", key, err)
	}
	return parsedVal, nil
}

func GetGoEnv() string {
	goEnv := os.Getenv("GO_ENV")
	switch goEnv {
	case "production":
		return "production"
	case "development":
		return "development"
	default:
		return "development"
	}
}

func GetBoolEnv(key string) (bool, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return false, fmt.Errorf("env does not exist: %s", key)
	}
	parsedVal, err := strconv.ParseBool(val)
	if err != nil {
		return false, fmt.Errorf("parsing %s to bool: %w", key, err)
	}
	return parsedVal, nil
}
