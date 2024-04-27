package threads

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func deleteThread(c chan *openai.ThreadDeleteResponse, key string, threadID string, orgID string) {

	deleteResponse, err := openai.DeleteThread(key, threadID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	c <- deleteResponse
}

func DeleteThreads(key string, threadIDs []string, orgID string) int {
	c := make(chan *openai.ThreadDeleteResponse, len(threadIDs))

	for _, threadID := range threadIDs {
		go deleteThread(c, key, threadID, orgID)
	}

	results := make([]*openai.ThreadDeleteResponse, len(threadIDs))
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
