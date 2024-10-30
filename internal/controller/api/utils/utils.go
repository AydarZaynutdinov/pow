package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
	_    Encoder = (*jsonEncoder)(nil)
	JSON Encoder = &jsonEncoder{}
)

type Encoder interface {
	Encode(w http.ResponseWriter, v any) error
	ContentType() string
}

type jsonEncoder struct{}

func (*jsonEncoder) Encode(w http.ResponseWriter, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func (*jsonEncoder) ContentType() string {
	return "application/json; charset=utf-8"
}

type ControllerResult struct {
	HTTPCode int
	Response any
}

func NewControllerResultOK(response any) *ControllerResult {
	return &ControllerResult{
		HTTPCode: http.StatusOK,
		Response: response,
	}
}

type ControllerStandardError struct {
	HTTPCode int
	Err      error
}

func NewControllerStandardErrorBadRequest(err error) *ControllerStandardError {
	return &ControllerStandardError{
		HTTPCode: http.StatusBadRequest,
		Err:      err,
	}
}

func NewControllerStandardErrorInternalServerError(err error) *ControllerStandardError {
	return &ControllerStandardError{
		HTTPCode: http.StatusInternalServerError,
		Err:      err,
	}
}

func (e *ControllerStandardError) Error() string {
	return e.Err.Error()
}

func (e *ControllerStandardError) Is(err error) bool {
	return errors.Is(e.Err, err)
}
