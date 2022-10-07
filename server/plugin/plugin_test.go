package plugin

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitBotUser(t *testing.T) {
	p := Plugin{}
	mockHelper := &plugintest.Helpers{}
	p.Helpers = mockHelper
	for _, testCase := range []struct {
		description string
	}{
		{
			description: "InitBotUser: valid",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockHelper.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.Anything).Return("mockBotID", nil)

			resp := p.initBotUser()
			assert.Nil(t, resp)
		})
	}
}

func TestMessageWillBePosted(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}
	for _, testCase := range []struct {
		description            string
		message                string
		taskData               []string
		pullRequestData        []string
		isValidTaskLink        bool
		isValidPullRequestLink bool
	}{
		{
			description:     "test change post for valid link",
			taskData:        []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			message:         "mockMessage",
			isValidTaskLink: true,
		},
		{
			description:            "test change post for valid link",
			pullRequestData:        []string{"https:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			message:                "mockMessage",
			isValidPullRequestLink: true,
		},
		{
			description: "MessageWillBePosted: invalid link",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(isValidTaskLink, func(_ string) ([]string, bool) {
				return testCase.taskData, testCase.isValidTaskLink
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "PostTaskPreview", func(_ *Plugin, _ []string, _, _ string) (*model.Post, string) {
				return &model.Post{}, testCase.message
			})
			monkey.Patch(isValidPullRequestLink, func(_ string) ([]string, bool) {
				return testCase.pullRequestData, testCase.isValidPullRequestLink
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "PostPullRequestPreview", func(_ *Plugin, _ []string, _, _, _ string) (*model.Post, string) {
				return &model.Post{}, testCase.message
			})

			post := &model.Post{
				ChannelId: "mockChannelID",
				UserId:    "mockUserID",
				Message:   testCase.message,
			}

			newPost, _ := p.MessageWillBePosted(&plugin.Context{}, post)
			if testCase.isValidTaskLink || testCase.isValidPullRequestLink {
				assert.NotNil(t, newPost)
				return
			}

			assert.Nil(t, newPost)
		})
	}
}

func TestIsValidTaskLink(t *testing.T) {
	for _, testCase := range []struct {
		description  string
		msg          string
		expectedData []string
		isValid      bool
	}{
		{
			description:  "IsValidTaskLink: valid link 1",
			msg:          "https://dev.azure.com/abc/xyz/_workitems/edit/1/",
			expectedData: []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			isValid:      true,
		},
		{
			description:  "IsValidTaskLink: valid link 2",
			msg:          "https://dev.azure.com/abc/xyz/_workitems/edit/1",
			expectedData: []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			isValid:      true,
		},
		{
			description:  "IsValidTaskLink: valid link 3",
			msg:          "http://dev.azure.com/abc/xyz/_workitems/edit/1",
			expectedData: []string{"http:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			isValid:      true,
		},
		{
			description: "IsValidTaskLink: invalid link 1",
			msg:         "https://abc/xyz/_workitems/edit/1",
		},
		{
			description: "IsValidTaskLink: invalid link 2",
			msg:         "https://dev.azure.com/abc/xyz/_workitems/edit",
		},
		{
			description: "IsValidTaskLink: invalid link 3",
			msg:         "https://dev.azure.com/xyz/_workitems/edit/1",
		},
		{
			description: "IsValidTaskLink: invalid link 4",
			msg:         "http://dev.azure/abc/xyz/items/it/1",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			data, isValid := isValidTaskLink(testCase.msg)
			assert.Equal(t, testCase.expectedData, data)
			assert.Equal(t, testCase.isValid, isValid)
		})
	}
}

func TestIsValidPullRequestLink(t *testing.T) {
	for _, testCase := range []struct {
		description  string
		msg          string
		expectedData []string
		isValid      bool
	}{
		{
			description:  "IsValidPullRequestLink: valid link 1",
			msg:          "https://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1/",
			expectedData: []string{"https:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			isValid:      true,
		},
		{
			description:  "IsValidPullRequestLink: valid link 2",
			msg:          "https://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1",
			expectedData: []string{"https:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			isValid:      true,
		},
		{
			description:  "IsValidPullRequestLink: valid link 3",
			msg:          "http://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1",
			expectedData: []string{"http:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			isValid:      true,
		},
		{
			description: "IsValidPullRequestLink: invalid link 1",
			msg:         "https://abc/xyz/_git/pullrequest/1",
		},
		{
			description: "IsValidPullRequestLink: invalid link 2",
			msg:         "https://dev.azure.com/abc/xyz/_git/pullrequest",
		},
		{
			description: "IsValidPullRequestLink: invalid link 3",
			msg:         "https://dev.azure.com/xyz/_git/xyz/pullrequest/1",
		},
		{
			description: "IsValidPullRequestLink: invalid link 4",
			msg:         "http://dev.azure/abc/xyz/pull/it/1",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			data, isValid := isValidPullRequestLink(testCase.msg)
			assert.Equal(t, testCase.expectedData, data)
			assert.Equal(t, testCase.isValid, isValid)
		})
	}
}
