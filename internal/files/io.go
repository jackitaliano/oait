package files

import (
	"errors"
	"os"
	"strings"
)

func SingleInput(fileID string) (string, error) {
	trimmedString := strings.Trim(fileID, " ")

	if trimmedString == "" {
		errMsg := "Invalid file id passed ' '"
		err := errors.New(errMsg)
		return "", err
	}

	return trimmedString, nil
}

func ListInput(fileIDs []string) ([]string, error) {
	var allFileIDs []string

	for _, idsStr := range fileIDs {
		ids := splitIDs(idsStr, " ")

		for _, id := range ids {
			allFileIDs = append(allFileIDs, id)
		}
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
	splitStrings := splitIDs(stringData, "\n")

	var fileIDs []string
	for _, val := range splitStrings {
		if val != "" {
			fileIDs = append(fileIDs, val)
		}
	}

	return fileIDs, nil
}

func splitIDs(str string, delimeter string) []string {
	trimmedString := strings.Trim(str, " ")
	splitString := strings.Split(trimmedString, delimeter)

	for i, val := range splitString {
		splitString[i] = strings.Trim(val, " ")
	}

	return splitString
}
