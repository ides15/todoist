package main

type Response struct {
	SyncToken     string         `json:"sync_token"`
	TempIDMapping map[string]int `json:"temp_id_mapping"`
	FullSync      bool           `json:"full_sync"`
	Projects      *[]Project     `json:"projects"`
}
