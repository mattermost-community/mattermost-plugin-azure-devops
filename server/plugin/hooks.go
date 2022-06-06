package plugin

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/api"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/command"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/config"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/utils"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// "OnConfigurationChange" is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(config.Configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	if err := configuration.ProcessConfiguration(); err != nil {
		p.API.LogError("Error in processing configuration.", "Error", err.Error())
		return err
	}

	if err := configuration.IsValid(); err != nil {
		p.API.LogError("Error in validating configuration.", "Error", err.Error())
		return err
	}

	p.setConfiguration(configuration)

	return nil
}

// "OnActivate" is invoked when the plugin is activated
func (p *Plugin) OnActivate() error {
	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	if err := p.registerCommand(); err != nil {
		p.API.LogError(err.Error())
		return err
	}

	p.router = p.InitAPI()
	return nil
}

// "InitAPI" initializes the plugin REST API
func (p *Plugin) InitAPI() *mux.Router {

	r := mux.NewRouter()
	r.Use(p.WithRecovery)

	s := r.PathPrefix(constants.API_PREFIX).Subrouter()
	api.InitRoutes(s)
	api.HandleStaticFiles(p.API, r)

	// 404 handler
	r.Handle(constants.WILD_ROUTE, http.NotFoundHandler())
	return r
}

// Handles executing a slash command
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split, argErr := utils.SplitArgs(args.Command)
	if argErr != nil {
		return utils.SendEphemeralCommandResponse(argErr.Error())
	}

	cmdName := split[0][1:]
	var params []string

	if len(split) > 1 {
		params = split[1:]
	}

	cmd := command.GetCommand("")
	if cmd.Trigger != cmdName {
		return utils.SendEphemeralCommandResponse("Unknown command: [" + cmdName + "] encountered")
	}

	p.API.LogDebug("Executing command: " + cmdName + "]")
	cmdContext := command.NewContext(args, c, p.API, p.Helpers)
	return command.AzureDevopsCommandHandler.Handle(cmdContext, params...)
}
