package serializers

import "time"

type TasksIDList struct {
	TaskList []TaskIDListValue `json:"workItems"`
}

type TaskIDListValue struct {
	Id int `json:"id"`
}

type TasksList struct {
	Count int          `json:"count"`
	Tasks []TasksValue `json:"value"`
}

type TasksValue struct {
	Id     int             `json:"id"`
	Fields TasksFieldValue `json:"fields"`
}

type TasksFieldValue struct {
	Title      string           `json:"System.Title"`
	Project    string           `json:"System.TeamProject"`
	Type       string           `json:"System.WorkItemType"`
	State      string           `json:"System.State"`
	Reason     string           `json:"System.Reason"`
	AssignedTo TasksUserDetails `json:"System.AssignedTo"`
	CreatedAt  time.Time        `json:"System.CreatedDate"`
	CreatedBy  TasksUserDetails `json:"System.CreatedBy"`
	UpdatedAt  time.Time        `json:"System.ChangedDate"`
	UpdatedBy  TasksUserDetails `json:"System.ChangedBy"`
}

type TasksUserDetails struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
}
