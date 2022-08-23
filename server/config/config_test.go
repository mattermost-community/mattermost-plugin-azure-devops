package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsValid(t *testing.T) {
	for _, testCase := range []struct {
		description string
		config      *Configuration
		errMsg      string
	}{
		{
			description: "valid configuration",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "mockAzureDevopsAPIBaseURL",
				AzureDevopsOAuthAppID:        "mockAzureDevopsOAuthAppID",
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				EncryptionSecret:             "mockEncryptionSecret",
			},
			errMsg: "",
		},
		{
			description: "invalid configuration with AzureDevopsAPIBaseURL empty",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "",
				AzureDevopsOAuthAppID:        "mockAzureDevopsOAuthAppID",
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				EncryptionSecret:             "mockEncryptionSecret",
			},
			errMsg: "azure devops API base URL should not be empty",
		},
		{
			description: "invalid configuration with AzureDevopsOAuthAppID empty",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "mockAzureDevopsAPIBaseURL",
				AzureDevopsOAuthAppID:        "",
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				EncryptionSecret:             "mockEncryptionSecret",
			},
			errMsg: "azure devops OAuth app id should not be empty",
		},
		{
			description: "invalid configuration with AzureDevopsOAuthClientSecret empty",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "mockAzureDevopsAPIBaseURL",
				AzureDevopsOAuthAppID:        "mockAzureDevopsOAuthAppID",
				AzureDevopsOAuthClientSecret: "",
				EncryptionSecret:             "mockEncryptionSecret",
			},
			errMsg: "azure devops OAuth client secret should not be empty",
		},
		{
			description: "invalid configuration with EncryptionSecret empty",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "mockAzureDevopsAPIBaseURL",
				AzureDevopsOAuthAppID:        "mockAzureDevopsOAuthAppID",
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				EncryptionSecret:             "",
			},
			errMsg: "encryption secret should not be empty",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			err := testCase.config.IsValid()
			if testCase.errMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestProcessConfiguration(t *testing.T) {
	for _, testCase := range []struct {
		description        string
		config             *Configuration
		afterProcessConfig *Configuration
	}{
		{
			description: "valid process configuration for AzureDevopsAPIBaseURL",
			config: &Configuration{
				AzureDevopsAPIBaseURL: "  mockAzureDevopsAPIBaseURL/  ",
			},
			afterProcessConfig: &Configuration{
				AzureDevopsAPIBaseURL: "mockAzureDevopsAPIBaseURL",
			},
		},
		{
			description: "valid process configuration for AzureDevopsOAuthAppID",
			config: &Configuration{
				AzureDevopsOAuthAppID: "  mockAzureDevopsOAuthAppID  ",
			},
			afterProcessConfig: &Configuration{
				AzureDevopsOAuthAppID: "mockAzureDevopsOAuthAppID",
			},
		},
		{
			description: "valid process configuration for AzureDevopsOAuthClientSecret",
			config: &Configuration{
				AzureDevopsOAuthClientSecret: "  mockAzureDevopsOAuthClientSecret  ",
			},
			afterProcessConfig: &Configuration{
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
			},
		},
		{
			description: "valid process configuration for EncryptionSecret",
			config: &Configuration{
				EncryptionSecret: "  mockEncryptionSecret  ",
			},
			afterProcessConfig: &Configuration{
				EncryptionSecret: "mockEncryptionSecret",
			},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			err := testCase.config.ProcessConfiguration()
			assert.Nil(t, err)
			assert.Equal(t, testCase.afterProcessConfig, testCase.config)
		})
	}
}

func TestCloneConfiguration(t *testing.T) {
	for _, testCase := range []struct {
		description string
		config      *Configuration
	}{
		{
			description: "valid clone configuration",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "mockAzureDevopsAPIBaseURL",
				AzureDevopsOAuthAppID:        "mockAzureDevopsOAuthAppID",
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				EncryptionSecret:             "mockEncryptionSecret",
			},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			configuration := testCase.config.Clone()
			assert.NotNil(t, configuration)
			assert.Equal(t, configuration, testCase.config)
		})
	}
}
