package testutils

import (
	"github.com/stretchr/testify/mock"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
)

const (
	MockOrganization     = "mockOrganization"
	MockProjectName      = "mockProjectName"
	MockMattermostUserID = "mockMattermostUserID"
	MockChannelID        = "mockChannelID"
)

func GetMockArgumentsWithType(typeString string, num int) []interface{} {
	ret := make([]interface{}, num)
	for i := 0; i < len(ret); i++ {
		ret[i] = mock.AnythingOfTypeArgument(typeString)
	}
	return ret
}

func GetSuscriptionDetailsPayload(userID, serviceType, eventType string) []*serializers.SubscriptionDetails {
	return []*serializers.SubscriptionDetails{
		{
			ChannelID:        "mockChannelID",
			MattermostUserID: userID,
			ServiceType:      serviceType,
			SubscriptionID:   "mockSubscriptionID",
			OrganizationName: "mockOrganization",
			ProjectName:      "mockProjectName",
			EventType:        eventType,
			CreatedBy:        "mockCreatedBy",
			ChannelName:      "mockChannelName",
		},
	}
}

func GetProjectDetailsPayload() []serializers.ProjectDetails {
	return []serializers.ProjectDetails{
		{
			MattermostUserID: "mockMattermostUserID",
			OrganizationName: "mockOrganization",
			ProjectName:      "mockProjectName",
			ProjectID:        "mockProjectID",
		},
	}
}
