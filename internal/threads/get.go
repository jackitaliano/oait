package threads

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func retrieveThread(c chan *[]openai.Message, key string, threadID string, orgID string) {

	messageResponse, err := openai.GetThreadMessages(key, threadID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- &[]openai.Message{}
		return
	}

	messageData := &(*messageResponse).Data
	c <- messageData
}

func RetrieveThreads(key string, threadIDs []string, orgID string) *[][]openai.Message {
	c := make(chan *[]openai.Message, len(threadIDs))

	for _, threadID := range threadIDs {
		go retrieveThread(c, key, threadID, orgID)
	}

	results := make([]*[]openai.Message, len(threadIDs))
	for i := range results {
		results[i] = <-c
	}

	threads := [][]openai.Message{}
	for _, thread := range results {
		if thread != nil {
			threads = append(threads, *thread)
		}
	}

	return &threads
}
