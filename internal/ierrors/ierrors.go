package ierrors

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"

	"errors"
)

type ErrResponse struct {
	Err        error  `json:"-"`
	HTTPCode   int    `json:"-"`
	StatusText string `json:"statusText"`
	AppCode    int    `json:"code,omitempty"`
	ErrorText  string `json:"message,omitempty"`
}

func NewErrorResponse(code int, message string) *ErrResponse {
	return &ErrResponse{
		HTTPCode:   code,
		StatusText: http.StatusText(code),
		ErrorText:  message,
	}
}

func (e *ErrResponse) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.StatusText
}

var (
	ErrBadRequest          = &ErrResponse{HTTPCode: http.StatusBadRequest, StatusText: "Bad Request"}
	ErrInternalServerError = &ErrResponse{HTTPCode: http.StatusInternalServerError, StatusText: "Internal Server Error"}
	ErrNotFound            = &ErrResponse{HTTPCode: http.StatusNotFound, StatusText: "Not Found", ErrorText: "No receipt found for that ID"}
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e FieldError) Error() string {
	return fmt.Sprintf("field %s: %s", e.Field, e.Message)
}

type ValidationErrors []FieldError

func (errs ValidationErrors) Error() string {
	if len(errs) == 0 {
		return "validation successful"
	}
	var sb strings.Builder
	sb.WriteString("validation errors:")
	for _, err := range errs {
		sb.WriteString("\n- ")
		sb.WriteString(err.Error())
	}
	return sb.String()
}

type ValidationErrorResponse struct {
	HTTPCode   int          `json:"-"`
	StatusText string       `json:"statusText"`
	ErrorText  string       `json:"errorText"`
	Errors     []FieldError `json:"errors"`
}

func NewValidationErrorResponse(err error) *ValidationErrorResponse {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fieldErrors := make([]FieldError, 0, len(validationErrors))
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			message := fieldError.Error()
			fieldErrors = append(fieldErrors, FieldError{
				Field:   field,
				Message: message,
			})
		}
		return &ValidationErrorResponse{
			HTTPCode:   http.StatusBadRequest,
			StatusText: "Bad Request",
			ErrorText:  "The receipt is invalid.",
			Errors:     fieldErrors,
		}
	}

	return &ValidationErrorResponse{
		HTTPCode:   http.StatusBadRequest,
		StatusText: "Bad Request",
		ErrorText:  "The receipt is invalid.",
		Errors:     []FieldError{{Message: err.Error()}},
	}
}
