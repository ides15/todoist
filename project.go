package main

type Project struct {
	ID             int    `json:"id"`
	LegacyID       int    `json:"legacy_id"`
	Name           string `json:"name"`
	Color          int    `json:"color"`
	ParentID       int    `json:"parent_id"`
	LegacyParentID int    `json:"legacy_parent_id"`
	ChildOrder     int    `json:"child_order"`
	Collapsed      int    `json:"collapsed"`
	Shared         bool   `json:"shared"`
	IsDeleted      int    `json:"is_deleted"`
	IsArchived     int    `json:"is_archived"`
	IsFavorite     int    `json:"is_favorite"`
	InboxProject   bool   `json:"inbox_project"`
	TeamInbox      bool   `json:"team_inbox"`
}

type NewProject struct {
	Name       string `json:"name"`
	Color      int    `json:"color"`
	ParentID   int    `json:"parent_id"`
	ChildOrder int    `json:"child_order"`
	IsFavorite int    `json:"is_favorite"`
}

type UpdatedProject struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Color      int    `json:"color"`
	Collapsed  int    `json:"collapsed"`
	IsFavorite int    `json:"is_favorite"`
}

type MovedProject struct {
	ID       int `json:"id"`
	ParentID int `json:"parent_id"`
}

type DeletedProject struct {
	ID int `json:"id"`
}

type ReorderedProject struct {
	Projects *[]Reorder `json:"projects"`
}

type Reorder struct {
	ID         int `json:"id"`
	ChildOrder int `json:"child_order"`
}
