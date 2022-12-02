package plugin

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
)

func TestClientGenerateOAuthToken(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)

	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "GenerateOAuthToken: valid",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "GenerateOAuthToken: with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GenerateOAuthToken(mockAPI.TestData().URLValues())

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestCreateTask(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "CreateTask: valid",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "CreateTask: with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.CreateTask(&serializers.CreateTaskRequestPayload{
				Fields: serializers.CreateTaskFieldValue{
					Description: "mockDescription",
				},
			}, "mockMattermostUSerID")

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestGetTask(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "GetTask: valid",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "GetTask: with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetTask("mockOrganization", "mockTaskID", "mockProjectName", "mockMattermostUSerID")

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestGetPullRequest(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "GetPullRequest: valid",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "GetPullRequest: with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetPullRequest("mockOrganization", "mockPullRequestID", "mockProjectName", "mockMattermostUSerID")

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestLink(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "Link: valid",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "Link: with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.Link(&serializers.LinkRequestPayload{}, "mockMattermostUSerID")

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestCreateSubscription(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "CreateSubscription: valid",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "CreateSubscription: with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.CreateSubscription(&serializers.CreateSubscriptionRequestPayload{}, &serializers.ProjectDetails{}, "mockChannelID", "mockPluginURL", "mockMattermostUSerID")

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestDeleteSubscription(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "DeleteSubscription: valid",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "DeleteSubscription: with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			statusCode, err := p.Client.DeleteSubscription("mockOrganization", "mockSubscriptionID", "mockMattermostUSerID")

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestCall(t *testing.T) {
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description string
	}{
		{
			description: "Call: valid",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(p), "AddAuthorization", func(_ *Plugin, _ *http.Request, _ string) error {
				return nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(p), "IsAccessTokenExpired", func(_ *Plugin, _ string) (bool, string) {
				return false, ""
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(p), "RefreshOAuthToken", func(_ *Plugin, _ string, _ string) error {
				return nil
			})
			client := &client{
				plugin:     p,
				httpClient: &http.Client{},
			}

			_, _, err := client.Call("mockBasePath", "mockMethod", "mockPath", "mockContentType", "mockMattermostUserID", nil, nil, url.Values{})
			assert.Error(t, err)
		})
	}
}

func TestGetGitRepositories(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "GetGitRepositories: valid",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "GetGitRepositories: with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetGitRepositories("mockOrganization", "mockProjectName", "mockMattermostUSerID")

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestGetGitRepositoryBranches(t *testing.T) {
	defer monkey.UnpatchAll()
	p := setupTestPlugin(&plugintest.API{})
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "GetGitRepositoryBranches: valid",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "GetGitRepositoryBranches: with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetGitRepositoryBranches("mockOrganization", "mockProjectName", "mockRepository", "mockMattermostUSerID")

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func setupTestPlugin(api *plugintest.API) *Plugin {
	p := Plugin{}
	p.API = api
	c := InitClient(&p)
	p.Client = c
	return &p
}
