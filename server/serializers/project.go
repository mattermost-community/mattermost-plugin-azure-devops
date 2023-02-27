package serializers

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
)

type ProjectDetails struct {
	MattermostUserID    string `json:"mattermostUserID"`
	ProjectID           string `json:"projectID"`
	ProjectName         string `json:"projectName"`
	OrganizationName    string `json:"organizationName"`
	DeleteSubscriptions bool   `json:"deleteSubscriptions"`
}

func (t *ProjectDetails) IsValid() error {
	if t.OrganizationName == "" {
		return errors.New(constants.OrganizationRequired)
	}
	if t.ProjectName == "" {
		return errors.New(constants.ProjectRequired)
	}
	if t.ProjectID == "" {
		return errors.New(constants.ProjectIDRequired)
	}
	return nil
}

func ProjectPayloadFromJSON(data io.Reader) (*ProjectDetails, error) {
	var body *ProjectDetails
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}
