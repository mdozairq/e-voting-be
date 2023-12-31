package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/mdozairq/e-voting-be/utils"
)

type ServerConfig struct {
	Port              string
	ServerApiPrefixV1 string
	JwtSecret         string
	BasePath          string
}

type DBConfig struct {
	MongoUri string
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
}

type TwilioConfig struct {
	Sid       string
	AuthToken string
}

type AdminConfig struct {
	AdminEmail string
	AdminPassword string
	AdminAuthToken string
}

// NewServerConfig returns a pointer to a new ServerConfig struct initialized with values from environment variables.
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:              os.Getenv("APP_PORT"),
		ServerApiPrefixV1: os.Getenv("SERVER_API_PREFIX_V1"),
		JwtSecret:         os.Getenv("JWT_SECRET"),
		BasePath:          os.Getenv("SERVER_BASE_PATH"),
	}
}

// NewDBConfig returns a pointer to a new DBConfig struct initialized with values from environment variables.
func NewDBConfig() *DBConfig {
	return &DBConfig{
		MongoUri: os.Getenv("MONGODB_URI"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}
}

func NewTwilioConfig() *TwilioConfig {
	return &TwilioConfig{
		Sid:       os.Getenv("TWILIO_SID"),
		AuthToken: os.Getenv("TWILIO_AUTH_TOKEN"),
	}
}

func NewAdminConfig() *AdminConfig {
	return &AdminConfig{
		AdminEmail: os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		AdminAuthToken: os.Getenv("ADMIN_AUTH_TOKEN"),
	}
}

// LoadEnv loads environment variables from the .env file in the current directory.
func LoadEnv() {

	loadEnvError := godotenv.Load(".env")
	if loadEnvError != nil {
		utils.LogFatal(loadEnvError)
	}
}
