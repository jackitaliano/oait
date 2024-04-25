package threads

import (
	"fmt"

	"github.com/jackitaliano/oait-go/internal/openai"
)

func deleteThread(c chan *openai.ThreadDeleteResponse, key string, threadId string, orgId string) {

	deleteResponse, err := openai.DeleteThread(key, threadId, orgId)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	c <- deleteResponse
}

func DeleteThreads(key string, threadIds []string, orgId string) int {
	c := make(chan *openai.ThreadDeleteResponse)

	for _, threadId := range threadIds {
		go deleteThread(c, key, threadId, orgId)
	}

	results := make([]*openai.ThreadDeleteResponse, len(threadIds))
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
