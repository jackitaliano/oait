package files

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func deleteFile(c chan *openai.FileDeleteResponse, key string, fileId string, orgId string) {

	deleteResponse, err := openai.DeleteFile(key, fileId, orgId)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	c <- deleteResponse
}

func DeleteFiles(key string, fileIds []string, orgId string) int {
	c := make(chan *openai.FileDeleteResponse, len(fileIds))

	for _, threadId := range fileIds {
		go deleteFile(c, key, threadId, orgId)
	}

	results := make([]*openai.FileDeleteResponse, len(fileIds))
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
