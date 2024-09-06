package main

import (
	"bytes"
	"flag"
	"os"
	"io/ioutil"
	"net/http"
	"fmt"
)

func main() {
	// Define a string flag with a default value and a short description.
	name := flag.String("name", "World", "a name to say hello to")

	// Define a string flag for the model with a default value and a short description.
	model := flag.String("model", "gpt-4o-mini", "the model to use (gpt-4o-mini or gpt-4o)")

	flag.Parse()

	// Validate the model choice
	if *model != "gpt-4o-mini" && *model != "gpt-4o" {
		fmt.Println("Error: Invalid model choice. Use 'gpt-4o-mini' or 'gpt-4o'.")
		return
	}

	// Use the flag value in the program to greet the user.
	fmt.Printf("Hello, %s!\n", *name)
	// Retrieve the API key from the environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable is not set.")
		return
	}
	url := "https://api.openai.com/v1/chat/completions"
	reqBody := []byte(`{
		"model": "` + *model + `",
		"messages": [{"role": "user", "content": "Say this is a test!"}],
		"temperature": 0.7
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response:", string(body))
}
