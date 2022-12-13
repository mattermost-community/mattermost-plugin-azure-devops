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
	p := Plugin{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p.Client = mockedClient
	for _, testCase := range []struct {
		description string
		linkData    []string
	}{
		{
			description: "PostTaskPreview: valid",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockedClient.EXPECT().GetTask(gomock.Any(), gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.TaskValue{}, http.StatusOK, nil)
			resp, err := p.PostTaskPreview(testCase.linkData, "mockUserID", "mockChannelID")
			assert.Equal(t, "", err)
			assert.NotNil(t, resp)
		})
	}
}

func TestPostPullRequestPreview(t *testing.T) {
	p := Plugin{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p.Client = mockedClient
	for _, testCase := range []struct {
		description string
		linkData    []string
	}{
		{
			description: "PostPullRequestPreview: valid",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockedClient.EXPECT().GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.PullRequest{}, http.StatusOK, nil)
			resp, stringErr := p.PostPullRequestPreview(testCase.linkData, "mockPullRequestLink", "mockUserID", "mockChannelID")
			assert.Equal(t, "", stringErr)
			assert.NotNil(t, resp)
		})
	}
}

func TestPostReleaseDetailsPreview(t *testing.T) {
	p := Plugin{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p.Client = mockedClient
	for _, testCase := range []struct {
		description string
		linkData    []string
	}{
		{
			description: "PostReleaseDetailsPreview: valid",
			linkData:    []string{"https:", "", "dev.azure.com", "abc", "xyz", "_releaseProgress?_a=release-pipeline-progress&releaseId=20"},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockedClient.EXPECT().GetReleaseDetails(gomock.Any(), gomock.Any(), gomock.Any(), "mockUserID").Return(&serializers.ReleaseDetails{}, http.StatusOK, nil)
			resp, stringErr := p.PostReleaseDetailsPreview(testCase.linkData, "mockReleasePipelineLink", "mockUserID", "mockChannelID")
			assert.Equal(t, "", stringErr)
			assert.NotNil(t, resp)
		})
	}
}
