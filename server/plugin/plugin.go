package plugin

import (
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/config"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/store"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	Client Client

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *config.Configuration
	router        *mux.Router
	Store         store.KVStore

	// user ID of the bot account
	botUserID string
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *config.Configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &config.Configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing configuration. This almost
// certainly means that the configuration was modified without being cloned and may result in
// an unsafe access.
func (p *Plugin) setConfiguration(configuration *config.Configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		// Ignore assignment if the configuration struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

// Initializes a bot user
func (p *Plugin) initBotUser() error {
	botID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    constants.BotUsername,
		DisplayName: constants.BotDisplayName,
		Description: constants.BotDescription,
	}, plugin.ProfileImagePath(filepath.Join("assets", "azurebot.png")))
	if err != nil {
		return errors.Wrap(err, "cannot create bot")
	}

	p.botUserID = botID
	return nil
}

// ServeHTTP demonstrates a plugin that handles HTTP requests.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	// Check if message is a work item link.
	if taskData, isValid := isValidTaskLink(post.Message); isValid {
		newPost, msg := p.PostTaskPreview(taskData, post.UserId, post.ChannelId)
		return newPost, msg
	}
	return nil, ""
}

// Function to validate the work item link.
func isValidTaskLink(msg string) ([]string, bool) {
	trimmedLink := strings.TrimRight(msg, "/")
	data := strings.Split(trimmedLink, "/")
	if len(data) != 8 {
		return nil, false
	}
	if (data[0] != constants.HTTPS && data[0] != constants.HTTP) || data[2] != constants.AzureDevopsBaseURL || data[5] != constants.Workitems || data[6] != constants.Edit {
		return nil, false
	}
	return data, true
}
