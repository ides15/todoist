package types

type Response struct {
	SyncToken string     `json:"sync_token"`
	FullSync  bool       `json:"full_sync"`
	Projects  []*Project `json:"projects"`
}
