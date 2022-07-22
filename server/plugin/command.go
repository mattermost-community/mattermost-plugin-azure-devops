package plugin

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-api/experimental/command"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

type HandlerFunc func(p *Plugin, c *plugin.Context, header *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError)

type Handler struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

var azureDevopsCommandHandler = Handler{
	handlers: map[string]HandlerFunc{
		"help": azureDevopsHelpCommand,
	},
	defaultHandler: executeDefault,
}

// TODO: add comments to explain the below code or refactor it
func (ch *Handler) Handle(p *Plugin, c *plugin.Context, header *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(p, c, header, args[n:]...)
		}
	}
	return ch.defaultHandler(p, c, header, args...)
}

func (p *Plugin) getAutoCompleteData() *model.AutocompleteData {
	azureDevops := model.NewAutocompleteData(constants.CommandTriggerName, "[command]", "Available commands: help")

	help := model.NewAutocompleteData("help", "", fmt.Sprintf("Show %s slash command help", constants.CommandTriggerName))
	azureDevops.AddCommand(help)

	return azureDevops
}

func (p *Plugin) getCommand() (*model.Command, error) {
	iconData, err := command.GetIconData(p.API, "assets/azurebot.svg")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Azure Devops icon")
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

func azureDevopsHelpCommand(p *Plugin, c *plugin.Context, header *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	return p.sendEphemeralPost(header, constants.HelpText)
}

func executeDefault(p *Plugin, c *plugin.Context, header *model.CommandArgs, args ...string) (*model.CommandResponse, *model.AppError) {
	return p.sendEphemeralPost(header, constants.InvalidCommand)
}

// Handles executing a slash command
func (p *Plugin) ExecuteCommand(c *plugin.Context, commandArgs *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	args := strings.Fields(commandArgs.Command)

	if len(args) == 0 || args[0] != fmt.Sprintf("/%s", constants.CommandTriggerName) {
		commandName := args[0][1:]
		return p.sendEphemeralPost(commandArgs, fmt.Sprintf("unknown command %s\n%s", commandName, constants.HelpText))
	}

	return azureDevopsCommandHandler.Handle(p, c, commandArgs, args[1:]...)
}
