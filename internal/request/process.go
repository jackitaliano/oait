package request

import (
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
)

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}


type ErrorResponse struct {
	Error Error `json:"error"`
}

type Response interface {}

func Process[T Response](req *http.Request) (*T, error) {
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Error making request to '%v':\nError: %v", *req.URL, err)
		err = errors.New(errMsg)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		errMsg := fmt.Sprintf("Error: Request (%v). Status: (%v). Message: %v", *req.URL, res.StatusCode, errRes.Error.Message)
		err = errors.New(errMsg)

		return nil, err
	}

	var resBody T
	err = json.NewDecoder(res.Body).Decode(&resBody)

	if err != nil {
		errMsg := fmt.Sprintf("Error reading response body from '%v':\n%v\n", *req.URL, err)
		err = errors.New(errMsg)
		return nil, err
	}

	return &resBody, nil
}
