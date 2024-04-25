package openai

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jackitaliano/oait-go/internal/request"
)

type FileObject struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int    `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

type FileObjectsResponse struct {
	Data   []FileObject `json:"data"`
	Object string       `json:"object"`
}

type FileDeleteResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

func GetFileObject(key string, fileId string, orgId string) (*FileObject, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/files/%v", fileId)

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

	if orgId != "" {
		req.Header.Set("Openai-Organization", orgId)
	}

	resBody, err := request.Process[FileObject](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func GetAllFileObjects(key string, orgId string) (*FileObjectsResponse, error) {
	var url string

	url = "https://api.openai.com/v1/files"

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

	if orgId != "" {
		req.Header.Set("Openai-Organization", orgId)
	}

	resBody, err := request.Process[FileObjectsResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func DeleteFile(key string, fileId string, orgId string) (*FileDeleteResponse, error) {
	url := fmt.Sprintf("https://api.openai.com/v1/files/%v", fileId)

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

	if orgId != "" {
		req.Header.Set("Openai-Organization", orgId)
	}

	resBody, err := request.Process[FileDeleteResponse](req)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}
