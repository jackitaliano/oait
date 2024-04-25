package files

import (
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
)

func retrieveFile(c chan openai.FileObject, key string, fileId string, orgId string) {

	fileObject, err := openai.GetFileObject(key, fileId, orgId)

	if err != nil {
		fmt.Println(err)
		c <- openai.FileObject{}
		return
	}

	c <- *fileObject
}

func RetrieveFiles(key string, threadIds []string, orgId string) *[]openai.FileObject {
	c := make(chan openai.FileObject)

	for _, threadId := range threadIds {
		go retrieveFile(c, key, threadId, orgId)
	}

	files := make([]openai.FileObject, len(threadIds))
	for i := range files {
		files[i] = <-c
	}

	return &files
}

func RetrieveAllFiles(key string, orgId string) *[]openai.FileObject {
	files, err := openai.GetAllFileObjects(key, orgId)

	if err != nil {
		return nil
	}

	return &files.Data
}
