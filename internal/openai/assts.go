package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jackitaliano/oait/internal/request"
)

type Property struct {
	Type        string `json:"string"`
	Description string `json:"description"`
}

type FunctionParameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Function struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  FunctionParameters `json:"parameters"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function,omitempty"`
}

type AsstObject struct {
	ID            string                         `json:"id"`
	Object        string                         `json:"object"`
	Created       int64                          `json:"created_at"`
	Name          string                         `json:"name"`
	Description   string                         `json:"description"`
	Instructions  string                         `json:"instructions"`
	Model         string                         `json:"model"`
	Tools         []Tool                         `json:"tools"`
	ToolResources map[string]map[string][]string `json:"tool_resources,omitempty"`
	ResFormat     string                         `json:"response_format"`
	Temp          float64                        `json:"temperature"`
	TopP          float64                        `json:"top_p"`
}

func (a AsstObject) CreatedAt() int64 {
	return a.Created
}

type AsstObjectsResponse struct {
	Data   []AsstObject `json:"data"`
	Object string       `json:"object"`
}

func (a AsstObjectsResponse) Len() int {
	return len(a.Data)
}

type CreatedAssistant struct {
	Name          string              `json:"name,omitempty"`
	Description   string              `json:"description,omitempty"`
	Instructions  string              `json:"instructions,omitempty"`
	Model         string              `json:"model"`
	Tools         []Tool              `json:"tools,omitempty"`
	ToolResources map[string][]string `json:"tool_resources,omitempty"`
	ResFormat     string              `json:"response_format,omitempty"`
	Temp          float64             `json:"temperature,omitempty"`
	TopP          float64             `json:"top_p,omitempty"`
}

type AsstDeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

func GetAsstObject(key string, asstID string, orgID string) (*AsstObject, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/assistants/%v", asstID)

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
	req.Header.Set("Openai-Beta", "assistants=v2")

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[AsstObject](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func GetAllAsstObjects(key string, orgID string) (*AsstObjectsResponse, error) {
	var url string

	url = "https://api.openai.com/v1/assistants?limit=100"

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
	req.Header.Set("Openai-Beta", "assistants=v2")

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[AsstObjectsResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func CreateAssistant(key string, asst *CreatedAssistant, orgID string) (*AsstObject, error) {
	var url string

	url = "https://api.openai.com/v1/assistants"

	method := "POST"
	jsonData, err := json.Marshal(*asst)
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
	req.Header.Set("Openai-Beta", "assistants=v2")

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[AsstObject](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func DeleteAsst(key string, asstID string, orgID string) (*AsstDeleteResponse, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/assistants/%v", asstID)

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
	req.Header.Set("Openai-Beta", "assistants=v2")

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[AsstDeleteResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}
