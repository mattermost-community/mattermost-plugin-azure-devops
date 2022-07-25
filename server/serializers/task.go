package serializers

import (
	"time"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

// TODO: WIP.
// type TaskIDList struct {
// 	TaskList []TaskIDListValue `json:"workItems"`
// }

// type TaskIDListValue struct {
// 	ID int `json:"id"`
// }

// type TaskList struct {
// 	Count int         `json:"count"`
// 	Tasks []TaskValue `json:"value"`
// }

type TaskValue struct {
	ID     int            `json:"id"`
	Fields TaskFieldValue `json:"fields"`
	Link   Link           `json:"_links"`
}

type TaskFieldValue struct {
	Title       string          `json:"System.Title"`
	Project     string          `json:"System.TeamProject"`
	Type        string          `json:"System.WorkItemType"`
	State       string          `json:"System.State"`
	Reason      string          `json:"System.Reason"`
	AssignedTo  TaskUserDetails `json:"System.AssignedTo"`
	CreatedAt   time.Time       `json:"System.CreatedDate"`
	CreatedBy   TaskUserDetails `json:"System.CreatedBy"`
	UpdatedAt   time.Time       `json:"System.ChangedDate"`
	UpdatedBy   TaskUserDetails `json:"System.ChangedBy"`
	Description string          `json:"System.Description"`
}

type Link struct {
	Html Href `json:"html"`
}

type Href struct {
	Href string `json:"href"`
}

type TaskUserDetails struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
}

type TaskCreateRequestPayload struct {
	Organization string               `json:"organization"`
	Project      string               `json:"project"`
	Type         string               `json:"type"`
	Feilds       TaskCreateFieldValue `json:"fields"`
}

type TaskCreateFieldValue struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskCreateBodyPayload struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	From      string `json:"from"`
	Value     string `json:"value"`
}

// IsValid function to validate request payload.
func (t *TaskCreateRequestPayload) IsValid() string {
	if t.Organization == "" {
		return constants.OrganizationRequired
	}
	if t.Project == "" {
		return constants.ProjectRequired
	}
	if t.Type == "" {
		return constants.TaskTypeRequired
	}
	if t.Feilds.Title == "" {
		return constants.TaskTitleRequired
	}
	return ""
}
