package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
)

// Initializes the plugin REST API
func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.Use(p.WithRecovery)

	// 404 handler
	r.Handle(constants.WildRoute, http.NotFoundHandler())
	return r
}

// Add custom routes and corresponding handlers here
func (p *Plugin) InitRoutes() {
	p.Client = InitClient(p)

	s := p.router.PathPrefix(constants.APIPrefix).Subrouter()

	// OAuth
	s.HandleFunc(constants.PathOAuthConnect, p.OAuthConnect).Methods(http.MethodGet)
	s.HandleFunc(constants.PathOAuthCallback, p.OAuthComplete).Methods(http.MethodGet)
	s.HandleFunc(constants.PathCreateTasks, p.handleAuthRequired(p.checkOAuth(p.handleCreateTask))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathLinkProject, p.handleAuthRequired(p.checkOAuth(p.handleLink))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathGetAllLinkedProjects, p.handleAuthRequired(p.checkOAuth(p.handleGetAllLinkedProjects))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathUnlinkProject, p.handleAuthRequired(p.checkOAuth(p.handleUnlinkProject))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathUser, p.handleAuthRequired(p.checkOAuth(p.handleGetUserAccountDetails))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathGetSubscriptions, p.handleAuthRequired(p.checkOAuth(p.handleGetSubscriptions))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathCreateSubscriptions, p.handleAuthRequired(p.checkOAuth(p.handleCreateSubscriptions))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathNotificationSubscriptions, p.handleNotificationSubscriptions).Methods(http.MethodPost)
	s.HandleFunc(constants.PathDeleteSubscriptions, p.handleAuthRequired(p.checkOAuth(p.handleDeleteSubscriptions))).Methods(http.MethodDelete)
	// TODO: for testing purpose, remove later
	s.HandleFunc("/test", p.testAPI).Methods(http.MethodGet)
}

// handleAuthRequired verifies if the provided request is performed by an authorized source.
func (p *Plugin) handleAuthRequired(handleFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
		if mattermostUserID == "" {
			error := serializers.Error{Code: http.StatusUnauthorized, Message: constants.NotAuthorized}
			p.handleError(w, r, &error)
			return
		}

		handleFunc(w, r)
	}
}

func (p *Plugin) checkOAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
		user, err := p.Store.LoadUser(mattermostUserID)
		if err != nil || user.AccessToken == "" {
			if errors.Is(err, ErrNotFound) || user.AccessToken == "" {
				p.handleError(w, r, &serializers.Error{Code: http.StatusUnauthorized, Message: constants.ConnectAccountFirst})
			} else {
				p.API.LogError("Unable to get user", "Error", err.Error())
				p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: constants.GenericErrorMessage})
			}
			return
		}
		handler(w, r)
	}
}

func (p *Plugin) handleError(w http.ResponseWriter, r *http.Request, error *serializers.Error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(error.Code)
	message := map[string]string{constants.Error: error.Message}
	response, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// API to link a project and an organization to a user.
func (p *Plugin) handleLink(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	var body *serializers.LinkRequestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		p.API.LogError(constants.ErrorDecodingBody, "Error", err.Error())
		error := serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		p.handleError(w, r, &error)
		return
	}

	if err := body.IsLinkPayloadValid(); err != "" {
		error := serializers.Error{Code: http.StatusBadRequest, Message: err}
		p.handleError(w, r, &error)
		return
	}

	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	_, isProjectLinked := p.IsProjectLinked(projectList, serializers.ProjectDetails{OrganizationName: body.Organization, ProjectName: body.Project})
	if isProjectLinked {
		p.DM(mattermostUserID, constants.AlreadyLinkedProject)
		return
	}

	response, statusCode, err := p.Client.Link(body, mattermostUserID)
	if err != nil {
		error := serializers.Error{Code: statusCode, Message: err.Error()}
		p.handleError(w, r, &error)
		return
	}

	project := serializers.ProjectDetails{
		MattermostUserID: mattermostUserID,
		ProjectID:        response.ID,
		ProjectName:      response.Name,
		OrganizationName: body.Organization,
	}

	p.Store.StoreProject(&project)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}

// handleGetAllLinkedProjects returns all linked projects list
func (p *Plugin) handleGetAllLinkedProjects(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	if len(projectList) == 0 {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.NoProjectLinked})
		return
	}

	response, err := json.Marshal(projectList)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleUnlinkProject unlinks a project
