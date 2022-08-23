package plugin

import (
	"net/http"
	"testing"

	"github.com/Brightscout/mattermost-plugin-azure-devops/mocks"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPostTaskPreview(t *testing.T) {
	p := Plugin{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p.Client = mockedClient
	for _, testCase := range []struct {
		description string
		linkData    []string
	}{
		{
			description: "test CreateTask",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockedClient.EXPECT().GetTask(gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.TaskValue{}, http.StatusOK, nil)
			resp, err := p.PostTaskPreview(testCase.linkData, "mockUserID", "mockChannelID")
			assert.Equal(t, "", err)
			assert.NotNil(t, resp)
		})
	}
}
