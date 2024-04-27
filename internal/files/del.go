package files

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func deleteFile(c chan *openai.FileDeleteResponse, key string, fileID string, orgID string) {

	deleteResponse, err := openai.DeleteFile(key, fileID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	c <- deleteResponse
}

func DeleteFiles(key string, fileIDs []string, orgID string) int {
	c := make(chan *openai.FileDeleteResponse, len(fileIDs))

	for _, threadID := range fileIDs {
		go deleteFile(c, key, threadID, orgID)
	}

	results := make([]*openai.FileDeleteResponse, len(fileIDs))
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
