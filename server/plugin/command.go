package plugin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/mattermost/mattermost-plugin-api/experimental/command"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
)

type HandlerFunc func(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError)

type Handler struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

var azureDevopsCommandHandler = Handler{
	handlers: map[string]HandlerFunc{
		constants.CommandHelp:       azureDevopsHelpCommand,
		constants.CommandConnect:    azureDevopsConnectCommand,
		constants.CommandDisconnect: azureDevopsDisconnectCommand,
		constants.CommandLink:       azureDevopsAccountConnectionCheck,
		constants.CommandBoards:     azureDevopsBoardsCommand,
		constants.CommandRepos:      azureDevopsReposCommand,
	},
	defaultHandler: executeDefault,
}

// Handle function calls the respective handlers of the commands.
// It checks whether any HandlerFunc is present for the given command by checking in the "azureDevopsCommandHandler".
// If the command is present, it calls its handler function, else calls the default handler.
func (ch *Handler) Handle(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	for arg := len(args); arg > 0; arg-- {
		handler := ch.handlers[strings.Join(args[:arg], "/")]
		if handler != nil {
			return handler(p, c, commandArgs, args[arg:]...)
		}
	}
	return ch.defaultHandler(p, c, commandArgs, args...)
}

func (p *Plugin) getAutoCompleteData() *model.AutocompleteData {
	azureDevops := model.NewAutocompleteData(constants.CommandTriggerName, "[command]", fmt.Sprintf("Available commands: %s, %s, %s, %s, %s", constants.CommandHelp, constants.CommandConnect, constants.CommandDisconnect, constants.CommandLink, constants.CommandBoards))

	help := model.NewAutocompleteData(constants.CommandHelp, "", fmt.Sprintf("Show %s slash command help", constants.CommandTriggerName))
	azureDevops.AddCommand(help)

	connect := model.NewAutocompleteData(constants.CommandConnect, "", "Connect to your Azure DevOps account")
	azureDevops.AddCommand(connect)

	disconnect := model.NewAutocompleteData(constants.CommandDisconnect, "", "Disconnect your Azure DevOps account")
	azureDevops.AddCommand(disconnect)

	link := model.NewAutocompleteData(constants.CommandLink, "", "Link a project")
	link.AddTextArgument("URL of the project to be linked", "[projectURL]", "")
	azureDevops.AddCommand(link)

	subscription := model.NewAutocompleteData(constants.CommandSubscription, "", "Add/list/delete subscriptions")
	subscriptionAdd := model.NewAutocompleteData(constants.CommandAdd, "", "Add a new subscription")
	subscriptionList := model.NewAutocompleteData(constants.CommandList, "", "List subscriptions")
	subscriptionDelete := model.NewAutocompleteData(constants.CommandDelete, "", "Delete a subscription")
	subscriptionDelete.AddTextArgument("ID of the subscription to be deleted", "[subscription id]", "")
	subscriptionCreatedByMe := model.NewAutocompleteData(constants.FilterCreatedByMe, "", "Created By Me")
	subscriptionShowForAllChannels := model.NewAutocompleteData(constants.FilterAllChannels, "", "Show for all channels or You can leave this argument to show for the current channel only")
	subscriptionCreatedByMe.AddCommand(subscriptionShowForAllChannels)
	subscriptionList.AddCommand(subscriptionCreatedByMe)
	subscriptionCreatedByAnyone := model.NewAutocompleteData(constants.FilterCreatedByAnyone, "", "Created By Anyone")
	subscriptionCreatedByAnyone.AddCommand(subscriptionShowForAllChannels)
	subscriptionList.AddCommand(subscriptionCreatedByAnyone)
	subscription.AddCommand(subscriptionAdd)
	subscription.AddCommand(subscriptionList)
	subscription.AddCommand(subscriptionDelete)

	boards := model.NewAutocompleteData(constants.CommandBoards, "", "Create a new work-item or add/list/delete board subscriptions")
	create := model.NewAutocompleteData(constants.CommandCreate, "", "Create a new work-item")
	create.AddTextArgument("Title", "[title]", "")
	create.AddTextArgument("Description", "[description]", "")
	boards.AddCommand(create)
	boards.AddCommand(subscription)
	azureDevops.AddCommand(boards)

	repos := model.NewAutocompleteData(constants.CommandRepos, "", "Add/list/delete repo subscriptions")
	repos.AddCommand(subscription)
	azureDevops.AddCommand(repos)

	return azureDevops
}

