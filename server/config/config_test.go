package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
)

func TestIsValid(t *testing.T) {
	for _, testCase := range []struct {
		description string
		config      *Configuration
		errMsg      string
	}{
		{
			description: "configuration: valid",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "mockAzureDevopsAPIBaseURL",
				AzureDevopsOAuthAppID:        "mockAzureDevopsOAuthAppID",
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				EncryptionSecret:             "mockEncryptionSecret",
			},
		},
		{
			description: "configuration: empty AzureDevopsAPIBaseURL",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "",
				AzureDevopsOAuthAppID:        "mockAzureDevopsOAuthAppID",
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				EncryptionSecret:             "mockEncryptionSecret",
			},
			errMsg: constants.EmptyAzureDevopsAPIBaseURLError,
		},
		{
			description: "configuration: empty AzureDevopsOAuthAppID",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "mockAzureDevopsAPIBaseURL",
				AzureDevopsOAuthAppID:        "",
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				EncryptionSecret:             "mockEncryptionSecret",
			},
			errMsg: constants.EmptyAzureDevopsOAuthAppIDError,
		},
		{
			description: "configuration: empty AzureDevopsOAuthClientSecret",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "mockAzureDevopsAPIBaseURL",
				AzureDevopsOAuthAppID:        "mockAzureDevopsOAuthAppID",
				AzureDevopsOAuthClientSecret: "",
				EncryptionSecret:             "mockEncryptionSecret",
			},
			errMsg: constants.EmptyAzureDevopsOAuthClientSecretError,
		},
		{
			description: "configuration: empty EncryptionSecret",
			config: &Configuration{
				AzureDevopsAPIBaseURL:        "mockAzureDevopsAPIBaseURL",
				AzureDevopsOAuthAppID:        "mockAzureDevopsOAuthAppID",
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				EncryptionSecret:             "",
			},
			errMsg: constants.EmptyEncryptionSecretError,
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
			description: "ProcessConfiguration: valid AzureDevopsAPIBaseURL",
			config: &Configuration{
				AzureDevopsAPIBaseURL: "  mockAzureDevopsAPIBaseURL/  ",
			},
			afterProcessConfig: &Configuration{
				AzureDevopsAPIBaseURL: "mockAzureDevopsAPIBaseURL",
			},
		},
		{
			description: "ProcessConfiguration: valid AzureDevopsOAuthAppID",
			config: &Configuration{
				AzureDevopsOAuthAppID: "  mockAzureDevopsOAuthAppID  ",
			},
			afterProcessConfig: &Configuration{
				AzureDevopsOAuthAppID: "mockAzureDevopsOAuthAppID",
			},
		},
		{
			description: "ProcessConfiguration: valid AzureDevopsOAuthClientSecret",
			config: &Configuration{
				AzureDevopsOAuthClientSecret: "  mockAzureDevopsOAuthClientSecret  ",
			},
			afterProcessConfig: &Configuration{
				AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
			},
		},
		{
			description: "ProcessConfiguration: valid EncryptionSecret",
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
			description: "CloneConfiguration: valid",
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
