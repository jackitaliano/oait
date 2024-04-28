package openai

import (
	"fmt"
)

func deleteFile(c chan *FileDeleteResponse, key string, fileID string, orgID string) {

	deleteResponse, err := DeleteFile(key, fileID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- nil
		return
	}

	c <- deleteResponse
}

func DeleteFiles(key string, fileIDs []string, orgID string) int {
	c := make(chan *FileDeleteResponse, len(fileIDs))

	for _, threadID := range fileIDs {
		go deleteFile(c, key, threadID, orgID)
	}

	results := make([]*FileDeleteResponse, len(fileIDs))
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

func retrieveFile(c chan FileObject, key string, fileID string, orgID string) {

	fileObject, err := GetFileObject(key, fileID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- FileObject{}
		return
	}

	c <- *fileObject
}

func RetrieveFiles(key string, threadIDs []string, orgID string) *[]FileObject {
	c := make(chan FileObject, len(threadIDs))

	for _, threadID := range threadIDs {
		go retrieveFile(c, key, threadID, orgID)
	}

	files := make([]FileObject, len(threadIDs))
	for i := range files {
		files[i] = <-c
	}

	return &files
}

func RetrieveAllFiles(key string, orgID string) *[]FileObject {
	files, err := GetAllFileObjects(key, orgID)

	if err != nil {
		return nil
	}

	return &files.Data
}
