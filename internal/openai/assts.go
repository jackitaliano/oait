package openai

import (
	"fmt"
)

func retrieveAsst(c chan AsstObject, key string, fileID string, orgID string) {

	asstObject, err := GetAsstObject(key, fileID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- AsstObject{}
		return
	}

	c <- *asstObject
}

func RetrieveAssts(key string, threadIDs []string, orgID string) *[]AsstObject {
	c := make(chan AsstObject, len(threadIDs))

	for _, threadID := range threadIDs {
		go retrieveAsst(c, key, threadID, orgID)
	}

	files := make([]AsstObject, len(threadIDs))
	for i := range files {
		files[i] = <-c
	}

	return &files
}

func RetrieveAllAssts(key string, orgID string) (*[]AsstObject, error) {
	files, err := GetAllAsstObjects(key, orgID)

	if err != nil {
		return nil, err
	}

	return &files.Data, nil
}

func deleteAsst(c chan *AsstDeleteResponse, key string, asstID string, orgID string) {

	deleteResponse, err := DeleteAsst(key, asstID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	c <- deleteResponse
}

func DeleteAssts(key string, fileIDs []string, orgID string) int {
	c := make(chan *AsstDeleteResponse, len(fileIDs))

	for _, threadID := range fileIDs {
		go deleteAsst(c, key, threadID, orgID)
	}

	results := make([]*AsstDeleteResponse, len(fileIDs))
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

func CreateAssistant(key string, createdAsst *CreatedAssistant, orgID string) (*AsstObject, error) {
	asst, err := NewAssistant(key, createdAsst, orgID)

	return asst, err
}
