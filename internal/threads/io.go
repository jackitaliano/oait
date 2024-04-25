package threads

import (
	"errors"
	"os"
	"strings"

	"github.com/jackitaliano/oait-go/internal/openai"
)

func SingleInput(threadId string) (string, error) {
	trimmedString := strings.Trim(threadId, " ")

	if trimmedString == "" {
		errMsg := "Invalid thread id passed ' '"
		err := errors.New(errMsg)
		return "", err
	}

	return trimmedString, nil
}

func ListInput(threadIds []string) ([]string, error) {
	var allThreadIds []string

	for _, idsStr := range threadIds {
		ids := splitIds(idsStr, " ")

		for _, id := range ids {
			allThreadIds = append(allThreadIds, id)
		}
	}

	return allThreadIds, nil
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

func SessionInput(sessionId string, orgId string) ([]string, error) {
	sessionThreadsRes, err := openai.GetSessionThreads(sessionId, orgId)

	if err != nil {
		return nil, err
	}

	threadIds := make([]string, len(sessionThreadsRes.Data))

	for i, thread := range (*sessionThreadsRes).Data {
		id := thread.Id

		threadIds[i] = id
	}

	return threadIds, nil
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

	var threadIds []string
	for _, val := range splitStrings {
		if val != "" {
			threadIds = append(threadIds, val)
		}
	}

	return threadIds, nil
}

func splitIds(str string, delimeter string) []string {
	trimmedString := strings.Trim(str, " ")
	splitString := strings.Split(trimmedString, delimeter)

	for i, val := range splitString {
		splitString[i] = strings.Trim(val, " ")
	}

	return splitString
}

