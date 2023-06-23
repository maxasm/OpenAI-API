package main

import (
	"bufio"
	"errors"
	"strings"
	"time"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func chat(prompt string) (string, error) {
	if len(strings.Trim(prompt, " ")) == 0 {
		return "", errors.New("You can't give an empty prompt! (Also, this is how you exit 😉)")
	}
	// The API endpoint
	const URL = "https://api.openai.com/v1/chat/completions"

	timeout,err_timeout := time.ParseDuration("10s")
	if err_timeout != nil {
		log.Fatalf("%s\n", err_timeout)
	}

	// Create a client with a timeout of 10m
	client := &http.Client{Timeout:timeout}

	req, err_req := http.NewRequest("POST", URL, nil)
	if err_req != nil {
		log.Fatalf("%s\n", err_req)
	}

	// Add the required Headers
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	// Create the request
	api_req := OPEN_AI_API_REQUEST{
		Model: GPT3,
		Messages: []OPEN_AI_API_MESSAGE{
			{Role: "system", Content: "You are a helpful assistant who follows the given prompt. Use emojis where possible."},
			{Role: "user", Content:prompt},
		},
	}

	// Encode request to json
	json_str, err_json_str := json.Marshal(api_req)
	if err_json_str != nil {
		log.Fatalf("%s\n", err_json_str)
	}

	// You need an io.ReadCloser
	str_read := bytes.NewBuffer(json_str)
	str_read_closer := io.NopCloser(str_read)
	req.Body = str_read_closer

	// Send the request
	resp, err_resp := client.Do(req)
	if err_resp != nil {
		log.Fatalf("%s\n", err_resp)
	}

	// Read the response as []byte
	data, err_data := io.ReadAll(resp.Body)
	if err_data != nil {
		log.Fatalf("%s\n", err_data)
	}

	// Use json.Unmarshal to convert the json-encoded string to an object 
	resp_obj := &OPEN_AI_API_RESPONSE{}
	err_json := json.Unmarshal(data,resp_obj)
	if err_json != nil {
		log.Fatalf("%s\n",err_json)
	}	
	
	// Print out the response
	msg := resp_obj.Choices[0].Message
	return msg.Content, nil	
}

// Todo: Concatinate the messages
// Todo: Block the input when the system is waiting for the reponse
// Todo: Implement the CTRL + L to clear the screen
// Todo: Implment scrolling up to past prompts

func main() {
	// Simple chat bot on the terminal
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("BOT > Hey 👋, how can I help you.\n")
	fmt.Printf("You > ")

	for scanner.Scan() {
		// Read input from the user
		user_prompt := scanner.Text()
		bot_response, err_bot_response := chat(user_prompt)
		
		if err_bot_response != nil {
			log.Fatalf("%s\n", err_bot_response)
		}
	
		fmt.Printf("BOT > %s\n", bot_response)
		fmt.Printf("You > ")
	}
}