func (p *Plugin) getCommand() (*model.Command, error) {
	iconData, err := command.GetIconData(p.API, "assets/azurebot.svg")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Azure DevOps icon")
	}

	return &model.Command{
		Trigger:              constants.CommandTriggerName,
		AutoComplete:         true,
		AutoCompleteDesc:     fmt.Sprintf("Available commands: %s, %s, %s, %s, %s", constants.CommandHelp, constants.CommandConnect, constants.CommandDisconnect, constants.CommandLink, constants.CommandBoards),
		AutoCompleteHint:     "[command]",
		AutocompleteData:     p.getAutoCompleteData(),
		AutocompleteIconData: iconData,
	}, nil
}

func azureDevopsAccountConnectionCheck(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	if isConnected := p.UserAlreadyConnected(commandArgs.UserId); !isConnected {
		return p.sendEphemeralPostForCommand(commandArgs, p.getConnectAccountFirstMessage())
	}
	return &model.CommandResponse{}, nil
}

func azureDevopsBoardsCommand(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	// Check if user's Azure DevOps account is connected
	if isConnected := p.UserAlreadyConnected(commandArgs.UserId); !isConnected {
		return p.sendEphemeralPostForCommand(commandArgs, p.getConnectAccountFirstMessage())
	}

	// Validate commands and their arguments
	switch {
	case len(args) >= 1 && args[0] == constants.CommandCreate:
		return &model.CommandResponse{}, nil
		// For "subscription" command there must be at least 2 arguments
	case len(args) >= 2 && args[0] == constants.CommandSubscription:
		switch args[1] {
		case constants.CommandList:
			// For "list" command there must be at least 3 arguments
			if len(args) >= 3 && (args[2] == constants.FilterCreatedByMe || args[2] == constants.FilterCreatedByAnyone) {
				return azureDevopsListSubscriptionsCommand(p, c, commandArgs, constants.CommandBoards, args...)
			}
		case constants.CommandDelete:
			return azureDevopsDeleteCommand(p, c, commandArgs, constants.CommandBoards, args...)
		case constants.CommandAdd:
			return &model.CommandResponse{}, nil
		}
	}

	return executeDefault(p, c, commandArgs, args...)
}

func azureDevopsReposCommand(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	// Check if user's Azure DevOps account is connected
	if isConnected := p.UserAlreadyConnected(commandArgs.UserId); !isConnected {
		return p.sendEphemeralPostForCommand(commandArgs, p.getConnectAccountFirstMessage())
	}

	// Validate commands and their arguments
	// For "subscription" command there must be at least 2 arguments
	if len(args) >= 2 && args[0] == constants.CommandSubscription {
		switch args[1] {
		case constants.CommandList:
			// For "list" command there must be at least 3 arguments
			if len(args) >= 3 && (args[2] == constants.FilterCreatedByMe || args[2] == constants.FilterCreatedByAnyone) {
				return azureDevopsListSubscriptionsCommand(p, c, commandArgs, constants.CommandRepos, args...)
			}
		case constants.CommandDelete:
			return azureDevopsDeleteCommand(p, c, commandArgs, constants.CommandRepos, args...)
		case constants.CommandAdd:
			return &model.CommandResponse{}, nil
		}
	}

	return executeDefault(p, c, commandArgs, args...)
}

