package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"fmt"
)

func main() {
	// Define a string flag with a default value and a short description.
	name := flag.String("name", "World", "a name to say hello to")

	// Define a string flag for the API key.
	apiKey := flag.String("apikey", "", "API key for authentication")
	flag.Parse()

	// Use the flag value in the program.
	fmt.Printf("Hello, %s!\n", *name)
	// Make an API request
	url := "https://api.openai.com/v1/chat/completions"
	reqBody := []byte(`{
		"model": "gpt-4o-mini",
		"messages": [{"role": "user", "content": "Say this is a test!"}],
		"temperature": 0.7
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+*apiKey)

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
