package types

import "net/http"

type Response struct {
	Raw *http.Response

	SyncToken string     `json:"sync_token"`
	FullSync  bool       `json:"full_sync"`
	Projects  []*Project `json:"projects"`
}
