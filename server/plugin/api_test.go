package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-azure-devops/mocks"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/testutils"
)

type panicHandler struct {
}

func (ph panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("bad handler")
}

func setupMockPlugin(api *plugintest.API, store *mocks.MockKVStore, client *mocks.MockClient) *Plugin {
	p := &Plugin{}
	p.API = api
	if store != nil {
		p.Store = store
	}

	if client != nil {
		p.Client = client
	}
	p.router = p.InitAPI()
	return p
}

func TestInitRoutes(t *testing.T) {
	p := setupMockPlugin(&plugintest.API{}, nil, nil)
	p.InitRoutes()
}

func TestHandleStaticFiles(t *testing.T) {
	mockAPI := &plugintest.API{}
	p := setupMockPlugin(mockAPI, nil, nil)
	mockAPI.On("GetBundlePath").Return("/test-path", nil)
	p.HandleStaticFiles()
}

func TestWithRecovery(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			require.Fail(t, "got panic")
		}
	}()

	mockAPI := &plugintest.API{}
	p := setupMockPlugin(mockAPI, nil, nil)
	mockAPI.On("LogError", "Recovered from a panic", "url", "http://random", "error", "bad handler", "stack", mock.Anything)

	ph := panicHandler{}
	handler := p.WithRecovery(ph)

	req := httptest.NewRequest(http.MethodGet, "http://random", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.Body != nil {
		defer resp.Body.Close()
		_, err := io.Copy(io.Discard, resp.Body)
		require.NoError(t, err)
	}
}

