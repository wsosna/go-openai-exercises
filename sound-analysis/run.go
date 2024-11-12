package sound_analysis

import (
	"github.com/joho/godotenv"
	"log"
)

func RunExercise() {
	loadEnvFile()
	transcript()
	answer()
}

func loadEnvFile() {
	if godotenv.Load(".env") != nil {
		log.Println("Failed to load `.env` file. Using system environment variables.")
	}
}
