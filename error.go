package todoist

import (
	"encoding/json"
	"net/http"

	"github.com/ides15/todoist/types"
)

func CreateError(res *http.Response) (*types.HTTPError, error) {
	errorBody := &types.HTTPError{}
	if err := json.NewDecoder(res.Body).Decode(errorBody); err != nil {
		return nil, err
	}

	return errorBody, nil
}
