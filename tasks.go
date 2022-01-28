package todoist

import (
	"context"

	"github.com/google/uuid"
)

// TasksService handles communication with the tasks related
// methods of the Todoist API.
//
// Todoist API docs: https://developer.todoist.com/sync/v8/?shell#items
type TasksService service

type Task struct {
	// The ID of the task.
	ID int `json:"id"`

	// The legacy ID of the task
	// (only shown for objects created before 1 April 2017)
	LegacyID *int `json:"legacy_id"`

	// The owner of the task.
	UserID int `json:"user_id"`

	// The ID of the parent project.
	ProjectID int `json:"project_id"`

	// Legacy project ID for the project that the task resides in
	// (only shown for objects created before 1 April 2017)
	LegacyProjectID *int `json:"legacy_project_id"`

	// The text of the task. This value may contain markdown-formatted text and hyperlinks. Details on markdown support can be found in the Text Formatting article in the Help Center.
	Content string `json:"content"`

	// A description for the task. This value may contain markdown-formatted text and hyperlinks. Details on markdown support can be found in the Text Formatting article in the Help Center.
	Description string `json:"description"`

	// TODO: Date type
	// The due date of the task. See the Due dates section for more details.
	Due *interface{} `json:"due"`

	// The priority of the task (a number between 1 and 4, 4 for very urgent and 1 for natural).
	// Note: Keep in mind that very urgent is the priority 1 on clients. So, p1 will return 4 in the API.
	Priority int `json:"priority"`

	// The ID of the parent task. Set to null for root tasks.
	ParentID *int `json:"parent_id"`

	// The legacy ID of the parent task. Set to null for root tasks
	// (only shown for objects created before 1 April 2017)
	LegacyParentID *int `json:"legacy_parent_id"`

	// The order of the task. Defines the position of the task among all the tasks with the same parent.
	ChildOrder int `json:"child_order"`

	// The ID of the parent section. Set to null for tasks not belonging to a section.
	SectionID *int `json:"section_id"`

	// The order of the task inside the Today or Next 7 days view (a number, where the smallest value would place the task at the top).
	DayOrder int `json:"day_order"`

	// Whether the task's sub-tasks are collapsed (where 1 is true and 0 is false).
	Collapsed int `json:"collapsed"`

	// The task's labels (a list of label IDs such as [2324,2525]).
	Labels []int `json:"labels"`

	// The ID of the user who created the task. This makes sense for shared projects only. For tasks created before 31 Oct 2019 the value is set to null. Cannot be set explicitly or changed via API.
	AddedByUID *int `json:"added_by_uid"`

	// The ID of the user who assigned the task. This makes sense for shared projects only. Accepts any user ID from the list of project collaborators. If this value is unset or invalid, it will automatically be set up to your uid.
	AssignedByUID *int `json:"assigned_by_uid"`

	// The ID of user who is responsible for accomplishing the current task. This makes sense for shared projects only. Accepts any user ID from the list of project collaborators or null or an empty string to unset.
	ResponsibleUID *int `json:"responsible_uid"`

	// Whether the task is marked as completed (where 1 is true and 0 is false).
	Checked int `json:"checked"`

	// Whether the task has been marked as completed and is marked to be moved to history, because all the child tasks of its parent are also marked as completed (where 1 is true and 0 is false)
	InHistory int `json:"in_history"`

	// Whether the task is marked as deleted (where 1 is true and 0 is false).
	IsDeleted int `json:"is_deleted"`

	// Identifier to find the match between tasks in shared projects of different collaborators. When you share a task, its copy has a different ID in the projects of your collaborators. To find a task in another account that matches yours, you can use the "sync_id" attribute. For non-shared tasks, the attribute is null.
	SyncID *int `json:"sync_id"`

	// The date when the task was completed (or null if not completed).
	DateCompleted *string `json:"date_completed"`

	// The date when the task was created.
	DateAdded string `json:"date_added"`
}

// List the tasks for a user.
func (s *TasksService) List(ctx context.Context, syncToken string) ([]Task, ReadResponse, error) {
	s.client.Logln("---------- Tasks.List")

	req, err := s.client.NewRequest(syncToken, []string{"items"}, nil)
	if err != nil {
		return nil, ReadResponse{}, err
	}

	var readResponse ReadResponse
	_, err = s.client.Do(ctx, req, &readResponse)
	if err != nil {
		return nil, readResponse, err
	}

	return readResponse.Tasks, readResponse, nil
}

type AddTask struct {
	// The text of the task. This value may contain markdown-formatted text and hyperlinks. Details on markdown support can be found in the Text Formatting article in the Help Center.
	Content string `json:"content"`

	// A description for the task. This value may contain markdown-formatted text and hyperlinks. Details on markdown support can be found in the Text Formatting article in the Help Center.
	Description string `json:"description,omitempty"`

	// The ID of the project to add the task to (a number or a temp id). By default the task is added to the user’s Inbox project.
	ProjectID string `json:"project_id,omitempty"`

	// TODO: Date type
	// The due date of the task. See the Due dates section for more details.
	Due *interface{} `json:"due,omitempty"`

	// The priority of the task (a number between 1 and 4, 4 for very urgent and 1 for natural).
	// Note: Keep in mind that very urgent is the priority 1 on clients. So, p1 will return 4 in the API.
	Priority int `json:"priority,omitempty"`

	// The ID of the parent task. Set to null for root tasks.
	ParentID *int `json:"parent_id,omitempty"`

	// The order of task. Defines the position of the task among all the tasks with the same parent.
	ChildOrder int `json:"child_order,omitempty"`

	// The ID of the section. Set to null for tasks not belonging to a section.
	SectionID *int `json:"section_id,omitempty"`

	// The order of the task inside the Today or Next 7 days view (a number, where the smallest value would place the task at the top).
	DayOrder int `json:"day_order,omitempty"`

	// Whether the task's sub-tasks are collapsed (where 1 is true and 0 is false).
	Collapsed int `json:"collapsed,omitempty"`

	// The tasks labels (a list of label IDs such as [2324,2525]).
	Labels []int `json:"labels,omitempty"`

	// The ID of user who assigns the current task. This makes sense for shared projects only. Accepts 0 or any user ID from the list of project collaborators. If this value is unset or invalid, it will be automatically setup to your uid.
	AssignedByUID int `json:"assigned_by_uid,omitempty"`

	// The ID of user who is responsible for accomplishing the current task. This makes sense for shared projects only. Accepts any user ID from the list of project collaborators or null or an empty string to unset.
	ResponsibleUID *int `json:"responsible_uid,omitempty"`

	// When this option is enabled, the default reminder will be added to the new item if it has a due date with time set. See also the auto_reminder user option for more info about the default reminder.
	AutoReminder bool `json:"auto_reminder,omitempty"`

	// When this option is enabled, the labels will be parsed from the task content and added to the task. In case the label doesn't exist, a new one will be created.
	AutoParseLabels bool `json:"auto_parse_labels,omitempty"`

	TempID string `json:"-"`
}

// Add a new task to a project.
func (s *TasksService) Add(ctx context.Context, syncToken string, addTask AddTask) ([]Task, CommandResponse, error) {
	s.client.Logln("---------- Tasks.Add")

	id := uuid.New().String()
	tempID := addTask.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	addCommand := Command{
		Type:   "item_add",
		Args:   addTask,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{addCommand}

	req, err := s.client.NewRequest(syncToken, []string{"items"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Tasks, commandResponse, nil
}
