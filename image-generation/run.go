package image_generation

import (
	"fmt"
	"go-openai-exercises/client"
)

func RunExercise() {

	desc := client.GetRobotDescription()
	url := GenerateRobotImage(desc)
	fmt.Println("Generated image URL: ", url)

	centrala := client.Centrala{}
	resp := centrala.SendSolution("robotid", url)
	fmt.Println("Response ", "\n", resp)
}
