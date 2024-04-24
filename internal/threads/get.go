package threads

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackitaliano/oait-go/internal/openai"
)

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type Thread struct {
	ThreadId string    `json:"thread_id,omitempty"`
	Messages []Message `json:"messages,omitempty"`
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

func retrieveThread(c chan []openai.Message, key string, threadId string) {

	messageResponse, err := openai.GetThreadMessages(key, threadId)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	messageData := messageResponse.Data
	c <- messageData
}

func RetrieveThreads(key string, threadIds *[]string) *[][]openai.Message {
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

func reverse[T any](list []T) []T {
	for i, j := 0, len(list)-1; i < j; {
		list[i], list[j] = list[j], list[i]
		i++
		j--
	}
	return list
}

func parseThread(c chan Thread, thread []openai.Message) {
	if len(thread) < 1 {
		c <- Thread{"", []Message{}}
		return
	}

	messages := []Message{}

	var threadId string = thread[0].ThreadId

	for _, msg := range thread {
		for _, content := range msg.Content {
			if content.Type == "text" {
				message := Message{msg.Role, content.Text.Value}
				messages = append(messages, message)
			}
		}
	}
	reversedMessages := reverse(messages)

	parsedThread := Thread{threadId, reversedMessages}

	c <- parsedThread
}

func ParseThreads(threads *[][]openai.Message) *[]Thread {
	c := make(chan Thread)

	for _, thread := range *threads {
		go parseThread(c, thread)
	}

	results := make([]Thread, len(*threads))
	for i := range results {
		results[i] = <-c
	}

	return &results
}

func ThreadsToJson(threads *[]Thread) ([]byte, error) {
	b, err := json.MarshalIndent(threads, "", "\t")

	if err != nil {
		err = errors.New("JSON Marshal failed with error: " + err.Error())
		return nil, err
	}

	return b, nil
}

func FilterByDaysLTE(threads *[][]openai.Message, days float64) *[][]openai.Message {
	filteredThreads := [][]openai.Message{}

	const dayInSeconds float64 = 86400

	for _, thread := range *threads {
		if len(thread) == 0 {
			continue
		}

		mostRecentMessage := thread[0]
		recentTime := float64(mostRecentMessage.CreatedAt)
		var unixTime float64 = float64(time.Now().Unix())

		unixTimeXDaysAgo := unixTime - ((days) * dayInSeconds)

		if unixTimeXDaysAgo <= recentTime {
			filteredThreads = append(filteredThreads, thread)
		}
	}

	return &filteredThreads
}

func FilterByDaysGT(threads *[][]openai.Message, days float64) *[][]openai.Message {
	filteredThreads := [][]openai.Message{}

	const dayInSeconds float64 = 86400

	for _, thread := range *threads {
		if len(thread) == 0 {
			continue
		}

		mostRecentMessage := thread[0]
		recentTime := float64(mostRecentMessage.CreatedAt)
		var unixTime float64 = float64(time.Now().Unix())

		unixTimeXDaysAgo := unixTime - ((days) * dayInSeconds)

		if unixTimeXDaysAgo > recentTime {
			filteredThreads = append(filteredThreads, thread)
		}
	}

	return &filteredThreads
}

func FilterByLengthLTE(threads *[][]openai.Message, length float64) *[][]openai.Message {
	filteredThreads := [][]openai.Message{}

	for _, thread := range *threads {

		if float64(len(thread)) <= length {
			filteredThreads = append(filteredThreads, thread)
		}
	}

	return &filteredThreads
}

func FilterByLengthGT(threads *[][]openai.Message, length float64) *[][]openai.Message {
	filteredThreads := [][]openai.Message{}

	for _, thread := range *threads {
		if float64(len(thread)) > length {
			filteredThreads = append(filteredThreads, thread)
		}
	}

	return &filteredThreads
}

func FileOutput(fileName string, data *[]byte) error {
	err := os.WriteFile(fileName, *data, 0644)

	if err != nil {
		err = errors.New("Failed to write to file: " + fileName)
		return err
	}

	return nil
}
