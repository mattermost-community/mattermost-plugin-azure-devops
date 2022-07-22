package config

import (
	"errors"
	"strings"
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
	// TODO: Below configs are not final they are used as placeholder here
	AzureDevopsAPIBaseURL            string `json:"azureDevopsAPIBaseURL"`
	AzureDevopsOAuthAppID            string `json:"azureDevopsOAuthAppID"`
	AzureDevopsOAuthClientSecret     string `jso:"azureDevopsOAuthClientSecret"`
	AzureDevopsOAuthAuthorizationURL string `json:"azureDevopsOAuthAuthorizationURL"`
	AzureDevopsOAuthTokenURL         string `json:"azureDevopsOAuthTokenURL"`
	AzureDevopsOAuthCallbackURL      string `json:"azureDevopsOAuthCallbackURL"`
	EncryptionSecret                 string `json:"EncryptionSecret"`
	MattermostSiteURL                string
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
	c.AzureDevopsOAuthAuthorizationURL = strings.TrimRight(strings.TrimSpace(c.AzureDevopsOAuthAuthorizationURL), "/")
	c.AzureDevopsOAuthTokenURL = strings.TrimRight(strings.TrimSpace(c.AzureDevopsOAuthTokenURL), "/")
	c.EncryptionSecret = strings.TrimSpace(c.EncryptionSecret)

	return nil
}

// Used for config validations.
func (c *Configuration) IsValid() error {
	if c.AzureDevopsAPIBaseURL == "" {
		return errors.New("azure devops API base URL should not be empty")
	}
	if c.AzureDevopsOAuthAppID == "" {
		return errors.New("azure devops OAuth app id should not be empty")
	}
	if c.AzureDevopsOAuthClientSecret == "" {
		return errors.New("azure devops OAuth client secret should not be empty")
	}
	if c.AzureDevopsOAuthAuthorizationURL == "" {
		return errors.New("azure devops OAuth authorization URL should not be empty")
	}
	if c.AzureDevopsOAuthTokenURL == "" {
		return errors.New("azure devops OAuth token URL should not be empty")
	}
	if c.AzureDevopsOAuthCallbackURL == "" {
		return errors.New("azure devops OAuth callback URL should not be empty")
	}
	if c.EncryptionSecret == "" {
		return errors.New("encryption secret should not be empty")
	}

	return nil
}
