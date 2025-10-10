package config

import (
	"fmt"
	"os"

	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"github.com/joho/godotenv"
)

// Function to get all the .env variables
func LoadConfig() *modals.Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load the environment file.")
		return &modals.Config{
			Port: "NOT DEFINED",
			Env: "NOT DEFINED",
		}
	}

	//Getting the port variable
	Port := os.Getenv("PORT")
	if Port == ""{
		Port = "8080"
	}

	//getting the env variable
	env := os.Getenv("ENV")
	if env == ""{
		env = "development"
	}

	return &modals.Config{
		Port: Port,
		Env: env,
	}
}