package files

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func retrieveFile(c chan openai.FileObject, key string, fileID string, orgID string) {

	fileObject, err := openai.GetFileObject(key, fileID, orgID)

	if err != nil {
		fmt.Println(err)
		c <- openai.FileObject{}
		return
	}

	c <- *fileObject
}

func RetrieveFiles(key string, threadIDs []string, orgID string) *[]openai.FileObject {
	c := make(chan openai.FileObject, len(threadIDs))

	for _, threadID := range threadIDs {
		go retrieveFile(c, key, threadID, orgID)
	}

	files := make([]openai.FileObject, len(threadIDs))
	for i := range files {
		files[i] = <-c
	}

	return &files
}

func RetrieveAllFiles(key string, orgID string) *[]openai.FileObject {
	files, err := openai.GetAllFileObjects(key, orgID)

	if err != nil {
		return nil
	}

	return &files.Data
}
