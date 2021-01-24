package gomodprivate

import (
	"errors"
	"strings"
)

func _ExtractTag(packageName string) (string, string, error) {
	if len(packageName) == 0 {
		return "", "", errors.New("Package Name can't be empty")
	}
	splitted := strings.Split(packageName, "@")
	if len(splitted) == 1 {
		return splitted[0], "", nil
	}
	return splitted[0], splitted[1], nil
}
