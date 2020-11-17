package utils

import (
	"fmt"
	"os"
	"strconv"
)

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func GetEnvString(params ...string) (string, error) {
	if value, exists := os.LookupEnv(params[0]); exists {
		return value, nil
	} else if len(params) > 1 {
		return params[1], nil
	} else {
		return "", fmt.Errorf("Environment variable not set: %v", params[0])
	}
}

func GetEnvInt(params ...string) (int, error) {
	if value, exists := os.LookupEnv(params[0]); exists {
		valInt, _ := strconv.Atoi(value)
		return valInt, nil
	} else if len(params) > 1 {
		valInt, _ := strconv.Atoi(params[1])
		return valInt, nil
	} else {
		return 0, fmt.Errorf("Environment variable not set: %v", params[0])
	}
}

func CheckRequiredEnvVars(vars []string) error {
	for _, v := range vars {
		val := os.Getenv(v)
		if val == "" {
			return fmt.Errorf("Environment variable not set: %v", v)
		}
	}

	return nil
}