func azureDevopsDeleteCommand(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, command string, args ...string) (*model.CommandResponse, *model.AppError) {
	if len(args) < 3 {
		return p.sendEphemeralPostForCommand(commandArgs, "Subscription ID is not provided")
	}

	subscriptionList, err := p.Store.GetAllSubscriptions("")
	if err != nil {
		p.API.LogError(constants.FetchSubscriptionListError, "Error", err.Error())
		return p.sendEphemeralPostForCommand(commandArgs, constants.GenericErrorMessage)
	}

	subscriptionIDToBeDeleted := args[2]
	for _, subscription := range subscriptionList {
		if subscription.SubscriptionID == subscriptionIDToBeDeleted && subscription.ServiceType == command {
			if _, err := p.sendEphemeralPostForCommand(commandArgs, fmt.Sprintf("%s subscription with ID: %q is being deleted", cases.Title(language.Und).String(command), subscriptionIDToBeDeleted)); err != nil {
				p.API.LogError("Error in sending ephemeral post", "Error", err.Error())
			}

			if statusCode, err := p.Client.DeleteSubscription(subscription.OrganizationName, subscription.SubscriptionID, commandArgs.UserId); err != nil {
				if statusCode == http.StatusForbidden {
					return p.sendEphemeralPostForCommand(commandArgs, constants.ErrorAdminAccess)
				}
				p.API.LogError("Error in deleting subscription", "Error", err.Error())
				return p.sendEphemeralPostForCommand(commandArgs, constants.GenericErrorMessage)
			}

			if deleteErr := p.Store.DeleteSubscription(subscription); deleteErr != nil {
				p.API.LogError("Error in deleting subscription", "Error", deleteErr.Error())
				return p.sendEphemeralPostForCommand(commandArgs, constants.GenericErrorMessage)
			}

			p.API.PublishWebSocketEvent(
				constants.WSEventSubscriptionDeleted,
				nil,
				&model.WebsocketBroadcast{UserId: commandArgs.UserId},
			)

			return p.sendEphemeralPostForCommand(commandArgs, fmt.Sprintf("%s subscription with ID: %q is successfully deleted", cases.Title(language.Und).String(command), subscriptionIDToBeDeleted))
		}
	}

	return p.sendEphemeralPostForCommand(commandArgs, fmt.Sprintf("%s subscription with ID: %q does not exist", cases.Title(language.Und).String(command), subscriptionIDToBeDeleted))
}

func azureDevopsListSubscriptionsCommand(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, command string, args ...string) (*model.CommandResponse, *model.AppError) {
	// If 4th argument is present then it must be "all_channels"
	if len(args) >= 4 && args[3] != constants.FilterAllChannels {
		return executeDefault(p, c, commandArgs, args...)
	}

	subscriptionList, err := p.Store.GetAllSubscriptions("")
	if err != nil {
		p.API.LogError(constants.FetchSubscriptionListError, "Error", err.Error())
		return p.sendEphemeralPostForCommand(commandArgs, constants.GenericErrorMessage)
	}

	showForChannelID := commandArgs.ChannelId
	if len(args) >= 4 && args[3] == constants.FilterAllChannels {
		showForChannelID = ""
	}
	return p.sendEphemeralPostForCommand(commandArgs, p.ParseSubscriptionsToCommandResponse(subscriptionList, showForChannelID, args[2], commandArgs.UserId, command, commandArgs.TeamId))
}

func azureDevopsHelpCommand(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	return p.sendEphemeralPostForCommand(commandArgs, constants.HelpText)
}

func azureDevopsConnectCommand(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	message := fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect)
	if isConnected := p.UserAlreadyConnected(commandArgs.UserId); isConnected {
		message = constants.UserAlreadyConnected
	}
	return p.sendEphemeralPostForCommand(commandArgs, message)
}

func azureDevopsDisconnectCommand(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	message := constants.UserDisconnected
	if isConnected := p.UserAlreadyConnected(commandArgs.UserId); !isConnected {
		message = p.getConnectAccountFirstMessage()
	} else {
		if isDeleted, err := p.Store.DeleteUser(commandArgs.UserId); !isDeleted {
			if err != nil {
				p.API.LogError(constants.UnableToDisconnectUser, "Error", err.Error())
			}
			message = constants.GenericErrorMessage
		}

		p.API.PublishWebSocketEvent(
			constants.WSEventDisconnect,
			nil,
			&model.WebsocketBroadcast{UserId: commandArgs.UserId},
		)
	}
	return p.sendEphemeralPostForCommand(commandArgs, message)
}

func executeDefault(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	out := constants.InvalidCommand + "\n\n" + constants.HelpText

	return p.sendEphemeralPostForCommand(commandArgs, out)
}

// Handles executing a slash command
func (p *Plugin) ExecuteCommand(c *plugin.Context, commandArgs *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	args := strings.Fields(commandArgs.Command)

	if len(args) == 0 || args[0] != fmt.Sprintf("/%s", constants.CommandTriggerName) {
		commandName := args[0][1:]
		return p.sendEphemeralPostForCommand(commandArgs, fmt.Sprintf("unknown command %s\n%s", commandName, constants.HelpText))
	}

	return azureDevopsCommandHandler.Handle(p, c, commandArgs, args[1:]...)
}
