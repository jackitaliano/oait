package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jackitaliano/oait-go/internal/request"
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

type Thread struct {
	Object    string `json:"object"`
	Id        string `json:"id"`
	CreatedAt int    `json:"created_at"`
}

type SessionThreadsResponse struct {
	Object  string   `json:"object"`
	Data    []Thread `json:"data"`
	FirstId string   `json:"first_id"`
	LastId  string   `json:"last_id"`
	HasMore bool     `json:"has_more"`
}

type ThreadDeleteResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type CreatedMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func GetThreadMessages(key string, threadId string, orgId string) (*MessagesResponse, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%v/messages?limit=100", threadId)
	method := "GET"
	var reqBody io.Reader = nil

	req, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		errMsg := fmt.Sprintf("Error creating request to '%v':\nError: %v", url, err)
		err = errors.New(errMsg)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	if orgId != "" {
		req.Header.Set("Openai-Organization", orgId)
	}

	resBody, err := request.Process[MessagesResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func GetSessionThreads(sessionId string, orgId string) (*SessionThreadsResponse, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads?limit=100")
	method := "GET"
	var reqBody io.Reader = nil

	req, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		errMsg := fmt.Sprintf("Error creating request to '%v':\nError: %v", url, err)
		err = errors.New(errMsg)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+sessionId)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Openai-Beta", "assistants=v1")

	if orgId != "" {
		req.Header.Set("Openai-Organization", orgId)
	}

	resBody, err := request.Process[SessionThreadsResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func DeleteThread(key string, threadId string, orgId string) (*ThreadDeleteResponse, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%v", threadId)
	method := "DELETE"
	var reqBody io.Reader = nil

	req, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		errMsg := fmt.Sprintf("Error creating request to '%v':\nError: %v", url, err)
		err = errors.New(errMsg)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	if orgId != "" {
		req.Header.Set("Openai-Organization", orgId)
	}

	resBody, err := request.Process[ThreadDeleteResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func AddMessage(key string, threadId string, message *CreatedMessage, orgId string) (*Message, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%v/messages", threadId)
	method := "POST"

	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewReader(jsonData)

	req, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		errMsg := fmt.Sprintf("Error creating request to '%v':\nError: %v", url, err)
		err = errors.New(errMsg)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	if orgId != "" {
		req.Header.Set("Openai-Organization", orgId)
	}

	resBody, err := request.Process[Message](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}
