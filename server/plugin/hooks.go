package plugin

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/config"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/routes"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// OnConfigurationChange is invoked when configuration changes may have been made.
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

// OnActivate is invoked when the plugin is activated
func (p *Plugin) OnActivate() error {
	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	p.router = p.InitAPI()
	return nil
}

// InitAPI initializes the REST API
func (p *Plugin) InitAPI() *mux.Router {

	r := mux.NewRouter()
	r.Use(p.WithRecovery)

	s := r.PathPrefix(constants.API_PREFIX).Subrouter()
	routes.InitRoutes(s)

	// 404 handler
	r.Handle(constants.WILD_ROUTE, http.NotFoundHandler())
	return r
}
