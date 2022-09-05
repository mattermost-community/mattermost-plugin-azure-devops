package plugin

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/Brightscout/mattermost-plugin-azure-devops/mocks"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/config"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOAuthConnect(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	for _, testCase := range []struct {
		description      string
		isConnected      bool
		mattermostUserID string
		DMErr            error
		statusCode       int
	}{
		{
			description:      "OAuthConnect: valid",
			mattermostUserID: "mockMattermostUserID",
			statusCode:       http.StatusFound,
		},
		{
			description: "OAuthConnect: without mattermostUserID",
			statusCode:  http.StatusUnauthorized,
		},
		{
			description:      "OAuthConnect: user already connected",
			isConnected:      true,
			mattermostUserID: "mockMattermostUserID",
			statusCode:       http.StatusBadRequest,
		},
		{
			description:      "OAuthConnect: user already connected and failed to DM",
			isConnected:      true,
			DMErr:            &model.AppError{},
			mattermostUserID: "mockMattermostUserID",
			statusCode:       http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "UserAlreadyConnected", func(_ *Plugin, _ string) bool {
				return testCase.isConnected
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "GenerateOAuthConnectURL", func(_ *Plugin, _ string) string {
				return "mockRedirectURL"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "CloseBrowserWindowWithHTTPResponse", func(_ *Plugin, _ http.ResponseWriter) {})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DM", func(_ *Plugin, _, _ string, _ ...interface{}) (string, error) {
				return "", testCase.DMErr
			})

			req := httptest.NewRequest(http.MethodGet, "/oauth/connect", bytes.NewBufferString(`{}`))
			req.Header.Add(constants.HeaderMattermostUserID, testCase.mattermostUserID)

			res := httptest.NewRecorder()

			p.OAuthConnect(res, req)
			assert.Equal(t, testCase.statusCode, res.Code)
		})
	}
}

func TestOAuthComplete(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	for _, testCase := range []struct {
		description   string
		code          string
		state         string
		oAuthTokenErr error
		statusCode    int
	}{
		{
			description: "OAuthComplete: valid",
			code:        "mockCode",
			state:       "mock_State",
			statusCode:  http.StatusOK,
		},
		{
			description: "OAuthComplete: without code",
			state:       "mock_State",
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "OAuthComplete: without state",
			code:        "mockCode",
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "OAuthComplete: length of state not equal to 2",
			code:        "mockCode",
			state:       "mockState",
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "OAuthComplete: state second word empty",
			code:        "mockCode",
			state:       "mockState_",
			statusCode:  http.StatusBadRequest,
		},
		{
			description:   "OAuthComplete: with oAuthTokenErr",
			code:          "mockCode",
			state:         "mock_State",
			oAuthTokenErr: errors.New("mockError"),
			statusCode:    http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "GenerateOAuthToken", func(_ *Plugin, _, _ string) error {
				return testCase.oAuthTokenErr
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "CloseBrowserWindowWithHTTPResponse", func(_ *Plugin, _ http.ResponseWriter) {})

			req := httptest.NewRequest(http.MethodGet, "/oauth/complete", bytes.NewBufferString(`{}`))
			q := req.URL.Query()
			q.Add("code", testCase.code)
			q.Add("state", testCase.state)
			req.URL.RawQuery = q.Encode()
			res := httptest.NewRecorder()

			p.OAuthComplete(res, req)
			assert.Equal(t, testCase.statusCode, res.Code)
		})
	}
}

