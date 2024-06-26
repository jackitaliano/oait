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

type MessagesResponse struct {
	Object string    `json:"object"`
	Data   []Message `json:"data"`
}

type SessionThreadsResponse struct {
	Object  string   `json:"object"`
	Data    []Thread `json:"data"`
	FirstID string   `json:"first_id"`
	LastID  string   `json:"last_id"`
	HasMore bool     `json:"has_more"`
}

type ThreadDeleteResponse struct {
	Object  string `json:"object"`
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type Thread struct {
	Object    string            `json:"object"`
	ID        string            `json:"id"`
	CreatedAt int               `json:"created_at"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type Messages struct {
	Messages []Message
}

type Message struct {
	ID          string           `json:"id"`
	Object      string           `json:"object"`
	CreatedAt   int64            `json:"created_at"`
	AssistantID string           `json:"assistant_id,omitempty"`
	ThreadID    string           `json:"thread_id"`
	RunID       string           `json:"run_id,omitempty"`
	Role        string           `json:"role"`
	Content     []MessageContent `json:"content"`
	Attachments []Attachment     `json:"attachments,omitempty"`
}

type Attachment struct {
	FileId string `json:"file_id"`
	Tools  []Tool `json:"tools"` // defined in asstReq.go
}

type MessageContent struct {
	Type      string       `json:"type"`
	Text      *MessageText `json:"text,omitempty"`
	ImageFile *ImageFile   `json:"image_file,omitempty"`
}

type MessageText struct {
	Value       string       `json:"value"`
	Annotations []Annotation `json:"annotations,omitempty"`
}

type ImageFile struct {
	FileId string `json:"file_id,omitempty"`
}

type CreatedMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Annotation struct {
	Type       string   `json:"type"`
	FilePath   FilePath `json:"file_path"`
	Text       string   `json:"text"`
	StartIndex int      `json:"start_index"`
}

type FilePath struct {
	FileID string `json:"file_id"`
}

func (m Messages) GetCreatedAt() int64 {
	if m.GetLen() > 0 {
		return m.Messages[0].CreatedAt
	}
	return 0
}

func (m Messages) GetLen() int {
	return len(m.Messages)
}

func (m Messages) GetContent() []string {

	content := make([]string, len(m.Messages))

	for i, msg := range m.Messages {
		for _, c := range msg.Content {
			if c.Type == "text" {
				content[i] = c.Text.Value
				continue
			} else {
				content[i] = ""
			}
		}
	}

	return content
}

func (t Thread) GetMetadata() map[string]string {
	return t.Metadata;
}

func GetThreadMessages(key string, threadID string, orgID string) (*MessagesResponse, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%v/messages?limit=100", threadID)
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

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[MessagesResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func GetSessionThreads(sessionID string, orgID string) (*SessionThreadsResponse, error) {
	url := "https://api.openai.com/v1/threads?limit=100"
	method := "GET"
	var reqBody io.Reader = nil

	req, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		errMsg := fmt.Sprintf("Error creating request to '%v':\nError: %v", url, err)
		err = errors.New(errMsg)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+sessionID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Openai-Beta", "assistants=v1")

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[SessionThreadsResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func GetThread(key string, threadId string, orgID string) (*Thread, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%v", threadId);
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
	req.Header.Set("Openai-Beta", "assistants=v1")

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[Thread](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func DeleteThread(key string, threadID string, orgID string) (*ThreadDeleteResponse, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%v", threadID)
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

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[ThreadDeleteResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func PostMessage(key string, threadID string, message *CreatedMessage, orgID string) (*Message, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%v/messages", threadID)
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

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[Message](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}
