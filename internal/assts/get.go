package assts

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func retrieveAsst(c chan openai.AsstObject, key string, fileId string, orgId string) {

	asstObject, err := openai.GetAsstObject(key, fileId, orgId)

	if err != nil {
		fmt.Println(err)
		c <- openai.AsstObject{}
		return
	}

	c <- *asstObject
}

func RetrieveAssts(key string, threadIds []string, orgId string) *[]openai.AsstObject {
	c := make(chan openai.AsstObject, len(threadIds))

	for _, threadId := range threadIds {
		go retrieveAsst(c, key, threadId, orgId)
	}

	files := make([]openai.AsstObject, len(threadIds))
	for i := range files {
		files[i] = <-c
	}

	return &files
}

func RetrieveAllAssts(key string, orgId string) (*[]openai.AsstObject, error ) {
	files, err := openai.GetAllAsstObjects(key, orgId)

	if err != nil {
		return nil, err
	}

	return &files.Data, nil
}
