package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"moderator/models"
	"net/http"
	"os"
)

func Validate(requirements string, text string) (*models.GptResponse, *models.Error, error) {
	url := "https://api.openai.com/v1/chat/completions"
	messages := models.Message{
		Role: "system",
		Content: []models.Content{
			{
				Type: "text",
				Text: requirements,
			},
			{
				Type: "text",
				Text: text,
			},
		},
	}
	request := models.GptRequest{
		Model:            "gpt-4o",
		Messages:         []models.Message{messages},
		Temperature:      1,
		MaxTokens:        100,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		ResponseFormat: models.ResponseFormat{
			Type: "text",
		},
	}
	data, err := json.Marshal(request)
	if err != nil {
		log.Println("failed serializing request into json", err)
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Println("failed creating request", err)
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))

	client, err := GetClientWithSocks5Proxy()
	if err != nil {
		log.Println("failed creating client with socks5 proxy", err)
		return nil, nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("failed sending request", err)
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("failed closing body", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("failed reading response body", err)
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var response models.Error
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("failed parsing response body", err)
			return nil, nil, err
		}
		log.Println("failed validating request", string(body))
		return nil, &response, nil
	}

	var response models.GptResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("failed parsing response body", err)
		return nil, nil, err
	}

	return &response, nil, nil
}
