package config
import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURI string
	JWTSecret   string
	Port        string
	AWSRegion   string
	S3BucketName string
	AWSAccessKeyID string
	AWSSecretAccessKey string
	Host			string
}

func Load() (Config, error) {
	// In development this loads P2/.env. In production, real environment
	// variables can be used without a .env file.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	cfg := Config{
		DatabaseURI: os.Getenv("POSTGRES_DB_URI"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        os.Getenv("PORT"),
		AWSRegion:   os.Getenv("AWS_REGION"),
		S3BucketName: os.Getenv("S3_BUCKET_NAME"),
		AWSAccessKeyID: os.Getenv("AWS_ACCESS_KEY"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Host: os.Getenv("HOST"),
	}

	if cfg.DatabaseURI == "" {
		return Config{}, fmt.Errorf("POSTGRES_DB_URI is not configured")
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is not configured")
	}
	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	return cfg, nil
}
