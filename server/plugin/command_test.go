package plugin

import (
	"fmt"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/Brightscout/mattermost-plugin-azure-devops/mocks"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecuteCommand(t *testing.T) {
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Store = mockedStore
	for _, testCase := range []struct {
		description      string
		commandArgs      *model.CommandArgs
		ephemeralMessage string
		isConnected      bool
		patchAPICalls    func()
	}{
		{
			description:      "ExecuteCommand: empty command",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops"},
			ephemeralMessage: constants.InvalidCommand,
		},
		{
			description:      "ExecuteCommand: help command",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops help"},
			ephemeralMessage: constants.HelpText,
		},
		{
			description:      "ExecuteCommand: connect command",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops connect"},
			ephemeralMessage: fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect),
		},
		{
			description:      "ExecuteCommand: connect command with user already connected",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops connect"},
			isConnected:      true,
			ephemeralMessage: constants.UserAlreadyConnected,
		},
		{
			description:      "ExecuteCommand: disconnect command with user not connected",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops disconnect"},
			ephemeralMessage: fmt.Sprintf(constants.ConnectAccountFirst, fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect)),
		},
		{
			description:      "ExecuteCommand: disconnect command with user connected",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops disconnect", UserId: "mockUserID"},
			isConnected:      true,
			ephemeralMessage: constants.UserDisconnected,
		},
		{
			description:      "ExecuteCommand: boards create command",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops boards create [title] [description]"},
			ephemeralMessage: fmt.Sprintf(constants.ConnectAccountFirst, fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect)),
		},
		{
			description:      "ExecuteCommand: invalid command",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops abc"},
			ephemeralMessage: constants.InvalidCommand,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
				post := args.Get(1).(*model.Post)
				assert.Equal(t, testCase.ephemeralMessage, post.Message)
			}).Once().Return(&model.Post{})

			mockAPI.On("GetBundlePath").Return("/test-path", nil)
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			mockAPI.On("PublishWebSocketEvent", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("*model.WebsocketBroadcast")).Return(nil)

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "UserAlreadyConnected", func(_ *Plugin, _ string) bool {
				return testCase.isConnected
			})

			if testCase.ephemeralMessage == constants.UserDisconnected {
				mockedStore.EXPECT().DeleteUser("mockUserID").Return(true, nil)
			}

			_, err := p.getCommand()
			assert.NotNil(t, err)

			response := p.getAutoCompleteData()
			assert.NotNil(t, response)

			res, err := p.ExecuteCommand(&plugin.Context{}, testCase.commandArgs)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		})
	}
}
