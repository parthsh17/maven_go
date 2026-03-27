package config

import (
	"os"
	"strconv"
)

type Config struct {
	Addr         string
	MongoURI     string
	DBName       string
	WorkerCount  int
	BufferSize   int
	SuccessRate  float64
}

func Load() *Config {
	return &Config{
		Addr:        getEnv("PORT", ":8080"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:      getEnv("DB_NAME", "maven"),
		WorkerCount: getEnvAsInt("WORKER_COUNT", 5),
		BufferSize:  getEnvAsInt("BUFFER_SIZE", 100),
		SuccessRate: getEnvAsFloat("SUCCESS_RATE", 0.70),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsFloat(key string, fallback float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return fallback
}
