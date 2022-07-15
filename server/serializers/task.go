package serializers

import "time"

type TaskIDList struct {
	TaskList []TaskIDListValue `json:"workItems"`
}

type TaskIDListValue struct {
	ID int `json:"id"`
}

type TaskList struct {
	Count int         `json:"count"`
	Tasks []TaskValue `json:"value"`
}

type TaskValue struct {
	ID     int            `json:"id"`
	Fields TaskFieldValue `json:"fields"`
}

type TaskFieldValue struct {
	Title      string          `json:"System.Title"`
	Project    string          `json:"System.TeamProject"`
	Type       string          `json:"System.WorkItemType"`
	State      string          `json:"System.State"`
	Reason     string          `json:"System.Reason"`
	AssignedTo TaskUserDetails `json:"System.AssignedTo"`
	CreatedAt  time.Time       `json:"System.CreatedDate"`
	CreatedBy  TaskUserDetails `json:"System.CreatedBy"`
	UpdatedAt  time.Time       `json:"System.ChangedDate"`
	UpdatedBy  TaskUserDetails `json:"System.ChangedBy"`
}

type TaskUserDetails struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
}