func TestHandleAuthRequired(t *testing.T) {
	p := setupMockPlugin(&plugintest.API{}, nil, nil)
	for _, testCase := range []struct {
		description string
	}{
		{
			description: "HandleAuthRequired: valid",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			timerHandler := func(w http.ResponseWriter, r *http.Request) {}

			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{}`))
			req.Header.Add(constants.HeaderMattermostUserID, "mockUserID")

			res := httptest.NewRecorder()

			timerHandler(res, req)

			resp := p.handleAuthRequired(timerHandler)
			assert.NotNil(t, resp)
		})
	}
}

func TestHandleCreateTask(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(mockAPI, nil, mockedClient)
	for _, testCase := range []struct {
		description        string
		body               string
		err                error
		marshalError       error
		statusCode         int
		expectedStatusCode int
		clientError        error
	}{
		{
			description: "CreateTask: valid fields",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",
				"type": "mockType",
				"fields": {
					"title": "mockTitle",
					"description": "mockDescription"
					}
				}`,
			err:                nil,
			statusCode:         http.StatusOK,
			expectedStatusCode: http.StatusOK,
		},
		{
			description:        "CreateTask: empty body",
			body:               `{}`,
			err:                errors.New("mockError"),
			statusCode:         http.StatusBadRequest,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "CreateTask: invalid body",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",
				"type": "mockType",`,
			err:                errors.New("mockError"),
			statusCode:         http.StatusBadRequest,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "CreateTask: missing fields",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",
				"type": "mockType"
				}`,
			err:                errors.New("mockError"),
			statusCode:         http.StatusBadRequest,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "CreateTask: marshaling gives error",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",
				"type": "mockType",
				"fields": {
					"title": "mockTitle",
					"description": "mockDescription"
					}
				}`,
			marshalError:       errors.New("mockError"),
			statusCode:         http.StatusOK,
			expectedStatusCode: http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			mockAPI.On("GetDirectChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&model.Channel{}, nil)
			mockAPI.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)

			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})

			if testCase.statusCode == http.StatusOK {
				mockedClient.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(&serializers.TaskValue{}, testCase.statusCode, testCase.err)
			}

			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(testCase.body))
			req.Header.Add(constants.HeaderMattermostUserID, "mockUserID")

			w := httptest.NewRecorder()
			p.handleCreateTask(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestHandleLink(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, mockedClient)
	for _, testCase := range []struct {
		description     string
		body            string
		err             error
		statusCode      int
		projectList     []serializers.ProjectDetails
		project         serializers.ProjectDetails
		isProjectLinked bool
	}{
		{
			description: "HandleLink: valid",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject"
				}`,
			statusCode: http.StatusOK,
			projectList: []serializers.ProjectDetails{
				{
					MattermostUserID: "mockMattermostUserID",
					ProjectName:      "mockProject",
					OrganizationName: "mockOrganizationName",
					ProjectID:        "mockProjectID",
				},
			},
			project: serializers.ProjectDetails{
				MattermostUserID: "mockMattermostUserID",
				ProjectName:      "mockProject",
				OrganizationName: "mockOrganizationName",
				ProjectID:        "mockProjectID",
			},
		},
		{
			description: "HandleLink: empty body",
			body:        `{}`,
			err:         errors.New("mockError"),
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "HandleLink: invalid body",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",`,
			err:        errors.New("mockError"),
			statusCode: http.StatusBadRequest,
		},
		{
			description: "HandleLink: missing fields",
			body: `{
				"organization": "mockOrganization",
				}`,
			err:        errors.New("mockError"),
			statusCode: http.StatusBadRequest,
		},
		{
			description: "HandleLink: project is already linked",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject"
				}`,
			statusCode: http.StatusOK,
			projectList: []serializers.ProjectDetails{
				{
					MattermostUserID: "mockMattermostUserID",
					ProjectName:      "mockProject",
					OrganizationName: "mockOrganizationName",
					ProjectID:        "mockProjectID",
				},
			},
			isProjectLinked: true,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			mockAPI.On("GetDirectChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&model.Channel{}, nil)
			mockAPI.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)

			monkey.PatchInstanceMethod(reflect.TypeOf(p), "IsProjectLinked", func(*Plugin, []serializers.ProjectDetails, serializers.ProjectDetails) (*serializers.ProjectDetails, bool) {
				return &serializers.ProjectDetails{}, testCase.isProjectLinked
			})

			if testCase.statusCode == http.StatusOK {
				mockedStore.EXPECT().GetAllProjects("mockMattermostUserID").Return(testCase.projectList, nil)
			}

			if testCase.statusCode == http.StatusOK && !testCase.isProjectLinked {
				mockedClient.EXPECT().Link(gomock.Any(), gomock.Any()).Return(&serializers.Project{}, testCase.statusCode, testCase.err)
				mockedClient.EXPECT().CheckIfUserIsProjectAdmin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
				mockedStore.EXPECT().StoreProject(&serializers.ProjectDetails{
					MattermostUserID: "mockMattermostUserID",
					ProjectName:      "Mockproject",
					OrganizationName: "mockorganization",
				}).Return(nil)
			}

			req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewBufferString(testCase.body))
			req.Header.Add(constants.HeaderMattermostUserID, "mockMattermostUserID")

			w := httptest.NewRecorder()
			p.handleLink(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)
		})
	}
}

func TestHandleDeleteAllSubscriptions(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, mockedClient)
	for _, testCase := range []struct {
		description           string
		userID                string
		projectID             string
		err                   error
		statusCode            int
		getAllSubscriptionErr error
		subscriptionList      []*serializers.SubscriptionDetails
		subscription          *serializers.SubscriptionDetails
	}{
		{
			description: "HandleDeleteAllSubscriptions: valid",
			userID:      "mockMattermostUserID",
			projectID:   "mockProjectID",
			statusCode:  http.StatusOK,
			subscriptionList: []*serializers.SubscriptionDetails{
				{
					MattermostUserID: "mockMattermostUserID",
					ProjectID:        "mockProjectID",
					OrganizationName: "mockOrganization",
					EventType:        "mockEventType",
					ChannelID:        "mockChannelID",
					SubscriptionID:   "mockSubscriptionID",
				},
			},
		},
		{
			description:           "HandleDeleteAllSubscriptions: GetAllSubscriptions gives error",
			userID:                "mockMattermostUserID",
			projectID:             "mockProjectID",
			statusCode:            http.StatusInternalServerError,
			getAllSubscriptionErr: errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))

			if testCase.statusCode == http.StatusOK || testCase.statusCode == http.StatusInternalServerError {
				mockedStore.EXPECT().GetAllSubscriptions(testCase.userID).Return(testCase.subscriptionList, testCase.getAllSubscriptionErr)
			}

			if testCase.statusCode == http.StatusOK {
				mockedClient.EXPECT().DeleteSubscription(gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.statusCode, testCase.err)
				mockedStore.EXPECT().DeleteSubscription(gomock.Any()).Return(nil)
			}

			statusCode, _ := p.handleDeleteAllSubscriptions(testCase.userID, testCase.projectID)
			assert.Equal(t, testCase.statusCode, statusCode)
		})
	}
}

func TestHandleGetAllLinkedProjects(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, nil)
	for _, testCase := range []struct {
		description string
		projectList []serializers.ProjectDetails
		err         error
		statusCode  int
	}{
		{
			description: "HandleGetAllLinkedProjects: valid",
			projectList: []serializers.ProjectDetails{},
			statusCode:  http.StatusOK,
		},
		{
			description: "HandleGetAllLinkedProjects: error while fetching project list",
			err:         errors.New("mockError"),
			statusCode:  http.StatusInternalServerError,
		},
		{
			description: "HandleGetAllLinkedProjects: empty project list",
			statusCode:  http.StatusOK,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))

			mockedStore.EXPECT().GetAllProjects("mockMattermostUserID").Return(testCase.projectList, testCase.err)

			req := httptest.NewRequest(http.MethodGet, "/project/link", bytes.NewBufferString(`{}`))
			req.Header.Add(constants.HeaderMattermostUserID, "mockMattermostUserID")

			w := httptest.NewRecorder()
			p.handleGetAllLinkedProjects(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)
		})
	}
}

func TestHandleUnlinkProject(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, nil)
	for _, testCase := range []struct {
		description        string
		body               string
		err                error
		marshalError       error
		statusCode         int
		expectedStatusCode int
		projectList        []serializers.ProjectDetails
		project            serializers.ProjectDetails
	}{
		{
			description: "HandleUnlinkProject: valid",
			body: `{
				"organizationName": "mockOrganization",
				"projectName": "mockProject",
				"projectID" :"mockProjectID"
				}`,
			err:        nil,
			statusCode: http.StatusOK,
			projectList: []serializers.ProjectDetails{
				{
					MattermostUserID: "mockMattermostUserID",
					ProjectName:      "mockProject",
					OrganizationName: "mockOrganization",
					ProjectID:        "mockProjectID",
				},
			},
			project: serializers.ProjectDetails{
				MattermostUserID: "mockMattermostUserID",
				ProjectName:      "mockProject",
				OrganizationName: "mockOrganization",
				ProjectID:        "mockProjectID",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "HandleUnlinkProject: invalid body",
			body: `{
				"organizationName": "mockOrganization",
				"projectName": "mockProject",`,
			err:                errors.New("mockError"),
			statusCode:         http.StatusBadRequest,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "HandleUnlinkProject: missing fields",
			body: `{
				"organization": "mockOrganization",
				}`,
			err:                errors.New("mockError"),
			statusCode:         http.StatusBadRequest,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "HandleUnlinkProject: marshaling gives error",
			body: `{
				"organizationName": "mockOrganization",
				"projectName": "mockProject",
				"projectID" :"mockProjectID"
				}`,
			statusCode: http.StatusOK,
			projectList: []serializers.ProjectDetails{
				{
					MattermostUserID: "mockMattermostUserID",
					ProjectName:      "mockProject",
					OrganizationName: "mockOrganization",
					ProjectID:        "mockProjectID",
				},
			},
			project: serializers.ProjectDetails{
				MattermostUserID: "mockMattermostUserID",
				ProjectName:      "mockProject",
				OrganizationName: "mockOrganization",
				ProjectID:        "mockProjectID",
			},
			marshalError:       errors.New("mockError"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))

			monkey.PatchInstanceMethod(reflect.TypeOf(p), "IsProjectLinked", func(*Plugin, []serializers.ProjectDetails, serializers.ProjectDetails) (*serializers.ProjectDetails, bool) {
				return &serializers.ProjectDetails{}, true
			})

			if testCase.statusCode == http.StatusOK {
				mockedStore.EXPECT().GetAllProjects("mockMattermostUserID").Return(testCase.projectList, nil)
				mockedStore.EXPECT().DeleteProject(&testCase.project).Return(nil)
			}

			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})

			req := httptest.NewRequest(http.MethodPost, "/project/unlink", bytes.NewBufferString(testCase.body))
			req.Header.Add(constants.HeaderMattermostUserID, "mockMattermostUserID")

			w := httptest.NewRecorder()
			p.handleUnlinkProject(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestHandleGetUserAccountDetails(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, nil)
	for _, testCase := range []struct {
		description   string
		err           error
		marshalError  error
		statusCode    int
		user          *serializers.User
		loadUserError error
	}{
		{
			description: "HandleGetUserAccountDetails: valid",
			statusCode:  http.StatusOK,
			user: &serializers.User{
				MattermostUserID: "mockMattermostUserID",
			},
		},
		{
			description: "HandleGetUserAccountDetails: empty user details",
			err:         nil,
			statusCode:  http.StatusUnauthorized,
			user:        &serializers.User{},
		},
		{
			description:   "HandleGetUserAccountDetails: error while loading user",
			loadUserError: errors.New("mockError"),
			statusCode:    http.StatusInternalServerError,
		},
		{
			description: "HandleGetUserAccountDetails: marshaling gives error",
			statusCode:  http.StatusInternalServerError,
			user: &serializers.User{
				MattermostUserID: "mockMattermostUserID",
			},
			marshalError: errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			mockAPI.On("PublishWebSocketEvent", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("*model.WebsocketBroadcast")).Return(nil)

			mockedStore.EXPECT().LoadUser("mockMattermostUserID").Return(testCase.user, testCase.loadUserError)

			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})

			req := httptest.NewRequest(http.MethodGet, "/user", bytes.NewBufferString(`{}`))
			req.Header.Add(constants.HeaderMattermostUserID, "mockMattermostUserID")

			w := httptest.NewRecorder()
			p.handleGetUserAccountDetails(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)
		})
	}
}

func TestHandleCreateSubscriptions(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, mockedClient)
	for _, testCase := range []struct {
		description        string
		body               string
		err                error
		marshalError       error
		expectedStatusCode int
		statusCode         int
		projectList        []serializers.ProjectDetails
		project            serializers.ProjectDetails
		subscriptionList   []*serializers.SubscriptionDetails
		subscription       *serializers.SubscriptionDetails
		isProjectLinked    bool
	}{
		{
			description: "HandleCreateSubscriptions: valid",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",
				"eventType": "mockEventType",
				"serviceType": "mockServiceType",
				"channelID": "mockChannelID"
				}`,
			statusCode:         http.StatusOK,
			expectedStatusCode: http.StatusOK,
			projectList:        []serializers.ProjectDetails{},
			project:            serializers.ProjectDetails{},
			subscriptionList:   []*serializers.SubscriptionDetails{},
			subscription: &serializers.SubscriptionDetails{
				MattermostUserID: "mockMattermostUserID",
				ProjectName:      "mockProject",
				OrganizationName: "mockOrganization",
				EventType:        "mockEventType",
				ServiceType:      "mockServiceType",
				ChannelID:        "mockChannelID",
			},
		},
		{
			description:        "HandleCreateSubscriptions: empty body",
			body:               `{}`,
			err:                errors.New("mockError"),
			statusCode:         http.StatusBadRequest,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "HandleCreateSubscriptions: invalid body",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",`,
			err:                errors.New("mockError"),
			statusCode:         http.StatusBadRequest,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "HandleCreateSubscriptions: missing fields",
			body: `{
				"organization": "mockOrganization",
				}`,
			err:                errors.New("mockError"),
			statusCode:         http.StatusBadRequest,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "HandleCreateSubscriptions: marshaling gives error",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",
				"eventType": "mockEventType",
				"serviceType": "mockServiceType",
				"channelID": "mockChannelID"
				}`,
			statusCode:         http.StatusOK,
			marshalError:       errors.New("mockError"),
			expectedStatusCode: http.StatusInternalServerError,
			projectList:        []serializers.ProjectDetails{},
			project:            serializers.ProjectDetails{},
			subscriptionList:   []*serializers.SubscriptionDetails{},
			subscription: &serializers.SubscriptionDetails{
				MattermostUserID: "mockMattermostUserID",
				ProjectName:      "mockProject",
				OrganizationName: "mockOrganization",
				EventType:        "mockEventType",
				ServiceType:      "mockServiceType",
				ChannelID:        "mockChannelID",
			},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			mockAPI.On("GetChannel", mock.AnythingOfType("string")).Return(&model.Channel{}, nil)
			mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{}, nil)

			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(p), "IsProjectLinked", func(*Plugin, []serializers.ProjectDetails, serializers.ProjectDetails) (*serializers.ProjectDetails, bool) {
				return &serializers.ProjectDetails{}, true
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(p), "IsSubscriptionPresent", func(*Plugin, []*serializers.SubscriptionDetails, *serializers.SubscriptionDetails) (*serializers.SubscriptionDetails, bool) {
				return &serializers.SubscriptionDetails{}, false
			})

			if testCase.statusCode == http.StatusOK {
				mockedClient.EXPECT().CreateSubscription(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&serializers.SubscriptionValue{}, testCase.statusCode, testCase.err)
				mockedStore.EXPECT().GetAllProjects("mockMattermostUserID").Return(testCase.projectList, nil)
				mockedStore.EXPECT().GetAllSubscriptions("mockMattermostUserID").Return(testCase.subscriptionList, nil)
				mockedStore.EXPECT().StoreSubscription(testCase.subscription).Return(nil)
			}

			req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBufferString(testCase.body))
			req.Header.Add(constants.HeaderMattermostUserID, "mockMattermostUserID")

			w := httptest.NewRecorder()
			p.handleCreateSubscription(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestHandleGetSubscriptions(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, nil)
	for _, testCase := range []struct {
		description                                          string
		subscriptionList                                     []*serializers.SubscriptionDetails
		project                                              string
		err                                                  error
		marshalError                                         error
		GetSubscriptionsForAccessibleChannelsOrProjectsError error
		statusCode                                           int
		isTeamIDValid                                        bool
	}{
		{
			description:      "HandleGetSubscriptions: valid",
			subscriptionList: []*serializers.SubscriptionDetails{},
			statusCode:       http.StatusOK,
			isTeamIDValid:    true,
		},
		{
			description:   "HandleGetSubscriptions: project as a query param",
			project:       "mockProject",
			statusCode:    http.StatusOK,
			isTeamIDValid: true,
		},
		{
			description:   "HandleGetSubscriptions: error while fetching subscription list",
			err:           errors.New("mockError"),
			statusCode:    http.StatusInternalServerError,
			isTeamIDValid: true,
		},
		{
			description:   "HandleGetSubscriptions: empty subscription list",
			statusCode:    http.StatusOK,
			isTeamIDValid: true,
		},
		{
			description:   "HandleGetSubscriptions: marshaling gives error",
			marshalError:  errors.New("mockError"),
			statusCode:    http.StatusInternalServerError,
			isTeamIDValid: true,
		},
		{
			description: "HandleGetSubscriptions: GetSubscriptionsForAccessibleChannelsOrProjects gives error",
			project:     "mockProject",
			GetSubscriptionsForAccessibleChannelsOrProjectsError: errors.New("mockError"),
			statusCode:    http.StatusInternalServerError,
			isTeamIDValid: true,
		},
		{
			description:   "HandleGetSubscriptions: Team ID is invalid",
			statusCode:    http.StatusBadRequest,
			isTeamIDValid: false,
		},
		{
			description: "HandleGetSubscriptions: GetSubscriptionsForAccessibleChannelsOrProjects gives error",
			project:     "mockProject",
			GetSubscriptionsForAccessibleChannelsOrProjectsError: errors.New("mockError"),
			statusCode:    http.StatusInternalServerError,
			isTeamIDValid: true,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))

			if testCase.isTeamIDValid {
				mockedStore.EXPECT().GetAllSubscriptions("mockMattermostUserID").Return(testCase.subscriptionList, testCase.err)
			}

			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})

			monkey.Patch(model.IsValidId, func(_ string) bool {
				return testCase.isTeamIDValid
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(p), "GetSubscriptionsForAccessibleChannelsOrProjects", func(_ *Plugin, _ []*serializers.SubscriptionDetails, _, _, _ string) ([]*serializers.SubscriptionDetails, error) {
				return nil, testCase.GetSubscriptionsForAccessibleChannelsOrProjectsError
			})

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?project=%s", "/subscriptions", testCase.project), bytes.NewBufferString(`{}`))
			req.Header.Add(constants.HeaderMattermostUserID, "mockMattermostUserID")

			w := httptest.NewRecorder()
			p.handleGetSubscriptions(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)
		})
	}
}

func TestHandleSubscriptionNotifications(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupMockPlugin(mockAPI, nil, nil)
	for _, testCase := range []struct {
		description      string
		body             string
		channelID        string
		isValidChannelID bool
		err              error
		statusCode       int
	}{
		{
			description: "SubscriptionNotifications: valid",
			body: `{
				"detailedMessage": {
					"markdown": "mockMarkdown"
					}
				}`,
			channelID:        "mockChannelID",
			isValidChannelID: true,
			statusCode:       http.StatusOK,
		},
		{
			description:      "SubscriptionNotifications: empty body",
			body:             `{}`,
			err:              errors.New("mockError"),
			channelID:        "mockChannelID",
			isValidChannelID: true,
			statusCode:       http.StatusOK,
		},
		{
			description: "SubscriptionNotifications: invalid channel ID",
			body:        `{}`,
			err:         errors.New("mockError"),
			channelID:   "mockChannelID",
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "SubscriptionNotifications: invalid body",
			body: `{
				"detailedMessage": {
					"markdown": "mockMarkdown"`,
			err:              errors.New("mockError"),
			channelID:        "mockChannelID",
			isValidChannelID: true,
			statusCode:       http.StatusBadRequest,
		},
		{
			description: "SubscriptionNotifications: without channelID",
			body: `{
				"detailedMessage": {
					"markdown": "mockMarkdown"
					}
				}`,
			isValidChannelID: true,
			statusCode:       http.StatusBadRequest,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			mockAPI.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)

			monkey.Patch(model.IsValidId, func(string) bool {
				return testCase.isValidChannelID
			})

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/notification?channelID=%s", testCase.channelID), bytes.NewBufferString(testCase.body))

			w := httptest.NewRecorder()
			p.handleSubscriptionNotifications(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)
		})
	}
}

func TestHandleDeleteSubscriptions(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, mockedClient)
	for _, testCase := range []struct {
		description      string
		body             string
		err              error
		statusCode       int
		subscriptionList []*serializers.SubscriptionDetails
		subscription     *serializers.SubscriptionDetails
	}{
		{
			description: "HandleDeleteSubscriptions: valid",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",
				"eventType": "mockEventType",
				"channelID": "mockChannelID",
				"mmUserID": "mockMattermostUserID"
				}`,
			statusCode:       http.StatusOK,
			subscriptionList: []*serializers.SubscriptionDetails{},
			subscription: &serializers.SubscriptionDetails{
				MattermostUserID: "mockMattermostUserID",
				ProjectName:      "mockProject",
				OrganizationName: "mockOrganization",
				EventType:        "mockEventType",
				ChannelID:        "mockChannelID",
			},
		},
		{
			description: "HandleDeleteSubscriptions: empty body",
			body:        `{}`,
			err:         errors.New("mockError"),
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "HandleDeleteSubscriptions: invalid body",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject",`,
			err:        errors.New("mockError"),
			statusCode: http.StatusBadRequest,
		},
		{
			description: "HandleDeleteSubscriptions: missing fields",
			body: `{
				"organization": "mockOrganization",
				}`,
			err:        errors.New("mockError"),
			statusCode: http.StatusBadRequest,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			mockAPI.On("LogDebug", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))

			monkey.PatchInstanceMethod(reflect.TypeOf(p), "IsSubscriptionPresent", func(*Plugin, []*serializers.SubscriptionDetails, *serializers.SubscriptionDetails) (*serializers.SubscriptionDetails, bool) {
				return &serializers.SubscriptionDetails{}, true
			})

			if testCase.statusCode == http.StatusOK {
				mockedClient.EXPECT().DeleteSubscription(gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.statusCode, testCase.err)
				mockedStore.EXPECT().GetAllSubscriptions("mockMattermostUserID").Return(testCase.subscriptionList, nil)
				mockedStore.EXPECT().DeleteSubscription(gomock.Any()).Return(nil)
			}

			req := httptest.NewRequest(http.MethodDelete, "/subscriptions", bytes.NewBufferString(testCase.body))
			req.Header.Add(constants.HeaderMattermostUserID, "mockMattermostUserID")

			w := httptest.NewRecorder()
			p.handleDeleteSubscriptions(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)
		})
	}
}

