package endpoints

import (
	"strings"
)

func ParseEndpoint(fullPath string) string {
	splitPath := strings.Split(fullPath, "?")
	if splitPath == nil {
		return fullPath
	}

	return splitPath[0]
}
