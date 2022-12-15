package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p := setupMockPlugin(mockAPI, mockedStore, mockedClient)
	for _, testCase := range []struct {
		description string
		body        string
		err         error
		statusCode  int
		projectList []serializers.ProjectDetails
		project     serializers.ProjectDetails
	}{
		{
			description: "HandleLink: valid",
			body: `{
				"organization": "mockOrganization",
				"project": "mockProject"
				}`,
			err:        nil,
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
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			mockAPI.On("GetDirectChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&model.Channel{}, nil)
			mockAPI.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)

			if testCase.statusCode == http.StatusOK {
				mockedClient.EXPECT().Link(gomock.Any(), gomock.Any()).Return(&serializers.Project{}, testCase.statusCode, testCase.err)
				mockedClient.EXPECT().CheckIfUserIsProjectAdmin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
				mockedStore.EXPECT().GetAllProjects("mockMattermostUserID").Return(testCase.projectList, nil)
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

			monkey.PatchInstanceMethod(reflect.TypeOf(p), "GetSubscriptionsForAccessibleChannelsOrProjects", func(_ *Plugin, _ []*serializers.SubscriptionDetails, _, _ string) ([]*serializers.SubscriptionDetails, error) {
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

func TestHandleGetGitRepositories(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(mockAPI, nil, mockedClient)
	for _, testCase := range []struct {
		description           string
		organization          string
		project               string
		getGitRepositoriesErr error
		statusCode            int
	}{
		{
			description:  "HandleGetGitRepositories: valid",
			organization: "mockOrganization",
			project:      "mockProject",
			statusCode:   http.StatusOK,
		},
		{
			description:  "HandleGetGitRepositories: Invalid organization or project name",
			organization: "mockOrganization",
			statusCode:   http.StatusBadRequest,
		},
		{
			description:           "HandleGetGitRepositories: GetGitRepositories returns error",
			organization:          "mockOrganization",
			project:               "mockProject",
			getGitRepositoriesErr: errors.New("failed to get git repository branches"),
			statusCode:            http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", testutils.GetMockArgumentsWithType("string", 3)...)

			if testCase.statusCode == http.StatusInternalServerError || testCase.statusCode == http.StatusOK {
				mockedClient.EXPECT().GetGitRepositories(gomock.Any(), gomock.Any(), gomock.Any()).Return(&serializers.GitRepositoriesResponse{}, testCase.statusCode, testCase.getGitRepositoriesErr)
			}

			req := httptest.NewRequest(http.MethodGet, "/mockPath", bytes.NewBufferString(`{}`))
			req.Header.Add(constants.HeaderMattermostUserID, "test-userID")

			pathParams := map[string]string{
				"organization": testCase.organization,
				"project":      testCase.project,
			}

			req = mux.SetURLVars(req, pathParams)

			w := httptest.NewRecorder()
			p.handleGetGitRepositories(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)
		})
	}
}

func TestHandleGetGitRepositoryBranches(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(mockAPI, nil, mockedClient)
	for _, testCase := range []struct {
		description                      string
		organization                     string
		repository                       string
		project                          string
		getGitRepositoryBranchesErr      error
		statusCode                       int
		getGitRepositoryBranchesResponse *serializers.GitBranchesResponse
		expectedResponse                 []*serializers.GitBranch
		expectedErrorResponse            interface{}
	}{
		{
			description:  "HandleGetGitRepositoryBranches: valid",
			organization: "mockOrganization",
			project:      "mockProject",
			repository:   "mockRepository",
			statusCode:   http.StatusOK,
			getGitRepositoryBranchesResponse: &serializers.GitBranchesResponse{
				Value: testutils.GetGitBranchesPayload(),
			},
			expectedResponse: []*serializers.GitBranch{
				{
					ID:   "mockID-1",
					Name: "mockName-1",
				},
				{
					ID:   "mockID-2",
					Name: "mockName-2",
				},
			},
		},
		{
			description:           "HandleGetGitRepositoryBranches: Invalid organization, project or repository params",
			organization:          "mockOrganization",
			project:               "mockProject",
			statusCode:            http.StatusBadRequest,
			expectedErrorResponse: map[string]interface{}{"Error": constants.ErrorRepositoryPathParam},
		},
		{
			description:                 "HandleGetGitRepositoryBranches: GetGitRepositoryBranches returns error",
			organization:                "mockOrganization",
			project:                     "mockProject",
			repository:                  "mockRepository",
			getGitRepositoryBranchesErr: errors.New("failed to get the git repository branches for a project"),
			statusCode:                  http.StatusInternalServerError,
			expectedErrorResponse:       map[string]interface{}{"Error": "failed to get the git repository branches for a project"},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", testutils.GetMockArgumentsWithType("string", 3)...)

			if testCase.statusCode == http.StatusInternalServerError || testCase.statusCode == http.StatusOK {
				mockedClient.EXPECT().GetGitRepositoryBranches(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.getGitRepositoryBranchesResponse, testCase.statusCode, testCase.getGitRepositoryBranchesErr)
			}

			req := httptest.NewRequest(http.MethodGet, "/mockPath", bytes.NewBufferString(`{}`))
			req.Header.Add(constants.HeaderMattermostUserID, "mockUserID")

			pathParams := map[string]string{
				"organization": testCase.organization,
				"project":      testCase.project,
				"repository":   testCase.repository,
			}

			req = mux.SetURLVars(req, pathParams)

			w := httptest.NewRecorder()
			p.handleGetGitRepositoryBranches(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)

			if testCase.expectedErrorResponse != nil {
				var actualResponse interface{}
				err := json.NewDecoder(resp.Body).Decode(&actualResponse)
				require.Nil(t, err)
				assert.Equal(t, testCase.expectedErrorResponse, actualResponse)
			}

			if testCase.expectedResponse != nil {
				var actualResponse []*serializers.GitBranch
				err := json.NewDecoder(resp.Body).Decode(&actualResponse)
				require.Nil(t, err)
				assert.Equal(t, testCase.expectedResponse, actualResponse)
			}
		})
	}
}

func TestHandlePipelineApproveOrRejectRequest(t *testing.T) {
	defer monkey.UnpatchAll()
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p := setupMockPlugin(mockAPI, nil, mockedClient)
	for _, testCase := range []struct {
		description                               string
		body                                      string
		statusCode                                int
		updatePipelineReleaseApprovalPostError    error
		updatePipelineReleaseApprovalRequestError error
		getApprovalDetailsError                   error
		updatePipelineApprovalRequestStatus       int
		getApprovalDetailsStatus                  int
		isPayloadInvalid                          bool
	}{
		{
			description: "HandlePipelineApproveOrRejectRequest: valid",
			body: `{
				"post_id": "mockPostID",
				"channel_id": "mockChannelID",
				"context": {
				  "approvalId": "mockApproverID",
				  "organization": "mockOrganization",
				  "projectID": "mockProjectID",
				  "requestType": "mockRequestType"
				}
			  }`,
			updatePipelineApprovalRequestStatus: http.StatusOK,
			statusCode:                          http.StatusOK,
			getApprovalDetailsStatus:            http.StatusOK,
		},
		{
			description: "HandlePipelineApproveOrRejectRequest: approved/rejected request successfully but failed to update post",
			body: `{
				"post_id": "mockPostID",
				"channel_id": "mockChannelID",
				"context": {
				  "approvalId": "mockApproverID",
				  "organization": "mockOrganization",
				  "projectID": "mockProjectID",
				  "requestType": "mockRequestType"
				}
			  }`,
			updatePipelineApprovalRequestStatus:    http.StatusOK,
			updatePipelineReleaseApprovalPostError: errors.New("approved/rejected request successfully but failed to update post"),
			statusCode:                             http.StatusInternalServerError,
		},
		{
			description: "HandlePipelineApproveOrRejectRequest: failed to approve/reject request",
			body: `{
				"post_id": "mockPostID",
				"channel_id": "mockChannelID",
				"context": {
				  "approvalId": "mockApproverID",
				  "organization": "mockOrganization",
				  "projectID": "mockProjectID",
				  "requestType": "mockRequestType"
				}
			  }`,
			updatePipelineApprovalRequestStatus:       http.StatusInternalServerError,
			updatePipelineReleaseApprovalRequestError: errors.New("not permitted to complete approval"),
			statusCode: http.StatusOK,
		},
		{
			description: "HandlePipelineApproveOrRejectRequest: failed to approve/reject request and update the post",
			body: `{
				"post_id": "mockPostID",
				"channel_id": "mockChannelID",
				"context": {
				  "approvalId": "mockApproverID",
				  "organization": "mockOrganization",
				  "projectID": "mockProjectID",
				  "requestType": "mockRequestType"
				}
			  }`,
			updatePipelineApprovalRequestStatus:       http.StatusInternalServerError,
			updatePipelineReleaseApprovalRequestError: errors.New("not permitted to complete approval"),
			updatePipelineReleaseApprovalPostError:    errors.New("failed to approve/reject request and update the post"),
			statusCode:                                http.StatusInternalServerError,
		},
		{
			description: "HandlePipelineApproveOrRejectRequest: failed to approve/reject request and fetch approval details",
			body: `{
				"post_id": "mockPostID",
				"channel_id": "mockChannelID",
				"context": {
				  "approvalId": "mockApproverID",
				  "organization": "mockOrganization",
				  "projectID": "mockProjectID",
				  "requestType": "mockRequestType"
				}
			  }`,
			updatePipelineApprovalRequestStatus:       http.StatusInternalServerError,
			updatePipelineReleaseApprovalRequestError: errors.New("not permitted to complete approval"),
			getApprovalDetailsError:                   errors.New("failed to approve/reject request and fetch approval details"),
			statusCode:                                http.StatusInternalServerError,
			getApprovalDetailsStatus:                  http.StatusInternalServerError,
		},
		{
			description: "HandlePipelineApproveOrRejectRequest: invalid payload",
			body: `{
				"post_id": "mockPostID",
				"channel_id": "mockChannelID",
				"context": {
    			  "approvalId": "mockApproverID",
    			  "wrong: "mockOrganization",
				  "projectID": "mockProjectID",
				  "requestType": "mockRequestType"
				}
			  }`,
			isPayloadInvalid: true,
			statusCode:       http.StatusInternalServerError,
		},
		{
			description: "HandlePipelineApproveOrRejectRequest: failed to approve/reject request due to some internal server error",
			body: `{
				"post_id": "mockPostID",
				"channel_id": "mockChannelID",
				"context": {
				  "approvalId": "mockApproverID",
				  "organization": "mockOrganization",
				  "projectID": "mockProjectID",
				  "requestType": "mockRequestType"
				}
			  }`,
			updatePipelineApprovalRequestStatus:       http.StatusBadRequest,
			statusCode:                                http.StatusInternalServerError,
			updatePipelineReleaseApprovalRequestError: errors.New("failed to approve/reject request due to some internal server error"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", testutils.GetMockArgumentsWithType("string", 3)...)
			mockAPI.On("GetDirectChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&model.Channel{}, nil)
			mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Return(&model.Post{Message: "mockMessage"})
			mockAPI.On("UpdateEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Return(nil)

			if !testCase.isPayloadInvalid {
				mockedClient.EXPECT().UpdatePipelineRunApprovalRequest(gomock.Any(), "mockOrganization", "mockProjectID", "test-userID").Return(&serializers.PipelineRunApproveResponse{
					Value: []*serializers.PipelineRunResponseValue{
						{},
					},
				}, testCase.updatePipelineApprovalRequestStatus, testCase.updatePipelineReleaseApprovalRequestError)
			}

			if testCase.updatePipelineApprovalRequestStatus == http.StatusInternalServerError {
				mockedClient.EXPECT().GetRunApprovalDetails("mockOrganization", "mockProjectID", "test-userID", "mockApproverID").Return(&serializers.PipelineRunApprovalDetails{}, testCase.getApprovalDetailsStatus, testCase.getApprovalDetailsError)
			}

			monkey.PatchInstanceMethod(reflect.TypeOf(p), "UpdatePipelineRunApprovalPost", func(_ *Plugin, _ []*serializers.ApprovalSteps, _ int, _, _, _ string) error {
				return testCase.updatePipelineReleaseApprovalPostError
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(p), "DM", func(_ *Plugin, _, _ string, _ bool, _ ...interface{}) (string, error) {
				return "", nil
			})

			req := httptest.NewRequest(http.MethodGet, "/mockPath", bytes.NewBufferString(testCase.body))
			req.Header.Add(constants.HeaderMattermostUserID, "test-userID")

			w := httptest.NewRecorder()
			p.handlePipelineApproveOrRejectRunRequest(w, req)
			resp := w.Result()
			assert.Equal(t, testCase.statusCode, resp.StatusCode)
		})
	}
}
