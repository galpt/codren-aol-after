package main

import (
	"errors"
	"strings"
)

// custom strings.Split
func kemoSplit(s, sep string) ([]string, error) {
	result := strings.Split(s, sep)

	if len(result) == 1 {
		return nil, errors.New("delimiter not found")
	}

	return result, nil
}

// custom func to remove content from a slice
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
