package threads

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func retrieveThread(c chan []openai.Message, key string, threadId string, orgId string) {

	messageResponse, err := openai.GetThreadMessages(key, threadId, orgId)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	messageData := (*messageResponse).Data
	c <- messageData
}

func RetrieveThreads(key string, threadIds []string, orgId string) *[][]openai.Message {
	c := make(chan []openai.Message)

	for _, threadId := range threadIds {
		go retrieveThread(c, key, threadId, orgId)
	}

	results := make([][]openai.Message, len(threadIds))
	for i := range results {
		results[i] = <-c
	}

	threads := [][]openai.Message{}
	for _, thread := range results {
		if thread != nil {
			threads = append(threads, thread)
		}
	}

	return &threads
}
