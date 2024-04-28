package io

import (
	"errors"
	"os"
)

func FileOutput(fileName string, data *[]byte) error {
	err := os.WriteFile(fileName, *data, 0644)

	if err != nil {
		err = errors.New("Failed to write to file: " + fileName)
		return err
	}

	return nil
}
