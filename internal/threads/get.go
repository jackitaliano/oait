package threads

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func retrieveThread(c chan *openai.Messages, key string, threadID string, orgID string) {

	messageResponse, err := openai.GetThreadMessages(key, threadID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- &openai.Messages{}
		return
	}

	messageData := &openai.Messages{Messages: (*messageResponse).Data}
	c <- messageData
}

func RetrieveThreads(key string, threadIDs []string, orgID string) *[]openai.Messages {
	c := make(chan *openai.Messages, len(threadIDs))

	for _, threadID := range threadIDs {
		go retrieveThread(c, key, threadID, orgID)
	}

	results := make([]*openai.Messages, len(threadIDs))
	for i := range results {
		results[i] = <-c
	}

	threads := []openai.Messages{}
	for _, thread := range results {
		if thread != nil {
			threads = append(threads, *thread)
		}
	}

	return &threads
}
