package filter

import (
	"errors"
	"strings"
	"time"
)

type CreatedAtProvider interface {
	GetCreatedAt() int64
}

type LenProvider interface {
	GetLen() int
}

type NameProvider interface {
	GetName() string
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
		if unixTimeXDaysAgo <= obj.GetCreatedAt() {
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
		if unixTimeXDaysAgo > obj.GetCreatedAt() {
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
		if float64(obj.GetLen()) <= length {
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
		if float64(obj.GetLen()) > length {
			filtered = append(filtered, obj)
		}
	}

	return &filtered, nil
}

func ContainsName[T NameProvider](list *[]T, names []string) (*[]T) {
	filtered := []T{}

	for _, obj := range *list {
		
		contains := true;
		for _, name := range names {
			if !strings.Contains(obj.GetName(), name) {
				contains = false;
				break;
			} 
		}

		if contains {
			filtered = append(filtered, obj)
		}
	}

	return &filtered

}

func NotContainsName[T NameProvider](list *[]T, names []string) (*[]T) {
	filtered := []T{}

	for _, obj := range *list {
		contains := false;

		for _, name := range names {
			if (strings.Contains(obj.GetName(), name)) {
				contains = true;
				break;
			}
		}

		if (!contains) {
			filtered = append(filtered, obj)
		}
	}

	return &filtered

}
