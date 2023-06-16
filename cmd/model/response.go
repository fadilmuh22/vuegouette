package model

type BasicResponse struct {
	Success bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func NewErrorMessage(message string) *ErrorMessage {
	return &ErrorMessage{
		Message: message,
	}
}

func (e *ErrorMessage) Error() string {
	return e.Message
}
