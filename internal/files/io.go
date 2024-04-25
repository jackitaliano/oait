package files

import (
	"errors"
	"os"
	"strings"
)

func SingleInput(fileId string) (string, error) {
	trimmedString := strings.Trim(fileId, " ")

	if trimmedString == "" {
		errMsg := "Invalid file id passed ' '"
		err := errors.New(errMsg)
		return "", err
	}

	return trimmedString, nil
}

func ListInput(fileIds []string) ([]string, error) {
	var allFileIds []string

	for _, idsStr := range fileIds {
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

	var fileIds []string
	for _, val := range splitStrings {
		if val != "" {
			fileIds = append(fileIds, val)
		}
	}

	return fileIds, nil
}

func splitIds(str string, delimeter string) []string {
	trimmedString := strings.Trim(str, " ")
	splitString := strings.Split(trimmedString, delimeter)

	for i, val := range splitString {
		splitString[i] = strings.Trim(val, " ")
	}

	return splitString
}
