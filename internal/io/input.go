package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jackitaliano/oait/internal/openai"
)

func SingleInput(ID string) (string, error) {
	trimmedString := strings.Trim(ID, " ")

	if trimmedString == "" {
		errMsg := "Invalid file id passed ' '"
		err := errors.New(errMsg)
		return "", err
	}

	return trimmedString, nil
}

func ListInput(IDs []string) ([]string, error) {
	var allFileIDs []string

	for _, idsStr := range IDs {
		ids := splitIDs(idsStr, " ")

		allFileIDs = append(allFileIDs, ids...)
	}

	return allFileIDs, nil
}

func FileInput(fileName string) ([]string, error) {
	var data []string
	var err error

	fileExt := fileName[len(fileName)-3:]
	if fileExt == "txt" {
		data, err = txtInput(fileName)

		if err != nil {
			data = []string{}
		}

	} else {
		err = errors.New("Invalid file input type")
	}

	return data, err

}

func JSONInput[T any](fileName string) (*T, error) {
	fileExt := fileName[len(fileName)-4:]
	if fileExt != "json" {
		errMsg := fmt.Sprintf("Invalid file name: '%v' Only JSON is valid.", fileName)
		err := errors.New(errMsg)
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jsonData T
	err = json.NewDecoder(file).Decode(&jsonData)

	if err != nil {
		return nil, err
	}

	return &jsonData, nil
}

func SessionInput(sessionID string, orgID string) ([]string, error) {
	sessionThreadsRes, err := openai.GetSessionThreads(sessionID, orgID)

	if err != nil {
		return nil, err
	}

	threadIds := make([]string, len(sessionThreadsRes.Data))

	for i, thread := range (*sessionThreadsRes).Data {
		id := thread.ID

		threadIds[i] = id
	}

	return threadIds, nil
}

func txtInput(fileName string) ([]string, error) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		err = errors.New("Failed reading file: " + fileName + ". Error: " + err.Error())
		return nil, err
	}

	stringData := string(data)
	splitStrings := splitIDs(stringData, "\n")

	var asstIDs []string
	for _, val := range splitStrings {
		if val != "" {
			asstIDs = append(asstIDs, val)
		}
	}

	return asstIDs, nil
}

func splitIDs(str string, delimeter string) []string {
	trimmedString := strings.Trim(str, " ")
	splitString := strings.Split(trimmedString, delimeter)

	for i, val := range splitString {
		splitString[i] = strings.Trim(val, " ")
	}

	return splitString
}
