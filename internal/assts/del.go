package assts

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func deleteAsst(c chan *openai.AsstDeleteResponse, key string, asstID string, orgID string) {

	deleteResponse, err := openai.DeleteAsst(key, asstID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	c <- deleteResponse
}

func DeleteAssts(key string, fileIDs []string, orgID string) int {
	c := make(chan *openai.AsstDeleteResponse, len(fileIDs))

	for _, threadID := range fileIDs {
		go deleteAsst(c, key, threadID, orgID)
	}

	results := make([]*openai.AsstDeleteResponse, len(fileIDs))
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
