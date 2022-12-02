package testutils

import (
	"github.com/stretchr/testify/mock"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
)

func GetMockArgumentsWithType(typeString string, num int) []interface{} {
	ret := make([]interface{}, num)
	for i := 0; i < len(ret); i++ {
		ret[i] = mock.AnythingOfTypeArgument(typeString)
	}
	return ret
}

func GetSuscriptionDetailsPayload(userID, serviceType string) []*serializers.SubscriptionDetails {
	return []*serializers.SubscriptionDetails{
		{
			ChannelID:        "mockChannelID",
			MattermostUserID: userID,
			ServiceType:      serviceType,
			SubscriptionID:   "mockSubscriptionID",
			OrganizationName: "mockOrganizationName",
			ProjectName:      "mockProjectName",
			EventType:        constants.SubscriptionEventWorkItemCreated,
			CreatedBy:        "mockCreatedBy",
			ChannelName:      "mockChannelName",
		},
	}
}

func GetGitBranchesPayload() []*serializers.GitBranch {
	return []*serializers.GitBranch{
		{
			ID:   "mockID-1",
			Name: "refs/heads/mockName-1",
		},
		{
			ID:   "mockID-2",
			Name: "refs/heads/mockName-2",
		},
	}
}
