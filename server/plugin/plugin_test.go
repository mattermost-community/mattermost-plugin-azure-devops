package plugin

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
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
		description string
		message     string
		data        []string
		isValidLink bool
		link        string
	}{
		{
			description: "MessageWillBePosted: test change post for valid link",
			data:        []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			message:     "mockMessage",
			isValidLink: true,
		},
		{
			description: "MessageWillBePosted: test change post for valid link",
			data:        []string{"https:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			message:     "mockMessage",
			isValidLink: true,
			link:        "https://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1",
		},
		{
			description: "MessageWillBePosted: invalid link",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "PostTaskPreview", func(_ *Plugin, _ []string, _, _ string) (*model.Post, string) {
				return &model.Post{}, testCase.message
			})
			monkey.Patch(IsLinkPresent, func(_, _ string) ([]string, string, bool) {
				return testCase.data, testCase.link, testCase.isValidLink
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
			if testCase.isValidLink {
				assert.NotNil(t, newPost)
				return
			}

			assert.Nil(t, newPost)
		})
	}
}

func TestIsLinkPresent(t *testing.T) {
	for _, testCase := range []struct {
		description  string
		msg          string
		expectedData []string
		isValid      bool
		expectedLink string
		regex        string
	}{
		{
			description:  "IsLinkPresent: valid task link 1",
			msg:          "https://dev.azure.com/abc/xyz/_workitems/edit/1/",
			expectedData: []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			isValid:      true,
			regex:        constants.TaskLinkRegex,
			expectedLink: "https://dev.azure.com/abc/xyz/_workitems/edit/1",
		},
		{
			description:  "IsLinkPresent: valid task link 2",
			msg:          "https://dev.azure.com/abc/xyz/_workitems/edit/1",
			expectedData: []string{"https:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			isValid:      true,
			regex:        constants.TaskLinkRegex,
			expectedLink: "https://dev.azure.com/abc/xyz/_workitems/edit/1",
		},
		{
			description:  "IsLinkPresent: valid task link 3",
			msg:          "http://dev.azure.com/abc/xyz/_workitems/edit/1",
			expectedData: []string{"http:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			isValid:      true,
			regex:        constants.TaskLinkRegex,
			expectedLink: "http://dev.azure.com/abc/xyz/_workitems/edit/1",
		},
		{
			description:  "IsLinkPresent: valid task link 4",
			msg:          "\n\nhttp://dev.azure.com/abc/xyz/_workitems/edit/1   mock-text",
			expectedData: []string{"http:", "", "dev.azure.com", "abc", "xyz", "_workitems", "edit", "1"},
			isValid:      true,
			regex:        constants.TaskLinkRegex,
			expectedLink: "http://dev.azure.com/abc/xyz/_workitems/edit/1",
		},
		{
			description: "IsLinkPresent: invalid task link 1",
			msg:         "https://abc/xyz/_workitems/edit/1",
			regex:       constants.TaskLinkRegex,
		},
		{
			description: "IsLinkPresent: invalid task link 2",
			msg:         "https://dev.azure.com/abc/xyz/_workitems/edit",
			regex:       constants.TaskLinkRegex,
		},
		{
			description: "IsLinkPresent: invalid task link 3",
			msg:         "https://dev.azure.com/xyz/_workitems/edit/1",
			regex:       constants.TaskLinkRegex,
		},
		{
			description: "IsLinkPresent: invalid task link 4",
			msg:         "http://dev.azure/abc/xyz/items/it/1",
			regex:       constants.TaskLinkRegex,
		},
		{
			description:  "IsLinkPresent: valid pull request link 1",
			msg:          "https://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1/",
			expectedData: []string{"https:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			isValid:      true,
			expectedLink: "https://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1",
			regex:        constants.PullRequestLinkRegex,
		},
		{
			description:  "IsLinkPresent: valid pull request link 2",
			msg:          "https://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1",
			expectedData: []string{"https:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			isValid:      true,
			expectedLink: "https://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1",
			regex:        constants.PullRequestLinkRegex,
		},
		{
			description:  "IsLinkPresent: valid pull request link 3",
			msg:          "http://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1",
			expectedData: []string{"http:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			isValid:      true,
			expectedLink: "http://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1",
			regex:        constants.PullRequestLinkRegex,
		},
		{
			description:  "IsLinkPresent: valid pull request link 4",
			msg:          "\n\nhttp://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1   mock-text",
			expectedData: []string{"http:", "", "dev.azure.com", "abc", "xyz", "_git", "xyz", "pullrequest", "1"},
			isValid:      true,
			expectedLink: "http://dev.azure.com/abc/xyz/_git/xyz/pullrequest/1",
			regex:        constants.PullRequestLinkRegex,
		},
		{
			description: "IsLinkPresent: invalid pull request link 1",
			msg:         "https://abc/xyz/_git/pullrequest/1",
			regex:       constants.PullRequestLinkRegex,
		},
		{
			description: "IsLinkPresent: invalid pull request link 2",
			msg:         "https://dev.azure.com/abc/xyz/_git/pullrequest",
			regex:       constants.PullRequestLinkRegex,
		},
		{
			description: "IsLinkPresent: invalid pull request link 3",
			msg:         "https://dev.azure.com/xyz/_git/xyz/pullrequest/1",
			regex:       constants.PullRequestLinkRegex,
		},
		{
			description: "IsLinkPresent: invalid pull request link 4",
			msg:         "http://dev.azure/abc/xyz/pull/it/1",
			regex:       constants.PullRequestLinkRegex,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			data, link, isValid := IsLinkPresent(testCase.msg, testCase.regex)
			assert.Equal(t, testCase.expectedData, data)
			assert.Equal(t, testCase.isValid, isValid)
			assert.Equal(t, link, testCase.expectedLink)
		})
	}
}
