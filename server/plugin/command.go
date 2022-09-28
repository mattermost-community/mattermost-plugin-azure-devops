package plugin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-api/experimental/command"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

type HandlerFunc func(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError)

type Handler struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

var azureDevopsCommandHandler = Handler{
	handlers: map[string]HandlerFunc{
		"help":       azureDevopsHelpCommand,
		"connect":    azureDevopsConnectCommand,
		"disconnect": azureDevopsDisconnectCommand,
		"link":       azureDevopsAccountConnectionCheck,
		"boards":     azureDevopsAccountConnectionCheck,
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
	azureDevops := model.NewAutocompleteData(constants.CommandTriggerName, "[command]", "Available commands: help, connect, disconnect, create, link, subscribe, subscriptions, unsubscribe")

	help := model.NewAutocompleteData("help", "", fmt.Sprintf("Show %s slash command help", constants.CommandTriggerName))
	azureDevops.AddCommand(help)

	connect := model.NewAutocompleteData("connect", "", "Connect to your Azure DevOps account")
	azureDevops.AddCommand(connect)

	disconnect := model.NewAutocompleteData("disconnect", "", "Disconnect your Azure DevOps account")
	azureDevops.AddCommand(disconnect)

	link := model.NewAutocompleteData("link", "[projectURL]", "Link a project")
	azureDevops.AddCommand(link)

	create := model.NewAutocompleteData("boards create [title] [description]", "", "Create a new task")
	azureDevops.AddCommand(create)

	subscribe := model.NewAutocompleteData("boards subscribe", "", "Add a boards subscription")
	azureDevops.AddCommand(subscribe)

	subscriptions := model.NewAutocompleteData(fmt.Sprintf("boards subscriptions [%s or %s] [%s or %s]", constants.FilterCreatedByMe, constants.FilterCreatedByAnyone, constants.FilterCurrentChannel, constants.FilterAllChannel), "", "View board's subscriptions")
	azureDevops.AddCommand(subscriptions)

	unsubscribe := model.NewAutocompleteData("boards unsubscribe [subscription id]", "", "Unsubscribe a board subscription")
	azureDevops.AddCommand(unsubscribe)

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
		AutoCompleteDesc:     "Available commands: help",
		AutoCompleteHint:     "[command]",
		AutocompleteData:     p.getAutoCompleteData(),
		AutocompleteIconData: iconData,
	}, nil
}

func azureDevopsAccountConnectionCheck(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	if isConnected := p.UserAlreadyConnected(commandArgs.UserId); !isConnected {
		return p.sendEphemeralPostForCommand(commandArgs, p.getConnectAccountFirstMessage())
	}

	if len(args) > 0 {
		switch {
		case args[0] == "subscriptions":
			return azureDevopsListSubscriptionsCommand(p, c, commandArgs, args...)
		case args[0] == "unsubscribe":
			return azureDevopsUnsubscribeCommand(p, c, commandArgs, args...)
		}
	}

	return &model.CommandResponse{}, nil
}

func azureDevopsUnsubscribeCommand(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	if len(args) < 2 {
		return p.sendEphemeralPostForCommand(commandArgs, "Subscription ID is not provided")
	}

	subscriptionList, err := p.Store.GetAllSubscriptions("")
	if err != nil {
		p.API.LogError(constants.FetchSubscriptionListError, "Error", err.Error())
		return p.sendEphemeralPostForCommand(commandArgs, constants.GenericErrorMessage)
	}

	for _, subscription := range subscriptionList {
		if subscription.SubscriptionID == args[1] {
			if _, err := p.sendEphemeralPostForCommand(commandArgs, fmt.Sprintf("Boards subscription with ID: %q is being deleted", args[1])); err != nil {
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

			return p.sendEphemeralPostForCommand(commandArgs, fmt.Sprintf("Boards subscription with ID: %q is successfully deleted", args[1]))
		}
	}

	return p.sendEphemeralPostForCommand(commandArgs, fmt.Sprintf("Boards subscription with ID: %q does not exist", args[1]))
}

func azureDevopsListSubscriptionsCommand(p *Plugin, c *plugin.Context, commandArgs *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	if !(len(args) == 3 && (args[1] == constants.FilterCreatedByMe || args[1] == constants.FilterCreatedByAnyone || args[2] == constants.FilterAllChannel || args[2] == constants.FilterCurrentChannel)) {
		executeDefault(p, c, commandArgs, args...)
		return &model.CommandResponse{}, nil
	}

	subscriptionList, err := p.Store.GetAllSubscriptions("")
	if err != nil {
		p.API.LogError(constants.FetchSubscriptionListError, "Error", err.Error())
		return p.sendEphemeralPostForCommand(commandArgs, constants.GenericErrorMessage)
	}

	showForChannelID := commandArgs.ChannelId
	if args[2] == constants.FilterAllChannel {
		showForChannelID = ""
	}
	return p.sendEphemeralPostForCommand(commandArgs, p.ParseSubscriptionsToCommandResponse(subscriptionList, showForChannelID, args[1], commandArgs.UserId))
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
	return p.sendEphemeralPostForCommand(commandArgs, constants.InvalidCommand)
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
