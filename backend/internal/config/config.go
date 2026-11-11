package config

import (
	"os"

	"github.com/joho/godotenv"
)

type MaxHelperStructConfig struct{
	MAX_token string
	MAX_deviceID string
	URL_webSocket string // по умолчанию "wss://api.oneme.ru/websocket"
}

type ConfigStruct struct{
	HTTP string//по умолчанию ":8080"
	ConnPgrs string
	MAX MaxHelperStructConfig
	URL_EndPoint string//по умолчанию "http://localhost:8080/changes/import-text"
}

func getEnv(key, defaultVal string) string{
	val := os.Getenv(key)
	
	if val == ""{
		return defaultVal
	}

	return val
}

func LoadConfig() *ConfigStruct{
	_ = godotenv.Load()

	var max MaxHelperStructConfig
	max.MAX_token = getEnv("MAX_TOKEN", "")
	max.MAX_deviceID = getEnv("MAX_DEVICE_ID", "")
	max.URL_webSocket = getEnv("MAX_WEB_SOCKET", 
	"wss://api.oneme.ru/websocket")

	var config ConfigStruct
	config.HTTP = getEnv("PORT", ":8080")
	config.ConnPgrs = getEnv("DATABASE_URL",
	 "postgres://postgres:postgres@localhost:5432/subly_db?sslmode=disable")
	config.URL_EndPoint = getEnv("SUBLY_API_URL", 
	"http://localhost:8080/changes/import-text")
	config.MAX = max

	return &config
}