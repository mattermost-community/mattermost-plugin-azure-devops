package serializers

import (
	"encoding/json"
	"io"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

type ProjectDetails struct {
	MattermostUserID string `json:"mattermostUserID"`
	ProjectID        string `json:"projectID"`
	ProjectName      string `json:"projectName"`
	OrganizationName string `json:"organizationName"`
}

func (t *ProjectDetails) IsValid() string {
	if t.OrganizationName == "" {
		return constants.OrganizationRequired
	}
	if t.ProjectName == "" {
		return constants.ProjectRequired
	}
	if t.ProjectID == "" {
		return constants.ProjectIDRequired
	}
	return ""
}

func ProjectPayloadFromJSON(data io.Reader) (*ProjectDetails, error) {
	var body *ProjectDetails
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}
