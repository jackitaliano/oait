package files

import (
	"errors"
	"time"

	"github.com/jackitaliano/oait-go/internal/openai"
)

func FilterByDaysLTE(files *[]openai.FileObject, days float64) (*[]openai.FileObject, error) {

	if (days < 0) {
		err := errors.New("Invalid number of days: negative numbers not supported")
		return nil, err
	}

	filteredFiles := []openai.FileObject{}

	const dayInSeconds float64 = 86400

	unixTime := time.Now().Unix()
	daysInSeconds := days * dayInSeconds

	unixTimeXDaysAgo := unixTime - int64(daysInSeconds)

	for _, file := range *files {
		fileCreatedTime := int64(file.CreatedAt)

		if unixTimeXDaysAgo <= fileCreatedTime {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return &filteredFiles, nil
}

func FilterByDaysGT(files *[]openai.FileObject, days float64) (*[]openai.FileObject, error) {
	if (days < 0) {
		err := errors.New("Invalid number of days: negative numbers not supported")
		return nil, err
	}

	filteredFiles := []openai.FileObject{}

	const dayInSeconds float64 = 86400

	unixTime := time.Now().Unix()
	daysInSeconds := days * dayInSeconds

	unixTimeXDaysAgo := unixTime - int64(daysInSeconds)

	for _, file := range *files {
		fileCreatedTime := int64(file.CreatedAt)

		if unixTimeXDaysAgo > fileCreatedTime {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return &filteredFiles, nil
}
