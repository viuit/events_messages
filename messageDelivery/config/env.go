package config

import (
	"os"
	"strconv"
)

type Environment struct {
	Rabbit    rabbitEnv
	CipherKey string
}

type rabbitEnv struct {
	Host     string
	Port     int
	VHost    string
	Exchange string
}

func NewEnv() *Environment {
	return &Environment{
		Rabbit: rabbitEnv{
			Host:     getEnv("ASD_RMQ_HOST", ""),
			Port:     getEnvAsInt("ASD_RMQ_PORT", 5672),
			VHost:    getEnv("ASD_RMQ_VHOST", ""),
			Exchange: getEnv("SERVICE_RMQ_EXCHANGE", ""),
		},
		CipherKey: getEnv("ASD_CIPHER_KEY", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
