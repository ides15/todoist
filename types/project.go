package types

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
	Name       string `json:"name,omitempty"`
	Color      int    `json:"color,omitempty"`
	ParentID   int    `json:"parent_id,omitempty"`
	ChildOrder int    `json:"child_order,omitempty"`
	IsFavorite int    `json:"is_favorite,omitempty"`
}

type UpdatedProject struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Color      int    `json:"color,omitempty"`
	Collapsed  int    `json:"collapsed,omitempty"`
	IsFavorite int    `json:"is_favorite,omitempty"`
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
