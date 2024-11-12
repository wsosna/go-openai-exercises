package sound_analysis

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
	"go-openai-exercises/utils"
	"log"
	"os"
	"time"
)

func getAndValidateVariables() (string, string) {
	dir, dirEx := os.LookupEnv("TRANSCRIPTION_DIR_PATH")
	tmp, tmpEx := os.LookupEnv("TMP_DIR")

	if !dirEx {
		log.Fatal("Environment variable `TRANSCRIPTION_DIR_PATH` is required.")
	}
	if !tmpEx {
		tmp = "tmp"
	}
	return dir, tmp
}

// transcriptAndSaveToFiles transcribes all mp3 files in the directory specified by TRANSCRIPTION_DIR_PATH.
func transcriptAndSaveToFiles() {
	dir, tmp := getAndValidateVariables()
	// Header
	start := time.Now()
	thoughtsPath := fmt.Sprintf("%v/sound-analysis/thoughts-%v.md", tmp, start.Format("20060102150405"))
	log.Println(fmt.Sprintf("Begin analisis of mp3 files in directory: `%v`", dir))
	log.Println(fmt.Sprintf("All models thoughts are saved in `%v`", thoughtsPath))

	client := openai.NewClient()
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		outputPath := fmt.Sprintf("%v/sound-analysis/output-%v.txt", tmp, entry.Name())
		log.Println(fmt.Sprintf("Reading: `%v`", entry.Name()))
		reader, err := utils.ReadFileToBuffer(dir + "/" + entry.Name())
		utils.HandleFatalError(err)
		body := openai.AudioTranscriptionNewParams{
			File:     openai.FileParam(reader, entry.Name(), "audio/mp4"),
			Model:    openai.F(openai.AudioModelWhisper1),
			Language: openai.F("pl"),
			Prompt:   openai.F(SystemPromptTranscript),
		}
		transcription, err := client.Audio.Transcriptions.New(context.TODO(), body)
		utils.HandleFatalError(err)

		err = utils.WriteStringToFile(transcription.Text, outputPath)
		utils.HandleFatalError(err)
		log.Println(fmt.Sprintf("Transcription saved to: `%v`", outputPath))
	}

	// Footer
	stop := time.Now()
	log.Println(fmt.Sprintf(
		"Finished analisis, started: `%v`, stopped: `%v`, took: `%v`ms",
		start.Format("2006-01-02-15:04:05:456"),
		stop.Format("2006-01-02-15:04:05:456"),
		stop.Sub(start).Milliseconds()))
}
