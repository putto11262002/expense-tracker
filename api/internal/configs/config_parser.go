package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigParser struct {
	EnvFile string
}

func NewConfigParser(envFile string) (*ConfigParser) {
	return &ConfigParser{
		EnvFile: envFile,
	}
}

func (cp *ConfigParser) LoadEnv() (error) {
	err :=  godotenv.Load(cp.EnvFile)
	if err != nil {
		return fmt.Errorf("loading %s: %w", cp.EnvFile, err)
	}
	return nil
}


func (cp *ConfigParser) GetStringEnv(key string) (string, error) {
	val, ok :=  os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("env does not exist: %s", key)
	}
	return val, nil
}


func (cp *ConfigParser) GetIntEnv(key string) (int64, error) {
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

func (cp *ConfigParser) GetBoolEnv(key string) (bool, error) {
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