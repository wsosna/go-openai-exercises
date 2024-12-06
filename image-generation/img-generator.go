package image_generation

import (
	"context"
	"github.com/openai/openai-go"
	"go-openai-exercises/utils"
)

func GenerateRobotImage(description string) string {

	client := openai.NewClient()
	params := openai.ImageGenerateParams{
		Model:          openai.F(openai.ImageModelDallE3),
		Prompt:         openai.F(generateRobotPrompt(description)),
		Size:           openai.F(openai.ImageGenerateParamsSize1024x1024),
		ResponseFormat: openai.F(openai.ImageGenerateParamsResponseFormatURL),
		Style:          openai.F(openai.ImageGenerateParamsStyleVivid),
	}
	resp, err := client.Images.Generate(context.TODO(), params)
	utils.HandleFatalError(err)
	return resp.Data[0].URL
}

func generateRobotPrompt(description string) string {
	prompt := `
		<system>
		You are helpful police artist who is creating a sketch of a robot based on the description provided by a witness.
		</system>
		<message>
		Please generate an image of a robot based on the following description:
		<description>
		` + description + `
		</description>
		</message>
`
	return prompt
}

func ImageToString(file string) string {
	systemMessage := `
		<system>
			You are an AI that is analyzing an image and generating text that is present in the image.
			<rules>
				- Respond with ONLY the text that is present in the image.
			</rules>
		</system>
	`
	client := openai.NewClient()
	bas64, err := utils.ReadFileToBase64(file)
	utils.HandleFatalError(err)

	body := openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModelGPT4o),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMessage),
			openai.UserMessageParts(openai.ImagePart(bas64)),
		}),
	}

	response, err := client.Chat.Completions.New(context.Background(), body)
	utils.HandleFatalError(err)
	return response.Choices[0].Message.Content
}
