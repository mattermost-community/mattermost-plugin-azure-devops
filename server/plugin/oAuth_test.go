package plugin

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mattermost/mattermost-plugin-azure-devops/mocks"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/config"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/testutils"
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
			mattermostUserID: testutils.MockMattermostUserID,
			statusCode:       http.StatusFound,
		},
		{
			description: "OAuthConnect: without mattermostUserID",
			statusCode:  http.StatusUnauthorized,
		},
		{
			description:      "OAuthConnect: user already connected",
			isConnected:      true,
			mattermostUserID: testutils.MockMattermostUserID,
			statusCode:       http.StatusBadRequest,
		},
		{
			description:      "OAuthConnect: user already connected and failed to DM",
			isConnected:      true,
			DMErr:            &model.AppError{},
			mattermostUserID: testutils.MockMattermostUserID,
			statusCode:       http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "MattermostUserAlreadyConnected", func(_ *Plugin, _ string) bool {
				return testCase.isConnected
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "GenerateOAuthConnectURL", func(_ *Plugin, _ string) string {
				return "mockRedirectURL"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "CloseBrowserWindowWithHTTPResponse", func(_ *Plugin, _ http.ResponseWriter) {})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DM", func(_ *Plugin, _, _ string, _ bool, _ ...interface{}) (string, error) {
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
			oAuthTokenErr: errors.New("oAuthTokenErr"),
			statusCode:    http.StatusInternalServerError,
		},
		{
			description:   "OAuthComplete: Azure DevOps user is already connected",
			code:          "mockCode",
			state:         "mock_State",
			oAuthTokenErr: errors.New(constants.ErrorMessageAzureDevopsAccountAlreadyConnected),
			statusCode:    http.StatusForbidden,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "GenerateOAuthToken", func(_ *Plugin, _, _ string) error {
				return testCase.oAuthTokenErr
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "CloseBrowserWindowWithHTTPResponse", func(_ *Plugin, _ http.ResponseWriter) {})
			mockAPI.On("LogError", testutils.GetMockArgumentsWithType("string", 3)...)

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
			mockAPI.On("PublishWebSocketEvent", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("*model.WebsocketBroadcast")).Return(nil)

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DM", func(_ *Plugin, _, _ string, _ bool, _ ...interface{}) (string, error) {
				return "", testCase.DMError
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "GenerateAndStoreOAuthToken", func(_ *Plugin, _ string, _ url.Values, _ bool) error {
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
			description:   "RefreshOAuthToken: token is not decoded successfully",
			user:          &serializers.User{},
			decodeError:   errors.New("error while decoding token"),
			expectedError: "error while decoding token",
		},
		{
			description:   "RefreshOAuthToken: token is not decoded successfully and DM error occurs",
			user:          &serializers.User{},
			decodeError:   errors.New("error decoding token"),
			DMErr:         errors.New("error sending direct message"),
			expectedError: "error sending direct message",
		},
		{
			description:   "RefreshOAuthToken: token is not decrypted successfully",
			user:          &serializers.User{},
			decryptError:  errors.New("error decrypting token"),
			expectedError: "error decrypting token",
		},
		{
			description:   "RefreshOAuthToken: token is not decrypted successfully and DM error occurs",
			user:          &serializers.User{},
			decryptError:  errors.New("error decrypting token"),
			DMErr:         errors.New("error sending direct message"),
			expectedError: "error sending direct message",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			p.setConfiguration(
				&config.Configuration{
					EncryptionSecret:             "mockEncryptionSecret",
					AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DM", func(_ *Plugin, _, _ string, _ bool, _ ...interface{}) (string, error) {
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
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "GenerateAndStoreOAuthToken", func(_ *Plugin, _ string, _ url.Values, _ bool) error {
				return nil
			})

			err := p.RefreshOAuthToken(testutils.MockMattermostUserID, "mockRefreshToken")
			if testCase.expectedError != "" {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGenerateAndStoreOAuthToken(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.Store = mockedStore
	p.Client = mockedClient
	for _, testCase := range []struct {
		description    string
		storeUserError error
		DMErr          error
		expectedError  string
		storeError     error
	}{
		{
			description: "GenerateAndStoreOAuthToken: valid",
		},
		{
			description:    "GenerateAndStoreOAuthToken: storing user gives error",
			storeUserError: errors.New("error storing user"),
			expectedError:  "error storing user",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockedClient.EXPECT().GenerateOAuthToken(gomock.Any()).Return(&serializers.OAuthSuccessResponse{}, 200, nil)
			mockedClient.EXPECT().GetUserProfile("me", "").Return(&serializers.UserProfile{}, 200, nil)
			mockedStore.EXPECT().LoadAzureDevopsUserDetails("").Return(&serializers.User{}, nil)

			monkey.Patch(strconv.Atoi, func(string) (int, error) {
				return 0, nil
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DM", func(_ *Plugin, _, _ string, _ bool, _ ...interface{}) (string, error) {
				return "", testCase.DMErr
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "Encrypt", func(_ *Plugin, _, _ []byte) ([]byte, error) {
				return nil, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "Encode", func(_ *Plugin, _ []byte) string {
				return ""
			})

			if testCase.storeError == nil {
				mockedStore.EXPECT().StoreAzureDevopsUserDetailsWithMattermostUserID(&serializers.User{
					ExpiresAt: time.Now().UTC().Add(time.Second * time.Duration(0)).Unix(),
				}).Return(testCase.storeUserError)
			}

			p.setConfiguration(
				&config.Configuration{
					EncryptionSecret:             "mockEncryptionSecret",
					AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				})

			err := p.GenerateAndStoreOAuthToken("", nil, false)
			if testCase.expectedError != "" {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestIsAccessTokenExpired(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Store = mockedStore
	for _, testCase := range []struct {
		description   string
		loadUserError error
		clientError   error
		DMErr         error
		expectedError string
	}{
		{
			description: "IsAccessTokenExpired: valid",
		},
		{
			description:   "IsAccessTokenExpired: loading user gives error",
			loadUserError: errors.New("error loading user"),
			expectedError: "error loading user",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", testutils.GetMockArgumentsWithType("string", 3)...)
			mockedStore.EXPECT().LoadAzureDevopsUserIDFromMattermostUser(testutils.MockMattermostUserID).Return(testutils.MockAzureDevopsUserID, nil)
			mockedStore.EXPECT().LoadAzureDevopsUserDetails(testutils.MockAzureDevopsUserID).Return(&serializers.User{}, nil)

			p.setConfiguration(
				&config.Configuration{
					EncryptionSecret:             "mockEncryptionSecret",
					AzureDevopsOAuthClientSecret: "mockAzureDevopsOAuthClientSecret",
				})

			isAccessTokenExpired, err := p.IsAccessTokenExpired(testutils.MockMattermostUserID)
			if testCase.expectedError != "" {
				assert.NotNil(t, err)
				assert.NotNil(t, isAccessTokenExpired)
				return
			}

			assert.NotNil(t, isAccessTokenExpired)
			assert.Equal(t, "", err)
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
			description: "MattermostUserAlreadyConnected: valid",
			user:        &serializers.User{},
		},
		{
			description:   "MattermostUserAlreadyConnected: user is not loaded successfully",
			loadUserError: errors.New("error loading user"),
		},
	} {
		mockAPI.On("LogError", testutils.GetMockArgumentsWithType("string", 3)...).Return(nil)

		mockedStore.EXPECT().LoadAzureDevopsUserIDFromMattermostUser(testutils.MockMattermostUserID).Return(testutils.MockAzureDevopsUserID, nil)
		mockedStore.EXPECT().LoadAzureDevopsUserDetails(testutils.MockAzureDevopsUserID).Return(&serializers.User{}, nil)

		t.Run(testCase.description, func(t *testing.T) {
			resp := p.MattermostUserAlreadyConnected(testutils.MockMattermostUserID)
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
