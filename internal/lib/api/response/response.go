package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error,omitempty"`
}

const (
	StatusOK    = "Ok"
	StatusError = "Error"
)

func (r *Response) Error() string {
	return r.ErrorMessage
}

func OkMessage() Response {
	return Response{
		Status: StatusOK,
	}
}

func ErrorMessage(msg string) Response {
	return Response{
		Status:       StatusError,
		ErrorMessage: msg,
	}
}

func ValidationError(errs validator.ValidationErrors) *Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return &Response{
		Status:       StatusError,
		ErrorMessage: strings.Join(errMsgs, ", "),
	}
}
