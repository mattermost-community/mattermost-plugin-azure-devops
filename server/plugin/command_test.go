package plugin

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mattermost/mattermost-plugin-azure-devops/mocks"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/testutils"
)

func TestExecuteCommand(t *testing.T) {
	monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, mockedClient)
	for _, testCase := range []struct {
		description                   string
		commandArgs                   *model.CommandArgs
		ephemeralMessage              string
		isConnected                   bool
		patchAPICalls                 func()
		isListCommand                 bool
		isDeleteCommand               bool
		serviceType                   string
		deleteSubscriptionClientError error
		deleteSubscriptionStoreError  error
		getAllSubscriptionError       error
	}{
		{
			description:      "ExecuteCommand: empty command",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops"},
			ephemeralMessage: constants.InvalidCommand + constants.HelpText,
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
			commandArgs:      &model.CommandArgs{Command: "/azuredevops disconnect", UserId: testutils.MockMattermostUserID},
			isConnected:      true,
			ephemeralMessage: constants.UserDisconnected,
		},
		{
			description:      "ExecuteCommand: boards command when user is not connected",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops boards wrong [title] [description]"},
			ephemeralMessage: fmt.Sprintf(constants.ConnectAccountFirst, fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect)),
		},
		{
			description:      "ExecuteCommand: invalid boards command",
			isConnected:      true,
			commandArgs:      &model.CommandArgs{Command: "/azuredevops boards wrong [title] [description]"},
			ephemeralMessage: constants.InvalidCommand + constants.HelpText,
		},
		{
			description:      "ExecuteCommand: boards create command",
			isConnected:      true,
			commandArgs:      &model.CommandArgs{Command: "/azuredevops boards workitem create [title] [description]"},
			ephemeralMessage: fmt.Sprintf(constants.ConnectAccountFirst, fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect)),
		},
		{
			description: "ExecuteCommand: boards add subscription command",
			isConnected: true,
			commandArgs: &model.CommandArgs{Command: "/azuredevops boards subscription add"},
		},
		{
			description:     "ExecuteCommand: boards delete subscription command",
			isConnected:     true,
			commandArgs:     &model.CommandArgs{Command: "/azuredevops boards subscription delete mockSubscriptionID"},
			isDeleteCommand: true,
			serviceType:     "boards",
		},
		{
			description:                  "ExecuteCommand: failed to delete subscription from store",
			isConnected:                  true,
			commandArgs:                  &model.CommandArgs{Command: "/azuredevops boards subscription delete mockSubscriptionID"},
			isDeleteCommand:              true,
			serviceType:                  "boards",
			deleteSubscriptionStoreError: errors.New("failed to delete subscription from store"),
		},
		{
			description:             "ExecuteCommand: failed to get all subscriptions while deleting a subscription",
			isConnected:             true,
			commandArgs:             &model.CommandArgs{Command: "/azuredevops boards subscription delete mockSubscriptionID"},
			isListCommand:           true,
			serviceType:             "boards",
			getAllSubscriptionError: errors.New("failed to get all subscriptions while deleting a subscription"),
			ephemeralMessage:        constants.GenericErrorMessage,
		},
		{
			description:      "ExecuteCommand: boards list subscription command",
			isConnected:      true,
			commandArgs:      &model.CommandArgs{Command: "/azuredevops boards subscription list me"},
			isListCommand:    true,
			ephemeralMessage: "mockSubscriptionList",
		},
		{
			description:             "ExecuteCommand: failed to get all the subscriptions",
			isConnected:             true,
			commandArgs:             &model.CommandArgs{Command: "/azuredevops boards subscription list me"},
			isListCommand:           true,
			getAllSubscriptionError: errors.New("failed to get all subscriptions"),
			ephemeralMessage:        constants.GenericErrorMessage,
		},
		{
			description:      "ExecuteCommand: repos command when user is not connected",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops repos wrong [title] [description]"},
			ephemeralMessage: fmt.Sprintf(constants.ConnectAccountFirst, fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect)),
		},
		{
			description:      "ExecuteCommand: invalid repos command",
			isConnected:      true,
			commandArgs:      &model.CommandArgs{Command: "/azuredevops repos wrong [title] [description]"},
			ephemeralMessage: constants.InvalidCommand + constants.HelpText,
		},
		{
			description: "ExecuteCommand: repos add subscription command",
			isConnected: true,
			commandArgs: &model.CommandArgs{Command: "/azuredevops repos subscription add"},
		},
		{
			description:     "ExecuteCommand: repos delete subscription command",
			isConnected:     true,
			commandArgs:     &model.CommandArgs{Command: "/azuredevops repos subscription delete mockSubscriptionID"},
			isDeleteCommand: true,
			serviceType:     "repos",
		},
		{
			description:      "ExecuteCommand: repos list subscription command",
			isConnected:      true,
			commandArgs:      &model.CommandArgs{Command: "/azuredevops repos subscription list me"},
			isListCommand:    true,
			ephemeralMessage: "mockSubscriptionList",
		},
		{
			description:      "ExecuteCommand: pipelines command when user is not connected",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops pipelines wrong [title] [description]"},
			ephemeralMessage: fmt.Sprintf(constants.ConnectAccountFirst, fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect)),
		},
		{
			description: "ExecuteCommand: pipelines add subscription command",
			isConnected: true,
			commandArgs: &model.CommandArgs{Command: "/azuredevops pipelines subscription add"},
		},
		{
			description:      "ExecuteCommand: invalid pipelines command",
			isConnected:      true,
			commandArgs:      &model.CommandArgs{Command: "/azuredevops pipelines wrong [title] [description]"},
			ephemeralMessage: constants.InvalidCommand + constants.HelpText,
		},
		{
			description:     "ExecuteCommand: pipelines delete subscription command",
			isConnected:     true,
			commandArgs:     &model.CommandArgs{Command: "/azuredevops pipelines subscription delete mockSubscriptionID"},
			isDeleteCommand: true,
			serviceType:     "pipelines",
		},
		{
			description:      "ExecuteCommand: pipelines list subscription command",
			isConnected:      true,
			commandArgs:      &model.CommandArgs{Command: "/azuredevops pipelines subscription list me"},
			isListCommand:    true,
			ephemeralMessage: "mockSubscriptionList",
		},
		{
			description:      "ExecuteCommand: invalid command",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops abc"},
			ephemeralMessage: constants.InvalidCommand + constants.HelpText,
		},
		{
			description: "ExecuteCommand: link command",
			commandArgs: &model.CommandArgs{Command: "/azuredevops link"},
			isConnected: true,
		},
		{
			description:      "ExecuteCommand: link command when user is not connected",
			commandArgs:      &model.CommandArgs{Command: "/azuredevops link"},
			ephemeralMessage: fmt.Sprintf(constants.ConnectAccountFirst, fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect)),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
				post := args.Get(1).(*model.Post)
				if !testCase.isDeleteCommand {
					assert.Equal(t, testCase.ephemeralMessage, post.Message)
				}
			}).Once().Return(&model.Post{})

			mockAPI.On("GetBundlePath").Return("/test-path", nil)
			mockAPI.On("LogError", testutils.GetMockArgumentsWithType("string", 3)...)
			mockAPI.On("PublishWebSocketEvent", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("*model.WebsocketBroadcast")).Return()

			monkey.PatchInstanceMethod(reflect.TypeOf(p), "UserAlreadyConnected", func(_ *Plugin, _ string) bool {
				return testCase.isConnected
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(p), "ParseSubscriptionsToCommandResponse", func(_ *Plugin, _ []*serializers.SubscriptionDetails, _, _, _, _, _ string) string {
				return "mockSubscriptionList"
			})

			if testCase.isListCommand || testCase.isDeleteCommand {
				mockedStore.EXPECT().GetAllSubscriptions("").Return(testutils.GetSuscriptionDetailsPayload(testutils.MockMattermostUserID, testCase.serviceType, testutils.MockEventType),
					testCase.getAllSubscriptionError)
			}

			if testCase.isDeleteCommand {
				mockedClient.EXPECT().DeleteSubscription(gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusOK, testCase.deleteSubscriptionClientError)
				mockedStore.EXPECT().DeleteSubscription(gomock.Any()).Return(testCase.deleteSubscriptionStoreError)
			}

			if testCase.ephemeralMessage == constants.UserDisconnected {
				mockedStore.EXPECT().DeleteUser(testutils.MockMattermostUserID).Return(true, nil)
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
