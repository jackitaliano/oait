package assts

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func retrieveAsst(c chan openai.AsstObject, key string, fileID string, orgID string) {

	asstObject, err := openai.GetAsstObject(key, fileID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- openai.AsstObject{}
		return
	}

	c <- *asstObject
}

func RetrieveAssts(key string, threadIDs []string, orgID string) *[]openai.AsstObject {
	c := make(chan openai.AsstObject, len(threadIDs))

	for _, threadID := range threadIDs {
		go retrieveAsst(c, key, threadID, orgID)
	}

	files := make([]openai.AsstObject, len(threadIDs))
	for i := range files {
		files[i] = <-c
	}

	return &files
}

func RetrieveAllAssts(key string, orgID string) (*[]openai.AsstObject, error) {
	files, err := openai.GetAllAsstObjects(key, orgID)

	if err != nil {
		return nil, err
	}

	return &files.Data, nil
}
