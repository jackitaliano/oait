package assts

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func SingleInput(asstId string) (string, error) {
	trimmedString := strings.Trim(asstId, " ")

	if trimmedString == "" {
		errMsg := "Invalid file id passed ' '"
		err := errors.New(errMsg)
		return "", err
	}

	return trimmedString, nil
}

func ListInput(asstIds []string) ([]string, error) {
	var allFileIds []string

	for _, idsStr := range asstIds {
		ids := splitIds(idsStr, " ")

		for _, id := range ids {
			allFileIds = append(allFileIds, id)
		}
	}

	return allFileIds, nil
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

func JsonInput[T any](fileName string) (*T, error) {
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

func FileOutput(fileName string, data *[]byte) error {
	err := os.WriteFile(fileName, *data, 0644)

	if err != nil {
		err = errors.New("Failed to write to file: " + fileName)
		return err
	}

	return nil
}

func txtInput(fileName string) ([]string, error) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		err = errors.New("Failed reading file: " + fileName + ". Error: " + err.Error())
		return nil, err
	}

	stringData := string(data)
	splitStrings := splitIds(stringData, "\n")

	var asstIds []string
	for _, val := range splitStrings {
		if val != "" {
			asstIds = append(asstIds, val)
		}
	}

	return asstIds, nil
}

func splitIds(str string, delimeter string) []string {
	trimmedString := strings.Trim(str, " ")
	splitString := strings.Split(trimmedString, delimeter)

	for i, val := range splitString {
		splitString[i] = strings.Trim(val, " ")
	}

	return splitString
}
