package vector

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/pkoukk/tiktoken-go"
	"github.com/qdrant/go-client/qdrant"
	"go-openai-exercises/client"
	"go-openai-exercises/utils"
	"log"
	"os"
	"strconv"
)

const model = openai.ChatModelGPT4o

func Run() {

}

func tokenize(text string) []int {
	tke, err := tiktoken.EncodingForModel(model)
	utils.HandleFatalError(err)
	token := tke.Encode(text, nil, nil)
	//tokens
	fmt.Println(token)
	// num_tokens
	fmt.Println(len(token))
	return token
}

func index() {
	openAiClient := openai.NewClient()
	params := openai.EmbeddingNewParams{
		Model:          openai.F(openai.EmbeddingModelTextEmbedding3Large),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatBase64),
	}
	response, err := openAiClient.Embeddings.New(context.Background(), params)
	utils.HandleFatalError(err)

	port, err := strconv.Atoi(os.Getenv("QDRANT_PORT"))
	utils.HandleFatalError(err)
	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host:   os.Getenv("QDRANT_HOST"),
		Port:   port,
		APIKey: os.Getenv("QDRANT_API_KEY"),
	})
	utils.HandleFatalError(err)

	points := make([]*qdrant.PointStruct, 0)
	for _, input := range response.Data {
		points = append(points, &qdrant.PointStruct{
			Vectors: input.Embedding,
			Payload: input.Text,
		})
	}
	upsert, err := qdrantClient.GetPointsClient().Upsert(context.Background(), nil)
	utils.HandleFatalError(err)
	log.Println("Upsert response:", upsert.GetResult().Status)
}

func sendSolution(answer string) {
	centrala := client.Centrala{}
	resp := centrala.SendSolution("wektory", answer)
	log.Println("Response from centrala:", "\n", resp)
}
