package threads

import (
	"encoding/json"
	"errors"

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

func ParseThreads(threadIds []string, threads *[][]openai.Message) *[]Thread {
	c := make(chan Thread)

	for i, thread := range *threads {
		threadId := threadIds[i]
		go parseThread(c, threadId, thread)
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

func parseThread(c chan Thread, threadId string, thread []openai.Message) {
	messages := []Message{}

	if len(thread) < 1 {
		c <- Thread{threadId, messages}
		return
	}

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

func reverse[T any](list []T) []T {
	for i, j := 0, len(list)-1; i < j; {
		list[i], list[j] = list[j], list[i]
		i++
		j--
	}
	return list
}
