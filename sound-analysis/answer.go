package sound_analysis

import (
	"context"
	"github.com/openai/openai-go"
	"go-openai-exercises/utils"
	"log"
	"os"
)

func answer() {
	client := openai.NewClient()
	recordings := readRecordings()
	message := createMessage(recordings)

	body := openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModelGPT4),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(SystemPromptAnalyser),
			openai.UserMessage(message),
		}),
	}

	response, err := client.Chat.Completions.New(context.Background(), body)
	utils.HandleFatalError(err)
	log.Println("Received answer from a model:")
	log.Println(response.Choices[0].Message.Content)
}

// Answer This function assumes that transcription was already executed and will look for files in the directory.
func readRecordings() string {
	dir := getSoundAnalysisDir()
	entries, err := os.ReadDir(dir)
	utils.HandleFatalError(err)
	result := ""
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// Read the file
		content, err := utils.ReadFileToString(dir + "/" + entry.Name())
		utils.HandleFatalError(err)

		// Append the file to string
		if len(content) > 0 {
			result += "<recording>\n"
			result += content
			result += "\n</recording>\n"
		}
	}
	log.Println("Read ", len(entries), " recordings from directory: ", dir)
	return result
}

func createMessage(recordings string) string {
	result := ""
	result += "<message>\n"
	result += "Analyse recordings and use your own knowledge, and figure out the address of university where professor Andrzej Maj has lectures.\n"
	result += "</message>\n"
	result += "<context>\n<recordings>\n"
	result += recordings
	result += "</context>\n</recordings>\n"
	return result
}

func getSoundAnalysisDir() string {
	tmp, tmpEx := os.LookupEnv("TMP_DIR")
	if !tmpEx {
		tmp = "tmp"
	}
	return tmp + "/sound-analysis"
}
