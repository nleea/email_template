package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

func ConfigEnv() map[string]string {
	envFile, err := godotenv.Read(".env")

	if err != nil {
		fmt.Println("There is a error reading the .env file")
	}

	return envFile
}