func TestGetUserChannelsForTeam(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	p := setupMockPlugin(mockAPI, nil, nil)
	for _, testCase := range []struct {
		description string
		teamID      string
		channels    []*model.Channel
		channelErr  *model.AppError
		statusCode  int
	}{
		{
			description: "GetUserChannelsForTeam: valid",
			teamID:      "qteks46as3befxj4ec1mip5ume",
			channels: []*model.Channel{
				{
					Id:   "mockChannelID",
					Type: model.CHANNEL_OPEN,
				},
			},
			channelErr: nil,
			statusCode: http.StatusOK,
		},
		{
			description: "GetUserChannelsForTeam: no channels",
			teamID:      "qteks46as3befxj4ec1mip5ume",
			channels:    nil,
			channelErr:  nil,
			statusCode:  http.StatusOK,
		},
		{
			description: "GetUserChannelsForTeam: invalid teamID",
			teamID:      "invalid-teamID",
			channelErr:  nil,
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "GetUserChannelsForTeam: no required channels",
			teamID:      "qteks46as3befxj4ec1mip5ume",
			channels: []*model.Channel{
				{
					Id:   "mockChannelID",
					Type: model.CHANNEL_PRIVATE,
				},
			},
			channelErr: nil,
			statusCode: http.StatusOK,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			mockAPI.On("GetChannelsForTeamForUser", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("bool")).Return(testCase.channels, testCase.channelErr)

			req := httptest.NewRequest(http.MethodGet, "/channels", bytes.NewBufferString(`{}`))
			req.Header.Add(constants.HeaderMattermostUserID, "test-userID")

			pathParams := map[string]string{
				"team_id": testCase.teamID,
			}

			req = mux.SetURLVars(req, pathParams)

			w := httptest.NewRecorder()
			p.getUserChannelsForTeam(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)
		})
	}
}

