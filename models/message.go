package models

import (
	"fmt"
)

type Message map[string]any

type ErrorMessage struct {
	Message string
	Error   string
}

func GetErrorMessage(msg string, err error) ErrorMessage {
	return ErrorMessage{
		Message: msg,
		Error:   fmt.Sprintf("%s", err),
	}
}

