package threads

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Thread struct {
	threadId string
	messages []map[string]string
}

func ListInput(threadIds *[]string) (*[]*Thread, error) {
	var threads []*Thread

	fmt.Println(threadIds)

	for _, threadId := range *threadIds {
		messages := []map[string]string{{"role": "user", "text": "test text"}}

		thread := &Thread{threadId, messages}

		threads = append(threads, thread)
	}

	return &threads, nil
}

func FileInput(fileName string) ( *[]*Thread, error ) {
	var threads = []*Thread{}

	messages := []map[string]string{{"role": "user", "text": "test text"}}
	thread := &Thread{fileName, messages}
	threads = append(threads, thread)

	return &threads, nil

}

func ThreadsToJson(threads *[]*Thread) (*[]byte, error) {
	var threadsExpanded = []Thread{}

	for _, thread := range *threads {
		threadsExpanded = append(threadsExpanded, *thread)
	}

	fmt.Println(threadsExpanded)

	b, err := json.Marshal(threadsExpanded)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func FileOutput(fileName string, data *[]byte) error {
	fmt.Printf("Outputting messages to: %v\n", fileName)

	err := os.WriteFile(fileName, *data, 0644)

	if err != nil {
		err = errors.New("Failed to write to file: " + fileName)
		return err
	}

	return nil
}
