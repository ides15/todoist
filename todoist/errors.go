package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// BaseError reports an error caused by a Todoist (sync) API request. BaseError
// will report on any status code response, therefore it is used for non-200 (OK)
// status code responses.
type BaseError struct {
	Response *http.Response `json:"-"` // HTTP response that caused this error

	Tag        string                 `json:"error_tag"`   // error tag
	Code       int                    `json:"error_code"`  // error code
	Message    string                 `json:"error"`       // error message
	HTTPCode   int                    `json:"http_code"`   // error HTTP code
	ErrorExtra map[string]interface{} `json:"error_extra"` // more detail on errors
}

func (e BaseError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.HTTPCode, e.Tag, e.Message)
}

// BadRequestError is used if the request was incorrect.
type BadRequestError struct {
	BaseError
}

func (e BadRequestError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.HTTPCode, e.Tag, e.Message)
}

// UnauthorizedError is used if authentication is required, and has failed, or has not yet been provided.
type UnauthorizedError struct {
	BaseError
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.HTTPCode, e.Tag, e.Message)
}

// ForbiddenError is used if the request was valid, but for something that is forbidden.
type ForbiddenError struct {
	BaseError
}

func (e ForbiddenError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.HTTPCode, e.Tag, e.Message)
}

// NotFoundError is used if the requested resource could not be found.
type NotFoundError struct {
	BaseError
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.HTTPCode, e.Tag, e.Message)
}

// TooManyRequestsError is used if the user has sent too many requests in a given amount of time.
type TooManyRequestsError struct {
	BaseError
}

func (e TooManyRequestsError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.HTTPCode, e.Tag, e.Message)
}

// InternalServerError is used if the request failed due to a server error.
type InternalServerError struct {
	BaseError
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.HTTPCode, e.Tag, e.Message)
}

// ServiceUnavailableError is used if the server is currently unable to handle the request.
type ServiceUnavailableError struct {
	BaseError
}

func (e ServiceUnavailableError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.HTTPCode, e.Tag, e.Message)
}

// SyncError reports an error caused by a Todoist (sync) API request
// with a 200 (OK) status code response, and contains an embedded BaseError.
// Todoist API docs: https://developer.todoist.com/sync/v8/?shell#response-error
type SyncError struct {
	BaseError // embedded original error

	ID string `json:"-"` // original command UUID
}

// checkResponseForErrors checks the API response for an error, and returns it if
// present. A response is considered an error if it has a status code not equal
// to 200 OK, or it has values in the `sync_status` field that are not equal
// to "ok".
//
// API error responses are expected to have response bodies, and a JSON response
// body that maps to SyncError.
func checkResponseForErrors(r *http.Response, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// TODO: handle this nicer
		return err
	}

	if body == nil {
		return errors.New("response body is nil")
	}

	switch c := r.StatusCode; c {
	// The request was processed successfully.
	// In the Todoist API, a 200 (OK) status code means that the response is at least
	// partially correct. If the response is a CommandResponse, there might be errors
	// in the sync_status field, so we still need to check that field for any errors.
	case http.StatusOK:
		if err = json.Unmarshal(body, &v); err != nil {
			// TODO: handle this nicer
			return err
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		switch vType := v.(type) {
		case CommandResponse:
			cr := vType

			// Range through each of the sync_status values, and map
			// each non "ok" value to a SyncError struct
			for cmdID, cmdResult := range cr.SyncStatus {
				if cmdResult != "ok" {
					// Serialize the command result back into an "unmarshallable" string
					cmdResultBytes, _ := json.Marshal(cmdResult)

					var syncErr SyncError
					if err = json.Unmarshal(cmdResultBytes, &syncErr); err != nil {
						return err
					}

					syncErr.ID = cmdID

					return syncErr
				}
			}

			return nil

		default:
			return nil
		}

	// The request was incorrect.
	case http.StatusBadRequest:
		var badRequestError BadRequestError
		if err := json.Unmarshal(body, &badRequestError); err != nil {
			return err
		}

		return badRequestError

	// Authentication is required, and has failed, or has not yet been provided.
	case http.StatusUnauthorized:
		var unauthorizedError UnauthorizedError
		if err := json.Unmarshal(body, &unauthorizedError); err != nil {
			return err
		}

		return unauthorizedError

	// The request was valid, but for something that is forbidden.
	case http.StatusForbidden:
		var forbiddenError ForbiddenError
		if err := json.Unmarshal(body, &forbiddenError); err != nil {
			return err
		}

		return forbiddenError

	// The requested resource could not be found.
	case http.StatusNotFound:
		var notFoundError NotFoundError
		if err := json.Unmarshal(body, &notFoundError); err != nil {
			return err
		}

		return notFoundError

	// The user has sent too many requests in a given amount of time.
	case http.StatusTooManyRequests:
		var tooManyRequestsError TooManyRequestsError
		if err := json.Unmarshal(body, &tooManyRequestsError); err != nil {
			return err
		}

		return tooManyRequestsError

	// The request failed due to a server error.
	case http.StatusInternalServerError:
		var internalServerError InternalServerError
		if err := json.Unmarshal(body, &internalServerError); err != nil {
			return err
		}

		return internalServerError

	// The server is currently unable to handle the request.
	case http.StatusServiceUnavailable:
		var serviceUnavailableError ServiceUnavailableError
		if err := json.Unmarshal(body, &serviceUnavailableError); err != nil {
			return err
		}

		return serviceUnavailableError

	default:
		unknownError := BaseError{
			Response:   r,
			Tag:        "UNKNOWN_ERROR",
			Message:    "Unknown error occured.",
			HTTPCode:   r.StatusCode,
			ErrorExtra: map[string]interface{}{},
		}

		return unknownError
	}
}
