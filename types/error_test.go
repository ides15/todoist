package types_test

import (
	"fmt"
	"testing"

	"github.com/ides15/todoist/types"
)

func TestErrorString(t *testing.T) {
	message := "This is the error message"
	err := types.HTTPError{
		ErrorMessage: message,
	}

	value := fmt.Sprint(err.Error())
	expected := fmt.Sprint(message)

	if value != expected {
		t.Fatalf("expected %s, received %s", expected, value)
	}
}
