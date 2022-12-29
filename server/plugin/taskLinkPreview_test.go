package plugin

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-plugin-azure-devops/mocks"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
)

func TestPostTaskPreview(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(nil, nil, mockedClient)
	linkData := []string{"https:", "", "test.com", "abc", "xyz", "_workitems", "edit", "1"}

	t.Run("PostTaskPreview: valid", func(t *testing.T) {
		mockedClient.EXPECT().GetTask(gomock.Any(), gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.TaskValue{}, http.StatusOK, nil)
		resp, err := p.PostTaskPreview(linkData, "mockUserID", "mockChannelID")
		assert.Equal(t, "", err)
		assert.NotNil(t, resp)
	})
}

func TestPostPullRequestPreview(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(nil, nil, mockedClient)
	linkData := []string{"https:", "", "test.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"}

	t.Run("PostPullRequestPreview: valid", func(t *testing.T) {
		mockedClient.EXPECT().GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.PullRequest{}, http.StatusOK, nil)
		resp, stringErr := p.PostPullRequestPreview(linkData, "mockPullRequestLink", "mockUserID", "mockChannelID")
		assert.Equal(t, "", stringErr)
		assert.NotNil(t, resp)
	})
}

func TestPostBuildDetailsPreview(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(nil, nil, mockedClient)
	linkData := []string{"https:", "", "test.com", "abc", "xyz", "_build", "results?buildId=50&view=results"}

	t.Run("PostBuildDetailsPreview: valid", func(t *testing.T) {
		mockedClient.EXPECT().GetBuildDetails(gomock.Any(), gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.BuildDetails{}, http.StatusOK, nil)
		resp, stringErr := p.PostBuildDetailsPreview(linkData, "mockBuildPipelineLink", "mockUserID", "mockChannelID")
		assert.Equal(t, "", stringErr)
		assert.NotNil(t, resp)
	})
}
