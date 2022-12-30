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
					AreaPath:    "mockAreaPath",
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

			_, statusCode, err := p.Client.GetPullRequest("mockOrganization", "mockPullRequestID", "mockProjectName", "mockMattermostUserID")

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestGetBuildDetails(t *testing.T) {
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
			description: "GetBuildDetails: valid",
			statusCode:  http.StatusOK,
		},
		{
			description:          "GetBuildDetails: with error",
			err:                  errors.New("failed to get build details"),
			statusCode:           http.StatusInternalServerError,
			expectedErrorMessage: "failed to get the pipeline build details: failed to get build details",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetBuildDetails("mockOrganization", "mockProjectName", "mockBuildID", "mockMattermostUserID")

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.expectedErrorMessage)
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

func TestCheckIfUserIsProjectAdmin(t *testing.T) {
	defer monkey.UnpatchAll()
	p := setupTestPlugin(&plugintest.API{})
	for _, testCase := range []struct {
		description   string
		err           error
		statusCode    int
		expectedError string
	}{
		{
			description: "CheckIfUserIsProjectAdmin: valid",
			statusCode:  http.StatusOK,
		},
		{
			description:   "CheckIfUserIsProjectAdmin: with error",
			err:           errors.New("failed to check user permissions"),
			statusCode:    http.StatusInternalServerError,
			expectedError: "failed to check if user is a project admin: failed to check user permissions",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			statusCode, err := p.Client.CheckIfUserIsProjectAdmin("mockOrganization", "mockProjectID", "mockProjectURL", "mockMattermostUSerID")

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestGetSubscriptionFilterPossibleValues(t *testing.T) {
	defer monkey.UnpatchAll()
	p := setupTestPlugin(&plugintest.API{})
	for _, testCase := range []struct {
		description          string
		err                  error
		statusCode           int
		request              *serializers.GetSubscriptionFilterPossibleValuesRequestPayload
		expectedErrorMessage string
	}{
		{
			description: "GetSubscriptionFilterPossibleValues: valid",
			statusCode:  http.StatusOK,
			request: &serializers.GetSubscriptionFilterPossibleValuesRequestPayload{
				Filters:      []string{"mockFilter1", "mockFilter2"},
				EventType:    "mockEventType",
				ProjectID:    "mockProjectID",
				RepositoryID: "mockRepositoryID",
			},
		},
		{
			description:          "GetSubscriptionFilterPossibleValues: with error",
			err:                  errors.New("error in getting subscription filter values"),
			statusCode:           http.StatusInternalServerError,
			request:              &serializers.GetSubscriptionFilterPossibleValuesRequestPayload{},
			expectedErrorMessage: "failed to get the subscription filter values: error in getting subscription filter values",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetSubscriptionFilterPossibleValues(testCase.request, "mockMattermostUserID")

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.expectedErrorMessage)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestUpdatePipelineApprovalRequest(t *testing.T) {
	defer monkey.UnpatchAll()
	p := setupTestPlugin(&plugintest.API{})
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "UpdatePipelineApprovalRequest: valid",
			statusCode:  http.StatusOK,
		},
		{
			description: "UpdatePipelineApprovalRequest: with error",
			err:         errors.New("failed to update pipeline approval request"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			statusCode, err := p.Client.UpdatePipelineApprovalRequest(&serializers.PipelineApproveRequest{}, "mockOrganization", "mockProjectID", "mockMattermostUSerID", 1234)

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestGetApprovalDetails(t *testing.T) {
	defer monkey.UnpatchAll()
	p := setupTestPlugin(&plugintest.API{})
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "GetApprovalDetails: valid",
			statusCode:  http.StatusOK,
		},
		{
			description: "GetApprovalDetails: with error",
			err:         errors.New("failed to get approval details"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetApprovalDetails("mockOrganization", "mockProjectID", "mockMattermostUSerID", 1234)

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.err.Error())
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
