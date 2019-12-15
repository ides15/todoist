package todoist

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreateErrorOK(t *testing.T) {
	res := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
			"error_tag": "ERROR",
			"error_code": 1,
			"http_code": 400,
			"error": "error message",
			"error_extra": {
				"retry_after": 1
			}
		}`))),
	}

	httpErr, err := CreateError(res)
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}

	expected := "error message"
	if httpErr.ErrorMessage != expected {
		t.Fatalf("expected %s, received %s", expected, httpErr.ErrorMessage)
	}
}

func TestCreateErrorDecodingError(t *testing.T) {
	res := &http.Response{
		Body: http.NoBody,
	}

	_, err := CreateError(res)
	if err == nil {
		t.Fatalf("expected decoding error, received nil")
	}
}
