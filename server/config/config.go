package config

import (
	"errors"
	"strings"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
)

// Configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration, as well as values computed from the configuration. Any public fields will be
// deserialized from the Mattermost server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// configuration can change at any time, access to the configuration must be synchronized. The
// strategy used in this plugin is to guard a pointer to the configuration, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
//
// If you add non-reference types to your configuration struct, be sure to rewrite Clone as a deep
// copy appropriate for your types.
type Configuration struct {
	AzureDevopsAPIBaseURL        string `json:"azureDevopsAPIBaseURL"`
	AzureDevopsOAuthAppID        string `json:"azureDevopsOAuthAppID"`
	AzureDevopsOAuthClientSecret string `json:"azureDevopsOAuthClientSecret"`
	EncryptionSecret             string `json:"EncryptionSecret"`
	WebhookSecret                string `json:"WebhookSecret"`
	MattermostSiteURL            string
}

// Clone shallow copies the configuration. Your implementation may require a deep copy if
// your configuration has reference types.
func (c *Configuration) Clone() *Configuration {
	var clone = *c
	return &clone
}

// ProcessConfiguration used for post-processing on the configuration.
func (c *Configuration) ProcessConfiguration() error {
	c.AzureDevopsAPIBaseURL = strings.TrimRight(strings.TrimSpace(c.AzureDevopsAPIBaseURL), "/")
	c.AzureDevopsOAuthAppID = strings.TrimSpace(c.AzureDevopsOAuthAppID)
	c.AzureDevopsOAuthClientSecret = strings.TrimSpace(c.AzureDevopsOAuthClientSecret)
	c.EncryptionSecret = strings.TrimSpace(c.EncryptionSecret)
	c.WebhookSecret = strings.TrimSpace(c.WebhookSecret)

	return nil
}

// Used for config validations.
func (c *Configuration) IsValid() error {
	if c.AzureDevopsAPIBaseURL == "" {
		return errors.New(constants.EmptyAzureDevopsAPIBaseURLError)
	}
	if c.AzureDevopsOAuthAppID == "" {
		return errors.New(constants.EmptyAzureDevopsOAuthAppIDError)
	}
	if c.AzureDevopsOAuthClientSecret == "" {
		return errors.New(constants.EmptyAzureDevopsOAuthClientSecretError)
	}
	if c.EncryptionSecret == "" {
		return errors.New(constants.EmptyEncryptionSecretError)
	}
	if c.WebhookSecret == "" {
		return errors.New(constants.EmptyWebhookSecretError)
	}

	return nil
}
