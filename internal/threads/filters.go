package threads

import (
	"time"

	"github.com/jackitaliano/oait/internal/openai"
)

func FilterByDaysLTE(threads *[][]openai.Message, days float64) *[][]openai.Message {
	filteredThreads := [][]openai.Message{}

	const dayInSeconds float64 = 86400

	for _, thread := range *threads {
		if len(thread) == 0 {
			continue
		}

		mostRecentMessage := thread[0]
		recentTime := float64(mostRecentMessage.CreatedAt)
		var unixTime float64 = float64(time.Now().Unix())

		unixTimeXDaysAgo := unixTime - ((days) * dayInSeconds)

		if unixTimeXDaysAgo <= recentTime {
			filteredThreads = append(filteredThreads, thread)
		}
	}

	return &filteredThreads
}

func FilterByDaysGT(threads *[][]openai.Message, days float64) *[][]openai.Message {
	filteredThreads := [][]openai.Message{}

	const dayInSeconds float64 = 86400

	for _, thread := range *threads {
		if len(thread) == 0 {
			continue
		}

		mostRecentMessage := thread[0]
		recentTime := float64(mostRecentMessage.CreatedAt)
		var unixTime float64 = float64(time.Now().Unix())

		unixTimeXDaysAgo := unixTime - ((days) * dayInSeconds)

		if unixTimeXDaysAgo > recentTime {
			filteredThreads = append(filteredThreads, thread)
		}
	}

	return &filteredThreads
}

func FilterByLengthLTE(threads *[][]openai.Message, length float64) *[][]openai.Message {
	filteredThreads := [][]openai.Message{}

	for _, thread := range *threads {

		if float64(len(thread)) <= length {
			filteredThreads = append(filteredThreads, thread)
		}
	}

	return &filteredThreads
}

func FilterByLengthGT(threads *[][]openai.Message, length float64) *[][]openai.Message {
	filteredThreads := [][]openai.Message{}

	for _, thread := range *threads {
		if float64(len(thread)) > length {
			filteredThreads = append(filteredThreads, thread)
		}
	}

	return &filteredThreads
}
