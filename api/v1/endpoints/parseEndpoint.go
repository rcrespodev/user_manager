package endpoints

import (
	"fmt"
	"strings"
)

func ParseEndpoint(fullPath string) string {
	splitPath := strings.Split(fullPath, "?")
	if splitPath == nil {
		return fullPath
	}

	return splitPath[0]
}

func BuildEndpointKey(relPath, method string) string {
	return fmt.Sprintf("%s-%s", method, relPath)
}
