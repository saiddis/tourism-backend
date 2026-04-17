package config

import "os"

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	ServerPort         string
	AccessTokenSecret  string
	RefreshTokenSecret string
	CookieDomain       string
	UploadDir          string
}

func Load() *Config {
	return &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "tourism_db"),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "tourism"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "tourism"),
		CookieDomain:       getEnv("COOKIE_DOMAIN", ""),
		UploadDir:          getEnv("UPLOAD_DIR", "./uploads/avatars"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
