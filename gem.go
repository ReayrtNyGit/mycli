package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type RequestPayload struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponsePayload struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type APIConfig struct {
	URL   string `json:"url"`
	Model string `json:"model"`
}

type Config struct {
	APIConfigs   []APIConfig `json:"api_configs"`
	MessagesList [][]Message `json:"messages_list"`
	Temperatures []float64   `json:"temperatures"`
}

var client = &http.Client{}

func fetchResponse(apiConfig APIConfig, messages []Message, temperature float64, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()

	apiKey := os.Getenv("API_KEY_" + apiConfig.Model) // Environmental variable as API key
	if apiKey == "" {
		log.Printf("API key for model %s is not set.", apiConfig.Model)
		return
	}

	data := RequestPayload{
		Model:       apiConfig.Model,
		Messages:    messages,
		Temperature: temperature,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		return
	}

	req, err := http.NewRequest("POST", apiConfig.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	apiCallDuration := time.Since(startTime)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}

		var responsePayload ResponsePayload
		err = json.Unmarshal(bodyBytes, &responsePayload)
		if err != nil {
			log.Printf("Error unmarshalling response: %v", err)
			return
		}

		assistantMessage := responsePayload.Choices[0].Message.Content
		fmt.Println(assistantMessage)
	} else {
		log.Printf("Request failed with status code %d: %s", resp.StatusCode, resp.Status)
	}

	fmt.Printf("API call took %v\n", apiCallDuration)
}

func readConfig(filePath string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	return config, err
}

func main() {
	config, err := readConfig("config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var wg sync.WaitGroup
	startTime := time.Now()

	for _, apiConfig := range config.APIConfigs {
		for i, messages := range config.MessagesList {
			wg.Add(1)
			go fetchResponse(apiConfig, messages, config.Temperatures[i], &wg)
		}
	}

	wg.Wait()
	totalTime := time.Since(startTime)

	fmt.Printf("Total execution time: %v\n", totalTime)
}