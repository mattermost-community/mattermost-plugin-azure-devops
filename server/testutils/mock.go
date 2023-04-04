package testutils

import (
	"github.com/stretchr/testify/mock"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
)

const (
	MockOrganization      = "mockOrganization"
	MockProjectName       = "mockProjectName"
	MockMattermostUserID  = "mockMattermostUserID"
	MockChannelID         = "mockChannelID"
	MockProjectID         = "mockProjectID"
	MockEventType         = "mockEventType"
	MockSubscriptionID    = "mockSubscriptionID"
	MockServiceType       = "mockServiceType"
	MockApproverID        = "mockApproverID"
	MockAzureDevopsUserID = "mockAzureDevopsUserID"
	MockTeamID            = "mockTeamID"
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
			ChannelID:        MockChannelID,
			MattermostUserID: userID,
			ServiceType:      serviceType,
			SubscriptionID:   MockSubscriptionID,
			OrganizationName: MockOrganization,
			ProjectName:      MockProjectName,
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

func GenerateStringOfSize(sizeInByte int) string {
	bytes := make([]byte, sizeInByte)
	for i := range bytes {
		bytes[i] = 'a'
	}

	return string(bytes)
}
