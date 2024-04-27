package threads

import (
	"errors"
	"os"
	"strings"

	"github.com/jackitaliano/oait/internal/openai"
)

func SingleInput(threadID string) (string, error) {
	trimmedString := strings.Trim(threadID, " ")

	if trimmedString == "" {
		errMsg := "Invalid thread id passed ' '"
		err := errors.New(errMsg)
		return "", err
	}

	return trimmedString, nil
}

func ListInput(threadIDs []string) ([]string, error) {
	var allThreadIDs []string

	for _, idsStr := range threadIDs {
		ids := splitIDs(idsStr, " ")

		for _, id := range ids {
			allThreadIDs = append(allThreadIDs, id)
		}
	}

	return allThreadIDs, nil
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

func SessionInput(sessionID string, orgID string) ([]string, error) {
	sessionThreadsRes, err := openai.GetSessionThreads(sessionID, orgID)

	if err != nil {
		return nil, err
	}

	threadIDs := make([]string, len(sessionThreadsRes.Data))

	for i, thread := range (*sessionThreadsRes).Data {
		id := thread.ID

		threadIDs[i] = id
	}

	return threadIDs, nil
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

	var threadIDs []string
	for _, val := range splitStrings {
		if val != "" {
			threadIDs = append(threadIDs, val)
		}
	}

	return threadIDs, nil
}

func splitIDs(str string, delimeter string) []string {
	trimmedString := strings.Trim(str, " ")
	splitString := strings.Split(trimmedString, delimeter)

	for i, val := range splitString {
		splitString[i] = strings.Trim(val, " ")
	}

	return splitString
}