func TestHandleGetSubscriptionFilterPossibleValues(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(mockAPI, nil, mockedClient)
	for _, testCase := range []struct {
		description                            string
		body                                   string
		getSubscriptionFilterPossibleValuesErr error
		statusCode                             int
		getGitRepositoryBranchesResponse       *serializers.SubscriptionFilterPossibleValuesResponseFromClient
		expectedResponse                       string
		expectedErrorResponse                  interface{}
	}{
		{
			description: "HandleGetSubscriptionFilterPossibleValues: valid",
			body: `{
				"organization": "mockOrganization",
				"projectId": "mockProjectID",
				"eventType": "mockEventType",
				"repositoryId": "mockRepositoryID",
				"filters": ["mockFilter1", "mockFilter2"]
				}`,
			statusCode: http.StatusOK,
			getGitRepositoryBranchesResponse: &serializers.SubscriptionFilterPossibleValuesResponseFromClient{
				InputValues: []*serializers.InputValues{
					{
						PossibleValues: []*serializers.PossibleValues{},
						SubscriptionFilter: serializers.SubscriptionFilter{
							InputID: "mockInputID1",
						},
					},
					{
						PossibleValues: []*serializers.PossibleValues{},
						SubscriptionFilter: serializers.SubscriptionFilter{
							InputID: "mockInputID2",
						},
					},
				},
			},
			expectedResponse: `{"mockInputID1":[],"mockInputID2":[]}`,
		},
		{
			description: "HandleGetSubscriptionFilterPossibleValues: missing fields",
			body: `{
				"projectId": "mockProjectID",
				"eventType": "mockEventType",
				"repositoryId": "mockRepositoryID",
				"filters": ["mockFilter1", "mockFilter2"]
				}`,
			statusCode:            http.StatusBadRequest,
			expectedErrorResponse: map[string]interface{}{"Error": constants.OrganizationRequired},
		},
		{
			description: "HandleGetSubscriptionFilterPossibleValues: Error fetching subscription filter possible values",
			body: `{
				"organization": "mockOrganization",
				"projectId": "mockProjectID",
				"eventType": "mockEventType",
				"repositoryId": "mockRepositoryID",
				"filters": ["mockFilter1", "mockFilter2"]
				}`,
			statusCode:                             http.StatusInternalServerError,
			getSubscriptionFilterPossibleValuesErr: errors.New("failed to fetch subscription filters possible values"),
			expectedErrorResponse:                  map[string]interface{}{"Error": "failed to fetch subscription filters possible values"},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", testutils.GetMockArgumentsWithType("string", 3)...)

			if testCase.statusCode == http.StatusOK || testCase.statusCode == http.StatusInternalServerError {
				mockedClient.EXPECT().GetSubscriptionFilterPossibleValues(gomock.Any(), gomock.Any()).Return(testCase.getGitRepositoryBranchesResponse, testCase.statusCode, testCase.getSubscriptionFilterPossibleValuesErr)
			}

			req := httptest.NewRequest(http.MethodGet, "/mockPath", bytes.NewBufferString(testCase.body))
			req.Header.Add(constants.HeaderMattermostUserID, "mockUserID")

			w := httptest.NewRecorder()
			p.handleGetSubscriptionFilterPossibleValues(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)

			if testCase.expectedErrorResponse != nil {
				var actualResponse interface{}
				err := json.NewDecoder(resp.Body).Decode(&actualResponse)
				require.Nil(t, err)
				assert.Equal(t, testCase.expectedErrorResponse, actualResponse)
			}

			if testCase.expectedResponse != "" {
				response, err := ioutil.ReadAll(resp.Body)
				require.Nil(t, err)
				assert.Contains(t, string(response), testCase.expectedResponse)
			}
		})
	}
}
