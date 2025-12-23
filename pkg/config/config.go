package config

import (
	"os"
	"strconv"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port                int
	GracefulShutdownSec int
}

type AppConfig struct {
	DB       DBConfig
	Server   ServerConfig
	Auth     AuthConfig
	LiveKit  LiveKitConfig
	AWS      AWSConfig
	Gemini   GeminiConfig
	LogLevel string
	Env      string
}

type AuthConfig struct {
	JwksURL string
}



type LiveKitConfig struct {
	Host      string
	APIKey    string
	APISecret string
}

type AWSConfig struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
}

type GeminiConfig struct {
	RealtimeModel string
	ChatModel     string
	APIKey        string
}

func LoadConfig() (*AppConfig, error) {
	portStr := os.Getenv("DB_PORT")
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		portInt = 5432
	}
	config := &AppConfig{
		DB: DBConfig{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST"),
			Port:     portInt,
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_DATABASE"),
		},
		Server: ServerConfig{
			Port:                9000,
			GracefulShutdownSec: 5,
		},
		Auth: AuthConfig{
			JwksURL: os.Getenv("JWKS_URL"),
		},
		LiveKit: LiveKitConfig{
			Host:      os.Getenv("LK_HOST"),
			APIKey:    os.Getenv("LK_API_KEY"),
			APISecret: os.Getenv("LK_API_SECRET"),
		},
		AWS: AWSConfig{
			AccessKey: os.Getenv("AWS_ACCESS_KEY"),
			SecretKey: os.Getenv("AWS_SECRET_KEY"),
			Region:    os.Getenv("AWS_REGION"),
			Bucket:    os.Getenv("AWS_S3_BUCKET"),
		},
		Gemini: GeminiConfig{
			RealtimeModel: os.Getenv("GEMINI_REALTIME_MODEL"),
			ChatModel:     os.Getenv("GEMINI_CHAT_MODEL"),
			APIKey:        os.Getenv("GEMINI_API_KEY"),
		},
		LogLevel: "info",
		Env:      os.Getenv("APP_ENV"),
	}
	return config, nil
}
