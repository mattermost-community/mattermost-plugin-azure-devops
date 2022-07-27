package serializers

type ProjectDetails struct {
	MattermostUserID string `json:"mattermostUserID"`
	ProjectID        string `json:"projectID"`
	ProjectName      string `json:"projectName"`
	OrganizationName string `json:"organizationName"`
}

// TODO: Remove later if not needed.
// import (
// 	"time"
// )

// type ProjectList struct {
// 	Count        int            `json:"count"`
// 	ProjectValue []ProjectValue `json:"value"`
// }

// type ProjectValue struct {
// 	ID             string    `json:"id"`
// 	URL            string    `json:"url"`
// 	Name           string    `json:"name"`
// 	State          string    `json:"state"`
// 	Revision       int       `json:"revision"`
// 	Visibility     string    `json:"visibility"`
// 	LastUpdateTime time.Time `json:"lastUpdateTime"`
// }
