package serializers

import (
	"encoding/json"
	"errors"
	"io"
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

type CreateTaskRequestPayload struct {
	Organization string               `json:"organization"`
	Project      string               `json:"project"`
	Type         string               `json:"type"`
	Fields       CreateTaskFieldValue `json:"fields"`
}

type CreateTaskFieldValue struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskBodyPayload struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	From      string `json:"from"`
	Value     string `json:"value"`
}

// IsValid function to validate request payload.
func (t *CreateTaskRequestPayload) IsValid() error {
	if t.Organization == "" {
		return errors.New(constants.OrganizationRequired)
	}
	if t.Project == "" {
		return errors.New(constants.ProjectRequired)
	}
	if t.Type == "" {
		return errors.New(constants.TaskTypeRequired)
	}
	if t.Fields.Title == "" {
		return errors.New(constants.TaskTitleRequired)
	}
	return nil
}

func CreateTaskRequestPayloadFromJSON(data io.Reader) (*CreateTaskRequestPayload, error) {
	var body *CreateTaskRequestPayload
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}
