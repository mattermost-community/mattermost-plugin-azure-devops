package plugin

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-plugin-azure-devops/mocks"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/testutils"
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

func TestPostReleaseDetailsPreview(t *testing.T) {
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(mockAPI, nil, mockedClient)
	for _, testCase := range []struct {
		description string
		linkData    []string
		err         error
		statusCode  int
	}{
		{
			description: "PostReleaseDetailsPreview: valid",
			linkData:    []string{"https:", "", "test.com", "abc", "xyz", "_releaseProgress?_a=release-pipeline-progress&releaseId=20"},
			statusCode:  http.StatusOK,
		},
		{
			description: "PostReleaseDetailsPreview: invalid",
			linkData:    []string{"https:", "", "text.com", "abc", "xyz", "_releaseProgress?_a=release-pipeline-progress&releaseId=20"},
			err:         errors.New("failed to post release details preview"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
			mockedClient.EXPECT().GetReleaseDetails(gomock.Any(), gomock.Any(), gomock.Any(), testutils.MockMattermostUserID).Return(&serializers.ReleaseDetails{}, testCase.statusCode, testCase.err)

			resp, msg := p.PostReleaseDetailsPreview(testCase.linkData, "mockReleasePipelineLink", testutils.MockMattermostUserID, testutils.MockChannelID)
			assert.Equal(t, "", msg)
			if testCase.err != nil {
				assert.Nil(t, resp)
			} else {
				assert.NotNil(t, resp)
			}
		})
	}
}
