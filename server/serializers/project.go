package serializers

type ProjectDetails struct {
	MattermostUserID string `json:"mattermostUserID"`
	ProjectID        string `json:"projectID"`
	ProjectName      string `json:"projectName"`
	OrganizationName string `json:"organizationName"`
}
