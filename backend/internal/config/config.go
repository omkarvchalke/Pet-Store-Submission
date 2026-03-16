package config

import "os"

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	ImageDir   string
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load() Config {
	return Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "petstore"),
		DBUser:     getEnv("DB_USER", "petuser"),
		DBPassword: getEnv("DB_PASSWORD", "petpass"),
		ImageDir:   getEnv("IMAGE_DIR", "./media/pet-images"),
	}
}
