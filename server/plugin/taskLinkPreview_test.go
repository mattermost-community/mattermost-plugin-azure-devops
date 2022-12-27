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
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(mockAPI, nil, mockedClient)
	for _, testCase := range []struct {
		description string
		linkData    []string
		err         error
	}{
		{
			description: "PostTaskPreview: valid",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
		},
		{
			description: "PostTaskPreview: invalid",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)

			mockedClient.EXPECT().GetTask(gomock.Any(), gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.TaskValue{}, http.StatusOK, testCase.err)
			resp, msg := p.PostTaskPreview(testCase.linkData, "mockUserID", "mockChannelID")
			assert.Equal(t, "", msg)
			if testCase.err != nil {
				assert.Nil(t, resp)
			} else {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestPostPullRequestPreview(t *testing.T) {
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(mockAPI, nil, mockedClient)
	for _, testCase := range []struct {
		description string
		linkData    []string
		err         error
	}{
		{
			description: "PostPullRequestPreview: valid",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
		},
		{
			description: "PostPullRequestPreview: invalid",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)

			mockedClient.EXPECT().GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.PullRequest{}, http.StatusOK, testCase.err)
			resp, msg := p.PostPullRequestPreview(testCase.linkData, "mockPullRequestLink", "mockUserID", "mockChannelID")
			assert.Equal(t, "", msg)
			if testCase.err != nil {
				assert.Nil(t, resp)
			} else {
				assert.NotNil(t, resp)
			}
		})
	}
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
	}{
		{
			description: "PostReleaseDetailsPreview: valid",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_releaseProgress?_a=release-pipeline-progress&releaseId=20"},
		},
		{
			description: "PostReleaseDetailsPreview: invalid",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_releaseProgress?_a=release-pipeline-progress&releaseId=20"},
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)

			mockedClient.EXPECT().GetReleaseDetails(gomock.Any(), gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.ReleaseDetails{}, http.StatusOK, testCase.err)
			resp, msg := p.PostReleaseDetailsPreview(testCase.linkData, "mockReleasePipelineLink", "mockUserID", "mockChannelID")
			assert.Equal(t, "", msg)
			if testCase.err != nil {
				assert.Nil(t, resp)
			} else {
				assert.NotNil(t, resp)
			}
		})
	}
}
