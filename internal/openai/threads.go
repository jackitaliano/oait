package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type MessageText struct {
	Value       string              `json:"value"`
	Annotations []map[string]string `json:"annotations"`
}

type MessageContent struct {
	Type string      `json:"type"`
	Text MessageText `json:"text"`
}

type Message struct {
	Id          string           `json:"id"`
	Object      string           `json:"object"`
	CreatedAt   int              `json:"created_at"`
	AssistantId string           `json:"assistant_id"`
	ThreadId    string           `json:"thread_id"`
	RunId       string           `json:"run_id"`
	Role        string           `json:"role"`
	Content     []MessageContent `json:"content"`
}

type MessagesResponse struct {
	Object string    `json:"object"`
	Data   []Message `json:"data"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

func GetThreadMessages(key string, threadId string) (MessagesResponse, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.openai.com/v1/threads/%v/messages?limit=100", threadId)
	method := "GET"
	var reqBody io.Reader = nil

	req, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		errMsg := fmt.Sprintf("Error creating request to '%v':\nError: %v", url, err)
		err = errors.New(errMsg)
		return MessagesResponse{"", []Message{}}, err
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	res, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Error making request to '%v':\nError: %v", url, err)
		err = errors.New(errMsg)
		return MessagesResponse{"", []Message{}}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		errMsg := fmt.Sprintf("Error: Request (%v). Status: (%v). Message: %v", url, res.StatusCode, errRes.Error.Message)
		err = errors.New(errMsg)
		return MessagesResponse{"", []Message{}}, err
	}

	var resBody MessagesResponse
	err = json.NewDecoder(res.Body).Decode(&resBody)

	if err != nil {
		errMsg := fmt.Sprintf("Error reading response body from '%v':\n%v\n", url, err)
		err = errors.New(errMsg)
		return MessagesResponse{"", []Message{}}, err
	}

	return resBody, nil
}
