package assts

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func deleteAsst(c chan *openai.AsstDeleteResponse, key string, asstId string, orgId string) {

	deleteResponse, err := openai.DeleteAsst(key, asstId, orgId)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	c <- deleteResponse
}

func DeleteAssts(key string, fileIds []string, orgId string) int {
	c := make(chan *openai.AsstDeleteResponse, len(fileIds))

	for _, threadId := range fileIds {
		go deleteAsst(c, key, threadId, orgId)
	}

	results := make([]*openai.AsstDeleteResponse, len(fileIds))
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
