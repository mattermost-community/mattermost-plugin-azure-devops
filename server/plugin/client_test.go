package plugin

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/testutils"
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
			}, testutils.MockMattermostUserID)

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

			_, statusCode, err := p.Client.GetTask(testutils.MockOrganization, "mockTaskID", testutils.MockProjectName, testutils.MockMattermostUserID)

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestGetReleaseDetails(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description          string
		err                  error
		statusCode           int
		expectedErrorMessage string
	}{
		{
			description: "GetReleaseDetails: valid",
			statusCode:  http.StatusOK,
		},
		{
			description:          "GetReleaseDetails: with error",
			err:                  errors.New("failed to get release details"),
			statusCode:           http.StatusInternalServerError,
			expectedErrorMessage: "failed to get the pipeline release details: failed to get release details",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetReleaseDetails(testutils.MockOrganization, testutils.MockProjectName, "mockReleaseID", testutils.MockMattermostUserID)

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.expectedErrorMessage)
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

			_, statusCode, err := p.Client.GetPullRequest(testutils.MockOrganization, "mockPullRequestID", testutils.MockProjectName, testutils.MockMattermostUserID)

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

			_, statusCode, err := p.Client.Link(&serializers.LinkRequestPayload{}, testutils.MockMattermostUserID)

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

			_, statusCode, err := p.Client.CreateSubscription(&serializers.CreateSubscriptionRequestPayload{}, &serializers.ProjectDetails{}, testutils.MockChannelID, "mockPluginURL", testutils.MockMattermostUserID)

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

			statusCode, err := p.Client.DeleteSubscription(testutils.MockOrganization, "mockSubscriptionID", testutils.MockMattermostUserID)

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

			_, _, err := client.Call("mockBasePath", "mockMethod", "mockPath", "mockContentType", testutils.MockMattermostUserID, nil, nil, url.Values{})
			assert.Error(t, err)
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
