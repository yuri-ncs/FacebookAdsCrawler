package setups

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var envErrors []string

const realtimeApi = "realtime-api"

func SetupEnv() {
	validateEnv("POSTGRES_HOST")
	validateEnv("POSTGRES_PORT")
	validateEnv("POSTGRES_USER")
	validateEnv("POSTGRES_PASSWORD")
	validateEnv("POSTGRES_DBNAME")
	validateEnv("POSTGRES_SSLMODE")
	defaultEnv("POSTGRES_TIMEZONE", "disable")
	defaultEnv("POSTGRES_DEBUG", "America/Sao_Paulo")
	defaultEnv("POSTGRES_LOGLEVEL", "info")

	defaultEnv("CRON_PERIOD", "0 0 */4 * * *")

	defaultEnv("X_FB_LSD", "AVo7QBln3V0")
	defaultEnv("USER_AGENT", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	defaultEnv("HTTP_SERVER_PORT", ":8080")

	for _, err := range envErrors {
		fmt.Println(err)
	}

	if len(envErrors) > 0 {
		os.Exit(0)
	}
}

func validateEnv(envName string) {
	env := os.Getenv(envName)
	if env == "" {
		envErrors = append(envErrors, "no env", envName)
	}
}

func defaultEnv(envName, defaultValue string) {
	env := os.Getenv(envName)
	if env == "" {
		os.Setenv(envName, defaultValue)
	}
}
