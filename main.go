package main

import (
	"github.com/joho/godotenv"
	"go-openai-exercises/vector"
	"log"
)

// Runs the program.
func main() {
	loadEnvFile()

	//sound_analysis.RunExercise()
	//image_generation.RunExercise()
	//many_formats.RunExercise()
	//metadata.Run()
	vector.Run()
}

func loadEnvFile() {
	if godotenv.Load(".env") != nil {
		log.Println("Failed to load `.env` file. Using system environment variables.")
	}
}
