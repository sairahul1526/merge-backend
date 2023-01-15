package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadConfig - loads all env vars
func LoadConfig() error {
	// load .env file from given path for local, else will be getting from env var
	if !strings.EqualFold(os.Getenv("prod"), "true") {
		configFile := ".test-env"
		if strings.EqualFold(os.Getenv("TESTING"), "true") {
			configFile = "../../.test-env"
		}
		err := godotenv.Load(configFile)
		if err != nil {
			return err
		}
	}

	// main postgres db
	MainDBConfig = os.Getenv("MAIN_DB_CONFIG")
	MainDBConnectionPool, _ = strconv.Atoi(os.Getenv("MAIN_DB_CONNECTION_POOL"))
	Log, _ = strconv.ParseBool(os.Getenv("LOG"))
	Migrate, _ = strconv.ParseBool(os.Getenv("MIGRATE"))
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	// s3

	// razorpay

	return nil
}
