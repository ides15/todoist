package todoist

import "context"

// SectionsService handles communication with the sections related
// methods of the Todoist API.
//
// Todoist API docs: https://developer.todoist.com/sync/v8/?shell#sections
type SectionsService service

// Section represnts a Todoist section.
type Section struct {
	// The ID of the section.
	ID int `json:"id"`

	// The name of the section.
	Name string `json:"name"`

	// Project that the section resides in.
	ProjectID int `json:"project_id"`

	// Legacy project ID for the project that the section resides in.
	// (only shown for objects created before 1 April 2017)
	LegacyProjectID *int `json:"legacy_project_id"`

	// The order of the section. Defines the position of the section among all the sections in the project.
	SectionOrder int `json:"section_order"`

	// Whether the section's tasks are collapsed (a true or false value).
	Collapsed bool `json:"collapsed"`

	// A special ID for shared sections (a number or null if not set). Used internally and can be ignored.
	SyncID *int `json:"sync_id"`

	// Whether the section is marked as deleted (a true or false value).
	IsDeleted bool `json:"is_deleted"`

	// Whether the section is marked as archived (a true or false value).
	IsArchived bool `json:"is_archived"`

	// The date when the section was archived (or null if not archived).
	DateArchived *string `json:"date_archived"`

	// The date when the section was created.
	DateAdded string `json:"date_added"`
}

func (s *SectionsService) List(ctx context.Context, syncToken string) ([]Section, ReadResponse, error) {
	s.client.Logln("---------- Sections.List")

	req, err := s.client.NewRequest(syncToken, []string{"sections"}, nil)
	if err != nil {
		return nil, ReadResponse{}, err
	}

	var readResponse ReadResponse
	_, err = s.client.Do(ctx, req, &readResponse)
	if err != nil {
		return nil, readResponse, err
	}

	return readResponse.Sections, readResponse, nil
}

// func (s *SectionsService) Add(ctx context.Context, syncToken string) ([])
