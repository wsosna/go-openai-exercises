package my_ai

import (
	"context"
	"github.com/openai/openai-go"
	"go-openai-exercises/utils"
)

type IOpenAiWrapper interface {
	AskMyAI(system string, user string) string
}

func NewOpenAiWrapper(chatModel openai.ChatModel) *OpenAiWrapper {
	return &OpenAiWrapper{
		Client: openai.NewClient(),
		Model:  chatModel,
	}
}

type OpenAiWrapper struct {
	Client *openai.Client
	Model  openai.ChatModel
}

func (w *OpenAiWrapper) AskMyAI(system string, user string) string {
	client := w.Client

	body := openai.ChatCompletionNewParams{
		Model: openai.F(w.Model),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(system),
			openai.UserMessage(user),
		}),
	}
	response, err := client.Chat.Completions.New(context.Background(), body)
	utils.HandleFatalError(err)
	return response.Choices[0].Message.Content
}
