package serializers

import (
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

type LinkRequestPayload struct {
	Organization string `json:"organization"`
	Project      string `json:"project"`
}

type Project struct {
	ID   string      `json:"id"`
	Name string      `json:"name"`
	Link ProjectLink `json:"_links"`
}

type ProjectLink struct {
	Web Href `json:"web"`
}

// IsLinkPayloadValid function to validate request payload.
func (t *LinkRequestPayload) IsLinkPayloadValid() string {
	if t.Organization == "" {
		return constants.OrganizationRequired
	}
	if t.Project == "" {
		return constants.ProjectRequired
	}
	return ""
}
