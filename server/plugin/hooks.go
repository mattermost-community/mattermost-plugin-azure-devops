package plugin

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/config"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/store"
)

// Invoked when configuration changes may have been made.
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

	oldEncryptionSecret := p.getConfiguration().EncryptionSecret
	mattermostSiteURL := p.API.GetConfig().ServiceSettings.SiteURL
	if mattermostSiteURL == nil {
		return errors.New("plugin requires Mattermost Site URL to be set")
	}
	configuration.MattermostSiteURL = *mattermostSiteURL
	p.setConfiguration(configuration)

	if oldEncryptionSecret != "" && oldEncryptionSecret != p.getConfiguration().EncryptionSecret {
		if err := p.Store.DeleteUserTokenOnEncryptionSecretChange(); err != nil {
			p.API.LogError("Error in deleting Users.", "Error", err.Error())
			return err
		}
	}

	return nil
}

// Invoked when the plugin is activated
func (p *Plugin) OnActivate() error {
	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	if err := p.initBotUser(); err != nil {
		return err
	}

	command, err := p.getCommand()
	if err != nil {
		return errors.Wrap(err, "failed to get command")
	}

	if err = p.API.RegisterCommand(command); err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	p.Store = store.NewStore(p.API)
	p.router = p.InitAPI()
	p.InitRoutes()
	return nil
}
