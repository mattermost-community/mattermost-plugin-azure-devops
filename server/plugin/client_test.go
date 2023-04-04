package plugin

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
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
			err:         errors.New("error generating oAuth token"),
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
			err:         errors.New("error creating the task"),
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
			err:         errors.New("error getting the task"),
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
			err:         errors.New("error getting the pull request"),
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

			_, statusCode, err := p.Client.GetBuildDetails(testutils.MockOrganization, testutils.MockProjectName, "mockBuildID", testutils.MockMattermostUserID)

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
			err:         errors.New("error linking the project"),
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
			err:         errors.New("error creating subscription"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.CreateSubscription(&serializers.CreateSubscriptionRequestPayload{}, &serializers.ProjectDetails{}, testutils.MockChannelID, "mockPluginURL", testutils.MockMattermostUserID, "mockUUID")

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
			err:         errors.New("error deleting the subscription"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			statusCode, err := p.Client.DeleteSubscription(testutils.MockOrganization, testutils.MockSubscriptionID, testutils.MockMattermostUserID)

			if testCase.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestOpenDialogRequest(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "OpenDialogRequest: valid",
			statusCode:  http.StatusOK,
		},
		{
			description: "OpenDialogRequest: with error",
			err:         errors.New("error in request to open comment dialog"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			statusCode, err := p.Client.OpenDialogRequest(&model.OpenDialogRequest{}, testutils.MockMattermostUserID)

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.err.Error())
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

func TestUpdatePipelineRunApprovalRequest(t *testing.T) {
	defer monkey.UnpatchAll()
	p := setupTestPlugin(&plugintest.API{})
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "UpdatePipelineRunApprovalRequest: valid",
			statusCode:  http.StatusOK,
		},
		{
			description: "UpdatePipelineRunApprovalRequest: with error",
			err:         errors.New("error updating pipeline run approval request"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.UpdatePipelineRunApprovalRequest([]*serializers.PipelineApproveRequest{}, testutils.MockProjectID, testutils.MockMattermostUserID, testutils.MockApproverID)

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestGetRunApprovalDetails(t *testing.T) {
	defer monkey.UnpatchAll()
	p := setupTestPlugin(&plugintest.API{})
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "GetRunApprovalDetails: valid",
			statusCode:  http.StatusOK,
		},
		{
			description: "GetRunApprovalDetails: with error",
			err:         errors.New("error getting run approval details"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetRunApprovalDetails(testutils.MockOrganization, testutils.MockProjectID, testutils.MockMattermostUserID, testutils.MockApproverID)

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.err.Error())
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

			statusCode, err := p.Client.UpdatePipelineApprovalRequest(&serializers.PipelineApproveRequest{}, testutils.MockOrganization, testutils.MockProjectID, testutils.MockMattermostUserID, 1234)

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

			_, statusCode, err := p.Client.GetApprovalDetails(testutils.MockOrganization, testutils.MockProjectID, testutils.MockMattermostUserID, 1234)

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.err.Error())
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

			_, statusCode, err := p.Client.GetSubscriptionFilterPossibleValues(testCase.request, testutils.MockMattermostUserID)

			if testCase.err != nil {
				assert.EqualError(t, err, testCase.expectedErrorMessage)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestMakeHTTPRequest(t *testing.T) {
	mockAPI := &plugintest.API{}
	p := setupTestPlugin(mockAPI)
	for _, testCase := range []struct {
		description                       string
		maxBytesSizeForClientResponseBody int
	}{
		{
			description:                       "MakeHTTPRequest: valid",
			maxBytesSizeForClientResponseBody: constants.MaxBytesSizeForReadingResponseBody,
		},
		{
			description:                       "MakeHTTPRequest: large response body",
			maxBytesSizeForClientResponseBody: constants.MaxBytesSizeForReadingResponseBody + 1,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			// Start a local HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Send response to be tested
				respBody := testutils.GenerateStringOfSize(testCase.maxBytesSizeForClientResponseBody)
				if _, err := rw.Write([]byte(respBody)); err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
				}
			}))
			// Close the server when test finishes
			defer server.Close()

			client := &client{
				plugin:     p,
				httpClient: server.Client(),
			}

			req := httptest.NewRequest(http.MethodGet, server.URL, nil)
			req.RequestURI = ""
			_, _, err := client.MakeHTTPRequest(req, "", nil)

			if testCase.maxBytesSizeForClientResponseBody > constants.MaxBytesSizeForReadingResponseBody {
				assert.Errorf(t, err, "http: request body too large")
			} else {
				assert.NoError(t, err)
			}
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
