package valueObjects

import "fmt"

type Status string

const (
	Error   = "error"
	Success = "Success"
)

var allowedStrings = []string{Error, Success}

func NewStatus(status string) (Status, error) {
	for _, allowedString := range allowedStrings {
		if status == allowedString {
			return Status(status), nil
		}
	}
	return "", fmt.Errorf("value %v not allowed for type Status", status)
}
