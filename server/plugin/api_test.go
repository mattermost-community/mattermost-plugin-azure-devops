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
	"github.com/Brightscout/mattermost-plugin-azure-devops/mocks"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type panicHandler struct {
}

func (ph panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("bad handler")
}

func TestInitRoutes(t *testing.T) {
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
	p.router = p.InitAPI()
	p.InitRoutes()
}

func TestHandleStaticFiles(t *testing.T) {
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI

	p.router = p.InitAPI()
	mockAPI.On("GetBundlePath").Return("/test-path", nil)
	p.HandleStaticFiles()
}

func TestWithRecovery(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			require.Fail(t, "got panic")
		}
	}()

	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockAPI.On("LogError", "Recovered from a panic", "url", "http://random", "error", "bad handler", "stack", mock.Anything)
	p.SetAPI(mockAPI)

	ph := panicHandler{}
	handler := p.WithRecovery(ph)

	req := httptest.NewRequest(http.MethodGet, "http://random", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.Body != nil {
		defer resp.Body.Close()
		_, err := io.Copy(ioutil.Discard, resp.Body)
		require.NoError(t, err)
	}
}

func TestHandleAuthRequired(t *testing.T) {
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	p.API = mockAPI
	p.Client = mockedClient
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Client = mockedClient
	p.Store = mockedStore
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
				mockedStore.EXPECT().GetAllProjects("mockMattermostUserID").Return(testCase.projectList, nil)
				mockedStore.EXPECT().StoreProject(&serializers.ProjectDetails{
					MattermostUserID: "mockMattermostUserID",
					OrganizationName: "mockOrganization",
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Store = mockedStore
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Store = mockedStore
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

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "IsProjectLinked", func(*Plugin, []serializers.ProjectDetails, serializers.ProjectDetails) (*serializers.ProjectDetails, bool) {
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Store = mockedStore
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Client = mockedClient
	p.Store = mockedStore
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
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "IsProjectLinked", func(*Plugin, []serializers.ProjectDetails, serializers.ProjectDetails) (*serializers.ProjectDetails, bool) {
				return &serializers.ProjectDetails{}, true
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "IsSubscriptionPresent", func(*Plugin, []*serializers.SubscriptionDetails, *serializers.SubscriptionDetails) (*serializers.SubscriptionDetails, bool) {
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Store = mockedStore
	for _, testCase := range []struct {
		description      string
		subscriptionList []*serializers.SubscriptionDetails
		project          string
		err              error
		marshalError     error
		statusCode       int
	}{
		{
			description:      "HandleGetSubscriptions: valid",
			subscriptionList: []*serializers.SubscriptionDetails{},
			statusCode:       http.StatusOK,
		},
		{
			description: "HandleGetSubscriptions: project as a query param",
			project:     "mockProject",
			statusCode:  http.StatusOK,
		},
		{
			description: "HandleGetSubscriptions: error while fetching subscription list",
			err:         errors.New("mockError"),
			statusCode:  http.StatusInternalServerError,
		},
		{
			description: "HandleGetSubscriptions: empty subscription list",
			statusCode:  http.StatusOK,
		},
		{
			description:  "HandleGetSubscriptions: marshaling gives error",
			marshalError: errors.New("mockError"),
			statusCode:   http.StatusInternalServerError,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))

			mockedStore.EXPECT().GetAllSubscriptions("mockMattermostUserID").Return(testCase.subscriptionList, testCase.err)

			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	mockCtrl := gomock.NewController(t)
	mockedClient := mocks.NewMockClient(mockCtrl)
	mockedStore := mocks.NewMockKVStore(mockCtrl)
	p.API = mockAPI
	p.Client = mockedClient
	p.Store = mockedStore
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

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "IsSubscriptionPresent", func(*Plugin, []*serializers.SubscriptionDetails, *serializers.SubscriptionDetails) (*serializers.SubscriptionDetails, bool) {
				return &serializers.SubscriptionDetails{}, true
			})

			if testCase.statusCode == http.StatusOK {
				mockedClient.EXPECT().DeleteSubscription(gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.statusCode, testCase.err)
				mockedStore.EXPECT().GetAllSubscriptions("mockMattermostUserID").Return(testCase.subscriptionList, nil)
				mockedStore.EXPECT().DeleteSubscription(testCase.subscription).Return(nil)
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
	p := Plugin{}
	mockAPI := &plugintest.API{}
	p.API = mockAPI
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
