package filter

import (
	"errors"
	"time"
)

type CreatedAtProvider interface {
	CreatedAt() int64
}

type LenProvider interface {
	Len() int
}

func DaysLTE[T CreatedAtProvider](list *[]T, days float64) (*[]T, error) {

	if days < 0 {
		err := errors.New("Invalid number of days: negative numbers not supported")
		return nil, err
	}

	filtered := []T{}

	const dayInSeconds float64 = 86400

	unixTime := time.Now().Unix()
	daysInSeconds := days * dayInSeconds

	unixTimeXDaysAgo := unixTime - int64(daysInSeconds)

	for _, obj := range *list {
		if unixTimeXDaysAgo <= obj.CreatedAt() {
			filtered = append(filtered, obj)
		}
	}

	return &filtered, nil
}

func DaysGT[T CreatedAtProvider](list *[]T, days float64) (*[]T, error) {
	if days < 0 {
		err := errors.New("Invalid number of days: negative numbers not supported")
		return nil, err
	}

	filtered := []T{}

	const dayInSeconds float64 = 86400

	unixTime := time.Now().Unix()
	daysInSeconds := days * dayInSeconds

	unixTimeXDaysAgo := unixTime - int64(daysInSeconds)

	for _, obj := range *list {
		if unixTimeXDaysAgo > obj.CreatedAt() {
			filtered = append(filtered, obj)
		}
	}

	return &filtered, nil
}

func LengthLTE[T LenProvider](list *[]T, length float64) (*[]T, error) {
	if length < 0 {
		err := errors.New("Invalid length: negative numbers not supported")
		return nil, err
	}

	filtered := []T{}

	for _, obj := range *list {
		if float64(obj.Len()) <= length {
			filtered = append(filtered, obj)
		}
	}

	return &filtered, nil
}

func LengthGT[T LenProvider](list *[]T, length float64) (*[]T, error) {
	if length < 0 {
		err := errors.New("Invalid length: negative numbers not supported")
		return nil, err
	}

	filtered := []T{}

	for _, obj := range *list {
		if float64(obj.Len()) > length {
			filtered = append(filtered, obj)
		}
	}

	return &filtered, nil
}
