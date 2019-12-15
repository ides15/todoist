package types

import (
	"errors"
	"fmt"
)

var (
	ErrRequiredToken = errors.New("must provide an API token")
	ErrBuildRequest  = errors.New("unable to build request")
	ErrUnknown       = errors.New("unknown error")
	ErrNotFound      = errors.New("not found")
)

type HTTPError struct {
	ErrorTag     string                 `json:"error_tag"`
	ErrorCode    int                    `json:"error_code"`
	HTTPCode     int                    `json:"http_code"`
	ErrorExtra   map[string]interface{} `json:"error_extra"`
	ErrorMessage string                 `json:"error"`
}

func (e *HTTPError) Error() string {
	return fmt.Sprint(e.ErrorMessage)
}
