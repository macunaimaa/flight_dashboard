package config

import (
	"os"
	"time"
)

type Config struct {
	Port                string
	MongoURI            string
	MongoDB             string
	RedisAddr           string
	JWTSecret           string
	OpenSkyPollInterval time.Duration
	OpenSkyUsername     string
	OpenSkyPassword     string
	LogLevel            string
	AllowedOrigins      string
	UseMockData         bool
	DataSource          string
	NewsAPIKey          string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		MongoURI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:             getEnv("MONGO_DB", "dashboard"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:           getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		OpenSkyPollInterval: parseDuration(getEnv("OPENSKY_POLL_INTERVAL", "10s")),
		OpenSkyUsername:     getEnv("OPENSKY_USERNAME", ""),
		OpenSkyPassword:     getEnv("OPENSKY_PASSWORD", ""),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		AllowedOrigins:      getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
		UseMockData:         getEnv("USE_MOCK_DATA", "false") == "true",
		DataSource:          getEnv("DATA_SOURCE", "adsblol"),
		NewsAPIKey:          getEnv("NEWS_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 10 * time.Second
	}
	return d
}
