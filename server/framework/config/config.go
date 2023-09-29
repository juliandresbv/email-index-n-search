package config

import (
	"github.com/joho/godotenv"

	customlogger "pkg/custom-logger"
)

var logger = customlogger.NewLogger()

func LoadEnvVars(path ...string) error {
	err := godotenv.Load(path...)

	if err != nil {
		logger.Fatalln("error loading .env file")

		return err
	}

	return nil
}
