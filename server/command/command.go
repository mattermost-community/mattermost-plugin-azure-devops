package command

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/utils"
)

// Context includes the context in which the slash command is executed and allows access to
// plugin API, helpers and services
type Context struct {
	*model.CommandArgs
	context *plugin.Context
	api     plugin.API
	helpers plugin.Helpers
}

func NewContext(args *model.CommandArgs, context *plugin.Context, api plugin.API, helpers plugin.Helpers) *Context {
	return &Context{
		args,
		context,
		api,
		helpers,
	}
}

type HandlerFunc func(context *Context, args ...string) (*model.CommandResponse, *model.AppError)

type Handler struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

func (ch Handler) Handle(context *Context, args ...string) (*model.CommandResponse, *model.AppError) {
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(context, args[n:]...)
		}
	}
	return ch.defaultHandler(context, args...)
}

const (
	invalidCommand = "Invalid command parameters. Please use `/azuredevops help` for more information."
	helpText       = "###### Azure DevOps - Slash Command Help\n\n" // TODO: Add proper text here
)

func GetCommand(iconData string) *model.Command {
	return &model.Command{
		Trigger:              "azuredevops",
		DisplayName:          "Azure DevOps",
		AutoComplete:         true,
		AutoCompleteDesc:     "Available commands: help",
		AutoCompleteHint:     "[command]",
		AutocompleteData:     getAutoCompleteData(),
		AutocompleteIconData: iconData,
	}
}

func getAutoCompleteData() *model.AutocompleteData {
	azuredevops := model.NewAutocompleteData("azuredevops", "[command]", "Available commands: help.")

	help := model.NewAutocompleteData("help", "", "Show azuredevops slash command help")
	azuredevops.AddCommand(help)

	return azuredevops
}

var AzureDevopsCommandHandler = Handler{
	handlers: map[string]HandlerFunc{
		"help": azureDevopsHelpCommand,
	},
	defaultHandler: func(context *Context, args ...string) (*model.CommandResponse, *model.AppError) {
		return utils.SendEphemeralCommandResponse(invalidCommand)
	},
}

func azureDevopsHelpCommand(ctx *Context, args ...string) (*model.CommandResponse, *model.AppError) {
	return utils.SendEphemeralCommandResponse(helpText)
}
