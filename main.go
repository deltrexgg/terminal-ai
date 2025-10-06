package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

func main() {

    token := 200

	if len(os.Args) < 2 {
		fmt.Println("Error: Please provide the question")
		return
	}

	question := os.Args[1]

	if len(os.Args) > 2 && os.Args[2] != "" {
		tokenStr := os.Args[2]
		var err error
		token, err = strconv.Atoi(tokenStr)
		if err != nil {
			fmt.Println("Error: token must be a number")
			return
		}
	}

    
	pcEnv := os.Getenv("TAIURL")
	if pcEnv == "" {
        fmt.Println("Error : Please provide the AI url in env as TAIURL")
	}

	body := map[string]interface{}{
		"model": "Qwen2.5-0.5B-Instruct-Q6_K",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert code debugging assistant. Always explain errors, suggest fixes, and provide corrected code.",
			},
			{
				"role":    "user",
				"content": question,
			},
		},
		"max_tokens":  token,
		"temperature": 0,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

    fmt.Println("AI is thinking .. wait!")
	resp, err := http.Post(pcEnv+"/v1/chat/completions", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var apiResp Response
	err = json.Unmarshal(respBytes, &apiResp)
	if err != nil {
		panic(err)
	}

	if len(apiResp.Choices) > 0 {
		fmt.Println(apiResp.Choices[0].Message.Content)
	} else {
		fmt.Println("No response from assistant.")
	}
}
