package config

import (
	"go.uber.org/zap"
	"os"
	"strconv"
	"varnish-cache-invalidator/pkg/logging"
)

var logger = logging.GetLogger()

func convertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Warn("An error occurred while converting from string to int, setting it as zero.", zap.String("string", s))
		i = 0
	}
	return i
}

func convertStringToBool(s string) bool {
	i, err := strconv.ParseBool(s)
	if err != nil {
		logger.Warn("An error occurred while converting from string to bool, setting it as false.", zap.String("string", s))
		i = false
	}
	return i
}

// GetStringEnv gets the specific environment variables with default value, returns default value if variable not set
func GetStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// GetIntEnv gets the specific environment variables with default value, returns default value if variable not set
func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return convertStringToInt(value)
}

// GetBoolEnv gets the specific environment variables with default value, returns default value if variable not set
func GetBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return convertStringToBool(value)
}
