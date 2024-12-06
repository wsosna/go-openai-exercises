package many_formats

import (
	"encoding/json"
	"fmt"
	"go-openai-exercises/client"
	image_generation "go-openai-exercises/image-generation"
	my_ai "go-openai-exercises/my-ai"
	sound_analysis "go-openai-exercises/sound-analysis"
	"go-openai-exercises/utils"
	"os"
	"strings"
)

const AllFormatsDir = "tmp/pliki_z_fabryki"
const OutputDir = "tmp/many-formats"

type Labels struct {
	People   []string `json:"people"`
	Hardware []string `json:"hardware"`
}

func RunExercise() {
	//transcriptAllSoundsFilesAndSave(false)
	answer := analyseFiles()

	answer = strings.ReplaceAll(answer, ".png.txt", ".png")
	answer = strings.ReplaceAll(answer, ".mp3.txt", ".mp3")
	fmt.Println("Req:\n" + answer)

	var jsonMap Labels
	err := json.Unmarshal([]byte(answer), &jsonMap)
	utils.HandleFatalError(err)

	centrala := client.Centrala{}
	resp := centrala.SendSolution("kategorie", jsonMap)
	fmt.Println("Response ", "\n", resp)
}

func transcriptAllSoundsFilesAndSave(operations ...bool) {
	transcribe, convert := readOptionals(operations...)
	dir, err := os.ReadDir(AllFormatsDir)
	utils.HandleFatalError(err)

	for _, entry := range dir {
		if transcribe && strings.HasSuffix(entry.Name(), ".mp3") {
			fmt.Println("Transcribing: ", entry.Name())
			text := sound_analysis.Transcript(AllFormatsDir+"/"+entry.Name(), "en")
			err := utils.WriteStringToFile(text, OutputDir+"/"+entry.Name()+".txt")
			utils.HandleFatalError(err)
			fmt.Println("Done ✓")
			continue
		}
		if convert && strings.HasSuffix(entry.Name(), ".png") {
			// todo: this could be just one call with many <images> in the message
			fmt.Println("Converting image to text: ", entry.Name())
			text := image_generation.ImageToString(AllFormatsDir + "/" + entry.Name())
			err := utils.WriteStringToFile(text, OutputDir+"/"+entry.Name()+".txt")
			utils.HandleFatalError(err)
			fmt.Println("Done ✓")
		}
	}

}

func analyseFiles() string {
	filesPrompt := readFilesToPrompt()
	systemPrompt := "<system> You are helpful investigator who is analyzing files from the factory. Your main ability if to categorize and label information.</system>"
	message :=
		`<message>Please analyze the following files and determine if each report was created by a robot. If a report was created by a robot, categorize the files into the following categories: "people" and "hardware".
			Ignore all other files. If a file was not written by a robot and does not contain information about captured people, ignore it. If a file contains information about captured people, categorize it as "people". 
			If a file contains information about hardware, categorize it as "hardware". Ignore files that contain information about software, code, or AI. Also, ignore files that mention searching for people without finding any.
			Respond only with the result in JSON format, as shown in the example below:
			<examples>
				{
				  "people": ["file1.txt", "file2.mp3", "fileN.png"],
				  "hardware": ["file4.txt", "file5.png", "file6.mp3"],
				  "ignored": ["file3.txt"]
				}
			</examples>
		<files> ` + filesPrompt + ` </files>
		</message>`

	ai := my_ai.NewOpenAiWrapper()
	answer := ai.AskMyAI(systemPrompt, message)
	return answer
}

func readFilesToPrompt() string {
	dir, err := os.ReadDir(OutputDir)
	utils.HandleFatalError(err)

	result := ""
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		content, err := utils.ReadFileToString(OutputDir + "/" + entry.Name())
		utils.HandleFatalError(err)

		content = strings.ReplaceAll(content, "\n", " ")
		if len(content) > 0 {
			result += "<file>\n"
			result += "<name>\n"
			result += entry.Name()
			result += "\n</name>\n"
			result += "<content>\n"
			result += content
			result += "\n</content>\n"
			result += "</file>\n"
		}
	}
	return result

}

// readOptionals reads two optional boolean values. If only one is provided, the second one is set to true.
// If none are provided, both are set to true.
// This is a way to limit execution of some operations when calling a function.
func readOptionals(optionals ...bool) (bool, bool) {
	length := len(optionals)
	if length == 0 {
		return true, true
	}
	if length == 1 {
		return optionals[0], true
	}
	return optionals[0], optionals[1]
}
