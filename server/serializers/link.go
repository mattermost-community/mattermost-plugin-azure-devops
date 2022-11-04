package serializers

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

type LinkRequestPayload struct {
	Organization string `json:"organization"`
	Project      string `json:"project"`
}

type Project struct {
	ID   string      `json:"id"`
	Name string      `json:"name"`
	URL  string      `json:"url"`
	Link ProjectLink `json:"_links"`
}

type ProjectLink struct {
	Web Href `json:"web"`
}

// IsLinkPayloadValid function to validate request payload.
func (t *LinkRequestPayload) IsLinkPayloadValid() error {
	if t.Organization == "" {
		return errors.New(constants.OrganizationRequired)
	}
	if t.Project == "" {
		return errors.New(constants.ProjectRequired)
	}
	return nil
}

func LinkPayloadFromJSON(data io.Reader) (*LinkRequestPayload, error) {
	var body *LinkRequestPayload
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}
