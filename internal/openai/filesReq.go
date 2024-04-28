package openai

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jackitaliano/oait/internal/request"
)

type FileObjectsResponse struct {
	Data   []FileObject `json:"data"`
	Object string       `json:"object"`
}

type FileDeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type FileObject struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Bytes    int    `json:"bytes"`
	Created  int64  `json:"created_at"`
	Filename string `json:"filename"`
	Purpose  string `json:"purpose"`
}

func (f FileObject) CreatedAt() int64 {
	return f.Created
}

func GetFileObject(key string, fileID string, orgID string) (*FileObject, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/files/%v", fileID)

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

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[FileObject](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func GetAllFileObjects(key string, orgID string) (*FileObjectsResponse, error) {
	url := "https://api.openai.com/v1/files"

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

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[FileObjectsResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func DeleteFile(key string, fileID string, orgID string) (*FileDeleteResponse, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/files/%v", fileID)

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

	if orgID != "" {
		req.Header.Set("Openai-Organization", orgID)
	}

	resBody, err := request.Process[FileDeleteResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}