func (p *Plugin) handleUnlinkProject(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)

	var project *serializers.ProjectDetails
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		p.API.LogError(constants.ErrorDecodingBody, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	_, isProjectLinked := p.IsProjectLinked(projectList, *project)
	if !isProjectLinked {
		p.API.LogError(constants.ProjectNotFound, "Error")
		p.handleError(w, r, &serializers.Error{Code: http.StatusNotFound, Message: constants.ProjectNotFound})
		return
	}

	if err := p.Store.DeleteProject(project); err != nil {
		p.API.LogError(constants.ErrorUnlinkProject, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	successResponse := &serializers.SuccessResponse{
		Message: "success",
	}
	response, err := json.Marshal(&successResponse)
	if err != nil {
		p.API.LogError("Error marhsalling response", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleUnlinkProject unlinks a project
func (p *Plugin) handleGetUserAccountDetails(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)

	userDetails, err := p.Store.LoadUser(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorDecodingBody, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if userDetails.MattermostUserID == "" {
		p.API.LogError(constants.ConnectAccountFirst, "Error")
		p.handleError(w, r, &serializers.Error{Code: http.StatusUnauthorized, Message: constants.ConnectAccountFirst})
		return
	}

	response, err := json.Marshal(&userDetails)
	if err != nil {
		p.API.LogError("Error marhsalling response", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// API to create task of a project in an organization.
func (p *Plugin) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	var body *serializers.TaskCreateRequestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		p.API.LogError("Error in decoding the body for creating a task", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if err := body.IsValid(); err != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	task, statusCode, err := p.Client.CreateTask(body, mattermostUserID)
	if err != nil {
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err = w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	message := fmt.Sprintf(constants.CreatedTask, task.Link.HTML.Href)

	// Send message to DM.
	_, err = p.DM(mattermostUserID, message)
	_ = errors.Wrap(err, "failed to DM the created Task")
}

func (p *Plugin) handleCreateSubscriptions(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	body, err := serializers.CreateSubscriptionRequestPayloadFromJSON(r.Body)
	if err != nil {
		p.API.LogError("Error in decoding the body for creating subscriptions", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if err := body.IsSubscriptionRequestPayloadValid(); err != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	project, isProjectLinked := p.IsProjectLinked(projectList, serializers.ProjectDetails{OrganizationName: body.Organization, ProjectName: body.Project})
	if !isProjectLinked {
		p.API.LogError(constants.ProjectNotFound, "Error")
		p.handleError(w, r, &serializers.Error{Code: http.StatusNotFound, Message: constants.ProjectNotLinked})
		return
	}

	// TODO: remove later
	teamID := "qteks46as3befxj4ec1mip5ume"
	channel, channelErr := p.API.GetChannelByName(teamID, body.ChannelName, false)
	if channelErr != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: channelErr.DetailedError})
		return
	}

	subscriptionList, err := p.Store.GetAllSubscriptions(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchSubscriptionList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	_, isSubscriptionPresent := p.IsSubscriptionPresent(subscriptionList, serializers.SubscriptionDetails{OrganizationName: body.Organization, ProjectName: body.Project, ChannelID: channel.Id, EventType: body.EventType})
	if isSubscriptionPresent {
		p.API.LogError(constants.SubscriptionAlreadyPresent, "Error")
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.SubscriptionAlreadyPresent})
		return
	}

	pluginURL := p.GetPluginURL()
	subscription, statusCode, err := p.Client.CreateSubscription(body, project, channel.Id, pluginURL, mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorCreateSubscription, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	p.Store.StoreSubscription(&serializers.SubscriptionDetails{
		MattermostUserID: mattermostUserID,
		ProjectName:      body.Project,
		ProjectID:        subscription.PublisherInputs.ProjectID,
		OrganizationName: body.Organization,
		EventType:        body.EventType,
		ChannelID:        channel.Id,
		SubscriptionID:   subscription.ID,
	})
	response, err := json.Marshal(subscription)
	if err != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err = w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Plugin) handleGetSubscriptions(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	subscriptionList, err := p.Store.GetAllSubscriptions(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchSubscriptionList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	response, err := json.Marshal(subscriptionList)
	if err != nil {
		p.API.LogError(constants.ErrorFetchSubscriptionList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Plugin) handleNotificationSubscriptions(w http.ResponseWriter, r *http.Request) {
	body, err := serializers.SubscriptionNotificationFromJSON(r.Body)
	if err != nil {
		p.API.LogError("Error in decoding the body for creating notifications", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	channelID := r.URL.Query().Get("channelID")
	if channelID == "" {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.ChannelIDRequired})
		return
	}

	attachment := &model.SlackAttachment{
		Text: body.DetailedMessage.Markdown,
	}
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: channelID,
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	p.API.CreatePost(post)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleDeleteSubscriptions(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	body, err := serializers.DeleteSubscriptionRequestPayloadFromJSON(r.Body)
	if err != nil {
		p.API.LogError("Error in decoding the body for deleting subscriptions", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if err := body.IsSubscriptionRequestPayloadValid(); err != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	subscriptionList, err := p.Store.GetAllSubscriptions(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchSubscriptionList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	_, isSubscriptionPresent := p.IsSubscriptionPresent(subscriptionList, serializers.SubscriptionDetails{OrganizationName: body.Organization, ProjectName: body.Project, ChannelID: body.ChannelID, EventType: body.EventType})
	if !isSubscriptionPresent {
		p.API.LogError(constants.SubscriptionNotFound, "Error")
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.SubscriptionNotFound})
		return
	}
	
	statusCode, err := p.Client.DeleteSubscription(body.Organization, body.SubscriptionID, mattermostUserID)
	if err != nil {
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	p.Store.DeleteSubscription(&serializers.SubscriptionDetails{
		MattermostUserID: mattermostUserID,
		ProjectName:      body.Project,
		OrganizationName: body.Organization,
		EventType:        body.EventType,
		ChannelID:        body.ChannelID,
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// TODO: for testing purpose, remove later
func (p *Plugin) testAPI(w http.ResponseWriter, r *http.Request) {
	// TODO: remove later
	response, err := p.Client.TestApi()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// TODO: for testing purposes, remove later
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	u, _ := p.Store.LoadUser(mattermostUserID)
	t, _ := p.ParseAuthToken(u.AccessToken)
	fmt.Println("\n\n\n", t)
	res, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

func (p *Plugin) WithRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				p.API.LogError("Recovered from a panic",
					"url", r.URL.String(),
					"error", x,
					"stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Handles the static files under the assets directory.
func (p *Plugin) HandleStaticFiles() {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		p.API.LogWarn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	// This will serve static files from the 'assets' directory under '/static/<filename>'
	p.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(bundlePath, "assets")))))
}
