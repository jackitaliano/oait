package threads

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jackitaliano/oait-go/internal/openai"
)

type Thread struct {
	ThreadId string `json:"thread_id"`
	Messages []map[string]string `json:"messages"`
}

func ListInput(threadIds *[]string) (*[]string, error) {
	var allThreadIds []string

	for _, idsStr := range *threadIds {
		ids := splitIds(idsStr, " ")

		for _, id := range *ids {
			allThreadIds = append(allThreadIds, id)
		}
	}

	return &allThreadIds, nil
}

func txtInput(fileName string) ( *[]string, error ) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		err = errors.New("Failed reading file: " + fileName + ". Error: " + err.Error())
		return nil, err
	}

	stringData := string(data)
	threadIds := splitIds(stringData, "\n");

	return threadIds, nil
}

func splitIds(str string, delimeter string) (*[]string) {
	splitString := strings.Split(str, delimeter)

	for i, val := range splitString {
		splitString[i] = strings.Trim(val, " ")
	}

	return &splitString
}

func FileInput(fileName string) ( *[]string, error ) {
	var data *[]string
	var err error

	fileExt := fileName[len(fileName)-3:]
	fmt.Printf("ext: %v\n", fileExt)
	if fileExt == "txt" {
		data, err = txtInput(fileName)

		if err != nil {
			data = &[]string{}
		}

	} else {
		err = errors.New("Invalid file input type")
	}

	return data, err

}

func retrieveThread(c chan []openai.Message, key string, threadId string) {
	messageResponse, err := openai.GetThreadMessages(key, threadId)

	if (err != nil) {
		fmt.Println(err)
	}

	messageData := messageResponse.Data
	c <- messageData
}

func RetrieveThreads(key string, threadIds *[]string) ( *[][]openai.Message ) {
	c := make(chan []openai.Message)

	for _, threadId := range *threadIds {
		go retrieveThread(c, key, threadId)
	}

	results := make([][]openai.Message, len(*threadIds))
	for i := range results {
		results[i] = <-c
	}

	return &results
}

func ThreadsToJson(threads *[]*Thread) (*[]byte, error) {
	var threadsExpanded = []Thread{}

	for _, thread := range *threads {
		threadsExpanded = append(threadsExpanded, *thread)
	}

	b, err := json.MarshalIndent(threadsExpanded, "", "	")

	if err != nil {
		err = errors.New("JSON Marshal failed with error: " + err.Error())
		return nil, err
	}

	return &b, nil
}

func FileOutput(fileName string, data *[]byte) error {
	err := os.WriteFile(fileName, *data, 0644)

	if err != nil {
		err = errors.New("Failed to write to file: " + fileName)
		return err
	}

	return nil
}
