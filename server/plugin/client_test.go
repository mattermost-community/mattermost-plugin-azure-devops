package plugin

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
)

func TestClientGenerateOAuthToken(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	c := InitClient(&p)
	p.Client = c
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "test ClientGenerateOAuthToken",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "test ClientGenerateOAuthToken with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(c), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	c := InitClient(&p)
	p.Client = c
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "test CreateTask",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "test CreateTask with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(c), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	c := InitClient(&p)
	p.Client = c
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "test GetTask",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "test GetTask with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(c), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
				return nil, testCase.statusCode, testCase.err
			})

			_, statusCode, err := p.Client.GetTask("mockOrganization", "mockTaskID", "mockMattermostUSerID")

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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	c := InitClient(&p)
	p.Client = c
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "test Link",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "test Link with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(c), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	c := InitClient(&p)
	p.Client = c
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "test CreateSubscription",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "test CreateSubscription with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(c), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	c := InitClient(&p)
	p.Client = c
	for _, testCase := range []struct {
		description string
		err         error
		statusCode  int
	}{
		{
			description: "test DeleteSubscription",
			err:         nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "test DeleteSubscription with error",
			err:         errors.New("mock-error"),
			statusCode:  http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(c), "Call", func(_ *client, basePath, method, path, contentType, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
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
	p := Plugin{}
	c := InitClient(&p)
	p.Client = c
	for _, testCase := range []struct {
		description string
	}{
		{
			description: "test Call",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "AddAuthorization", func(_ *Plugin, _ *http.Request, _ string) error {
				return nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "RefreshOAuthToken", func(_ *Plugin, _ string) error {
				return nil
			})
			client := &client{
				plugin:     &p,
				httpClient: &http.Client{},
			}

			_, _, err := client.Call("mockBasePath", "mockMethod", "mockPath", "mockContentType", "mockMattermostUserID", nil, nil, url.Values{})
			assert.Error(t, err)
		})
	}
}
