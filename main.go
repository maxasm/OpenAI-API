package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"io"
	"log"
	"net/http"
)

func main() {
	const URL = "https://api.openai.com/v1/chat/completions"

	client := &http.Client{}
	req, err_req := http.NewRequest("POST",URL, nil)	
	if err_req != nil {
		log.Fatalf("%s\n",err_req.Error())
	}
	
	// Add the required Headers
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	
	// sample request 
	api_req := OPEN_AI_API_REQUEST{
		Model: GPT3,
		Messages: []OPEN_AI_API_MESSAGE{
			{Role: "system", Content: "You are a football expert."},
			{Role: "user", Content: "When did Liverpool win the Champions League?"},
		},
	}
	
	// encode request to json
	json_str, err_json_str := json.Marshal(api_req)
	if err_json_str != nil {
		log.Fatalf("%s\n",err_json_str.Error())
	}
	
	// req.Body -> io.ReadCloser
	str_read := bytes.NewBuffer(json_str) 
	str_read_closer := io.NopCloser(str_read)

	req.Body = str_read_closer

	resp, err_resp := client.Do(req)
	if err_resp != nil {
		log.Fatalf("%s\n", err_resp.Error())
	}
	
	// type of data is []byte
	data, err_data := io.ReadAll(resp.Body)
	if err_data != nil {
		log.Fatalf("%s\n", err_data.Error())
	}

	data_as_string := string(data)
	fmt.Printf("%s\n", data_as_string)	
}

