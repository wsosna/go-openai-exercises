package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-openai-exercises/utils"
	"io"
	"log"
	"net/http"
	"os"
)

type ICentrala interface {
	SendSolution(task string, answer any) string
}

type Centrala struct {
}

func (c *Centrala) SendSolution(task string, answer any) string {
	req := TaskRequest{Task: task, ApiKey: os.Getenv("AI_DEVS_3_API_KEY"), Answer: answer}
	body, _ := json.Marshal(req)
	bodyReader := bytes.NewReader(body)

	resp, err := http.Post("https://centrala.ag3nts.org/report", "application/json", bodyReader)
	utils.HandleFatalError(err)
	log.Println("Received response status " + resp.Status)

	respBody, err := io.ReadAll(resp.Body)
	utils.HandleFatalError(err)

	return string(respBody)
}

func GetRobotDescription() string {
	resp, err := http.Get(os.Getenv("ROBOT_DESCRIPTION_URL"))
	utils.HandleFatalError(err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	utils.HandleFatalError(err)
	description := string(respBody)

	fmt.Println("Robots description response ", resp.Status, "\n", description)
	return description
}

type TaskRequest struct {
	Task   string `json:"task"`
	ApiKey string `json:"apikey"`
	Answer any    `json:"answer"`
}