func TestGenerateOAuthToken(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Store = mockedStore
	for _, testCase := range []struct {
		description      string
		code             string
		state            string
		verifyOAuthError error
		expectedError    string
		DMError          error
	}{
		{
			description: "GenerateOAuthToken: valid",
			code:        "mockCode",
			state:       "mock_state",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DM", func(_ *Plugin, _, _ string, _ ...interface{}) (string, error) {
				return "", testCase.DMError
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "StoreOAuthToken", func(_ *Plugin, _ string, _ url.Values) error {
				return nil
			})

			if testCase.expectedError == "" {
				mockedStore.EXPECT().VerifyOAuthState("state", testCase.state).Return(testCase.verifyOAuthError)
			}

			err := p.GenerateOAuthToken(testCase.code, testCase.state)
			if testCase.expectedError != "" {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestRefreshOAuthToken(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.Store = mockedStore
	for _, testCase := range []struct {
		description   string
		decodeError   error
		decryptError  error
		user          *serializers.User
		loadUserErr   error
		DMErr         error
		expectedError string
	}{
		{
			description: "RefreshOAuthToken: token is parsed successfully",
			user: &serializers.User{
				RefreshToken: "mockRefreshToken",
			},
		},
		{
			description:   "RefreshOAuthToken: user is not loaded successfully",
			loadUserErr:   errors.New("mockError"),
			expectedError: "mockError",
		},
		{
			description:   "RefreshOAuthToken: user is not loaded successfully and DM error occurs",
			loadUserErr:   errors.New("mockError"),
			DMErr:         errors.New("mockError"),
			expectedError: "mockError",
		},
		{
			description:   "RefreshOAuthToken: token is not decoded successfully",
			user:          &serializers.User{},
			decodeError:   errors.New("mockError"),
			expectedError: "mockError",
		},
		{
			description:   "RefreshOAuthToken: token is not decoded successfully and DM error occurs",
			user:          &serializers.User{},
			decodeError:   errors.New("mockError"),
			DMErr:         errors.New("mockError"),
			expectedError: "mockError",
		},
		{
			description:   "RefreshOAuthToken: token is not decrypted successfully",
			user:          &serializers.User{},
			decryptError:  errors.New("mockError"),
			expectedError: "mockError",
		},
		{
			description:   "RefreshOAuthToken: token is not decrypted successfully and DM error occurs",
			user:          &serializers.User{},
			decryptError:  errors.New("mockError"),
			DMErr:         errors.New("mockError"),
			expectedError: "mockError",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			p.setConfiguration(
				&config.Configuration{
					EncryptionSecret:             "mockEncryptionSecret",
					AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DM", func(_ *Plugin, _, _ string, _ ...interface{}) (string, error) {
				return "", testCase.DMErr
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "Decode", func(_ *Plugin, _ string) ([]byte, error) {
				return nil, testCase.decodeError
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "GetSiteURL", func(_ *Plugin) string {
				return "mockSiteURL"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "GetPluginURLPath", func(_ *Plugin) string {
				return "mockPluginURLPath"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "Decrypt", func(_ *Plugin, _, _ []byte) ([]byte, error) {
				return nil, testCase.decryptError
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "StoreOAuthToken", func(_ *Plugin, _ string, _ url.Values) error {
				return nil
			})

			mockedStore.EXPECT().LoadUser("mockMattermostUserID").Return(testCase.user, testCase.loadUserErr)

			err := p.RefreshOAuthToken("mockMattermostUserID", "mockRefreshToken")
			if testCase.expectedError != "" {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestStoreOAuthToken(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Store = mockedStore
	p.Client = mockedClient
	for _, testCase := range []struct {
		description    string
		user           *serializers.User
		storeUserError error
		DMErr          error
		expectedError  string
	}{
		{
			description: "StoreOAuthToken: valid",
			user:        &serializers.User{},
		},
		{
			description:    "StoreOAuthToken: storing user gives error",
			user:           &serializers.User{},
			storeUserError: errors.New("mockError"),
			expectedError:  "mockError",
		},
		{
			description:   "StoreOAuthToken: DM gives error",
			user:          &serializers.User{},
			DMErr:         errors.New("mockError"),
			expectedError: "mockError",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("PublishWebSocketEvent", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("*model.WebsocketBroadcast")).Return(nil)

			mockedClient.EXPECT().GenerateOAuthToken(gomock.Any()).Return(&serializers.OAuthSuccessResponse{}, 200, nil)

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DM", func(_ *Plugin, _, _ string, _ ...interface{}) (string, error) {
				return "", testCase.DMErr
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "Encrypt", func(_ *Plugin, _, _ []byte) ([]byte, error) {
				return nil, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "Encode", func(_ *Plugin, _ []byte) string {
				return ""
			})

			mockedStore.EXPECT().StoreUser(testCase.user).Return(testCase.storeUserError)

			p.setConfiguration(
				&config.Configuration{
					EncryptionSecret:             "mockEncryptionSecret",
					AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				})

			err := p.GenerateAndStoreOAuthToken("", nil)
			if testCase.expectedError != "" {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestUserAlreadyConnected(t *testing.T) {
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.Store = mockedStore
	for _, testCase := range []struct {
		description   string
		user          *serializers.User
		loadUserError error
	}{
		{
			description: "UserAlreadyConnected: valid",
			user:        &serializers.User{},
		},
		{
			description:   "UserAlreadyConnected: user is not loaded successfully",
			loadUserError: errors.New("mockError"),
		},
	} {
		mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

		mockedStore.EXPECT().LoadUser("mockMattermostUserID").Return(testCase.user, testCase.loadUserError)

		t.Run(testCase.description, func(t *testing.T) {
			resp := p.UserAlreadyConnected("mockMattermostUserID")
			assert.NotNil(t, resp)
		})
	}
}

func TestCloseBrowserWindowWithHTTPResponse(t *testing.T) {
	p := Plugin{}
	for _, testCase := range []struct {
		description string
		html        string
		statusCode  int
	}{
		{
			description: "CloseBrowserWindowWithHTTPResponse: valid",
			html:        "mockHTML",
			statusCode:  http.StatusOK,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			res := httptest.NewRecorder()

			p.CloseBrowserWindowWithHTTPResponse(res)
			assert.Equal(t, testCase.statusCode, res.Code)
		})
	}
}
