package openai

import (
	"fmt"
)

func AddMessage(key string, threadID string, createdMessage *CreatedMessage, orgID string) (*Message, error) {
	message, err := PostMessage(key, threadID, createdMessage, orgID)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func deleteThread(c chan *ThreadDeleteResponse, key string, threadID string, orgID string) {

	deleteResponse, err := DeleteThread(key, threadID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	c <- deleteResponse
}

func DeleteThreads(key string, threadIDs []string, orgID string) int {
	c := make(chan *ThreadDeleteResponse, len(threadIDs))

	for _, threadID := range threadIDs {
		go deleteThread(c, key, threadID, orgID)
	}

	results := make([]*ThreadDeleteResponse, len(threadIDs))
	numDeleted := 0

	for i := range results {
		res := <-c
		results[i] = res

		if res.Deleted {
			numDeleted += 1
		}
	}

	return numDeleted
}

func retrieveThread(c chan *Messages, key string, threadID string, orgID string) {

	messageResponse, err := GetThreadMessages(key, threadID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- &Messages{}
		return
	}

	messageData := &Messages{Messages: (*messageResponse).Data}
	c <- messageData
}

func RetrieveThreads(key string, threadIDs []string, orgID string) *[]Messages {
	c := make(chan *Messages, len(threadIDs))

	for _, threadID := range threadIDs {
		go retrieveThread(c, key, threadID, orgID)
	}

	results := make([]*Messages, len(threadIDs))
	for i := range results {
		results[i] = <-c
	}

	threads := []Messages{}
	for _, thread := range results {
		if thread != nil {
			threads = append(threads, *thread)
		}
	}

	return &threads
}
