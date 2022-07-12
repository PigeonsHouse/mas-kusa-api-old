package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	ServerUrl  string
	Port       string
	DbHost     string
	DbUser     string
	DbPass     string
	DbName     string
	DbPort     string
	SigningKey []byte
)

func InitEnv() {
	godotenv.Load(".env")
	apiProtocol := getEnv("API_PROTOCOL", "http")
	apiUrl := getEnv("API_URL", "localhost")
	ServerUrl = fmt.Sprintf("%s://%s", apiProtocol, apiUrl)
	Port = getEnv("API_PORT", "8000")

	DbHost = getEnv("POSTGRES_HOST", "localhost")
	DbUser = getEnv("POSTGRES_USER", "user")
	DbPass = getEnv("POSTGRES_PASS", "password")
	DbName = getEnv("POSTGRES_DB", "postgres")
	DbPort = getEnv("POSTGRES_PORT", "5432")
	secretKey := getEnv("SECRET_KEY", "secret")

	SigningKey = []byte(secretKey)
}

func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
