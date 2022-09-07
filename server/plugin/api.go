package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"

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
	// Plugin APIs
	s.HandleFunc(constants.PathCreateTasks, p.handleAuthRequired(p.checkOAuth(p.handleCreateTask))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathLinkProject, p.handleAuthRequired(p.checkOAuth(p.handleLink))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathGetAllLinkedProjects, p.handleAuthRequired(p.checkOAuth(p.handleGetAllLinkedProjects))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathUnlinkProject, p.handleAuthRequired(p.checkOAuth(p.handleUnlinkProject))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathUser, p.handleAuthRequired(p.checkOAuth(p.handleGetUserAccountDetails))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathSubscriptions, p.handleAuthRequired(p.checkOAuth(p.handleCreateSubscription))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathSubscriptions, p.handleAuthRequired(p.checkOAuth(p.handleGetSubscriptions))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathSubscriptionNotifications, p.handleSubscriptionNotifications).Methods(http.MethodPost)
	s.HandleFunc(constants.PathSubscriptions, p.handleAuthRequired(p.checkOAuth(p.handleDeleteSubscriptions))).Methods(http.MethodDelete)
	s.HandleFunc(constants.PathGetUserChannelsForTeam, p.handleAuthRequired(p.getUserChannelsForTeam)).Methods(http.MethodGet)
}

// API to create task of a project in an organization.
func (p *Plugin) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)

	body, err := serializers.CreateTaskRequestPayloadFromJSON(r.Body)
	if err != nil {
		p.API.LogError(constants.ErrorDecodingBody, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if validationErr := body.IsValid(); validationErr != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: validationErr.Error()})
		return
	}

	task, statusCode, err := p.Client.CreateTask(body, mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorCreateTask)
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	p.writeJSON(w, task)
	message := fmt.Sprintf(constants.CreatedTask, task.Fields.Type, task.ID, task.Link.HTML.Href, task.Fields.Title, task.Fields.CreatedBy.DisplayName)

	// Send message to DM.
	if _, DMErr := p.DM(mattermostUserID, message, true); DMErr != nil {
		p.API.LogError("Failed to DM", "Error", DMErr.Error())
	}
}

// API to link a project and an organization to a user.
func (p *Plugin) handleLink(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)

	body, err := serializers.LinkPayloadFromJSON(r.Body)
	if err != nil {
		p.API.LogError(constants.ErrorDecodingBody, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if linkValidationErr := body.IsLinkPayloadValid(); linkValidationErr != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: linkValidationErr.Error()})
		return
	}

	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	if _, isProjectLinked := p.IsProjectLinked(projectList, serializers.ProjectDetails{OrganizationName: body.Organization, ProjectName: body.Project}); isProjectLinked {
		if _, DMErr := p.DM(mattermostUserID, constants.AlreadyLinkedProject, true); DMErr != nil {
			p.API.LogError("Failed to DM", "Error", DMErr.Error())
		}
		return
	}

	response, statusCode, err := p.Client.Link(body, mattermostUserID)
	if err != nil {
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	project := serializers.ProjectDetails{
		MattermostUserID: mattermostUserID,
		ProjectID:        response.ID,
		ProjectName:      response.Name,
		OrganizationName: body.Organization,
	}

	if storeErr := p.Store.StoreProject(&project); storeErr != nil {
		p.API.LogError("Error in storing a project", "Error", storeErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: storeErr.Error()})
	}

	returnStatusOK(w)
}

// handleGetAllLinkedProjects returns all linked projects list
func (p *Plugin) handleGetAllLinkedProjects(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if len(projectList) == 0 {
		if _, err = w.Write([]byte("[]")); err != nil {
			p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		}
		return
	}

	p.writeJSON(w, projectList)
}

// handleUnlinkProject unlinks a project
func (p *Plugin) handleUnlinkProject(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)

	project, err := serializers.ProjectPayloadFromJSON(r.Body)
	if err != nil {
		p.API.LogError(constants.ErrorDecodingBody, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if validationErr := project.IsValid(); validationErr != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: validationErr.Error()})
		return
	}

	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	if _, isProjectLinked := p.IsProjectLinked(projectList, *project); !isProjectLinked {
		p.API.LogError(constants.ProjectNotFound, "Error")
		p.handleError(w, r, &serializers.Error{Code: http.StatusNotFound, Message: constants.ProjectNotFound})
		return
	}

	if deleteErr := p.Store.DeleteProject(&serializers.ProjectDetails{
		MattermostUserID: mattermostUserID,
		ProjectID:        project.ProjectID,
		ProjectName:      project.ProjectName,
		OrganizationName: project.OrganizationName,
	}); deleteErr != nil {
		p.API.LogError(constants.ErrorUnlinkProject, "Error", deleteErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: deleteErr.Error()})
	}

	successResponse := &serializers.SuccessResponse{
		Message: "success",
	}

	p.writeJSON(w, &successResponse)
}

func (p *Plugin) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	body, err := serializers.CreateSubscriptionRequestPayloadFromJSON(r.Body)
	if err != nil {
		p.API.LogError("Error in decoding the body for creating subscriptions", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if validationErr := body.IsSubscriptionRequestPayloadValid(); validationErr != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: validationErr.Error()})
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

	subscriptionList, err := p.Store.GetAllSubscriptions(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.FetchSubscriptionListError, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	if _, isSubscriptionPresent := p.IsSubscriptionPresent(subscriptionList, &serializers.SubscriptionDetails{OrganizationName: body.Organization, ProjectName: body.Project, ChannelID: body.ChannelID, EventType: body.EventType}); isSubscriptionPresent {
		p.API.LogError(constants.SubscriptionAlreadyPresent, "Error")
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.SubscriptionAlreadyPresent})
		return
	}

	subscription, statusCode, err := p.Client.CreateSubscription(body, project, body.ChannelID, p.GetPluginURL(), mattermostUserID)
	if err != nil {
		p.API.LogError(constants.CreateSubscriptionError, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	channel, channelErr := p.API.GetChannel(body.ChannelID)
	if channelErr != nil {
		p.API.LogError(constants.GetChannelError, "Error", channelErr.Error())
		http.Error(w, fmt.Sprintf("%s. Error: %s", constants.GetChannelError, channelErr.Error()), channelErr.StatusCode)
		return
	}

	if storeErr := p.Store.StoreSubscription(&serializers.SubscriptionDetails{
		MattermostUserID: mattermostUserID,
		ProjectName:      body.Project,
		ProjectID:        subscription.PublisherInputs.ProjectID,
		OrganizationName: body.Organization,
		EventType:        body.EventType,
		ChannelID:        body.ChannelID,
		SubscriptionID:   subscription.ID,
		ChannelName:      channel.DisplayName,
	}); storeErr != nil {
		p.API.LogError("Error in creating a subscription", "Error", storeErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: storeErr.Error()})
	}

	p.writeJSON(w, subscription)
}

func (p *Plugin) handleGetSubscriptions(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	subscriptionList, err := p.Store.GetAllSubscriptions(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.FetchSubscriptionListError, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	project := r.URL.Query().Get(constants.QueryParamProject)
	if project != "" {
		subscriptionByProject := []*serializers.SubscriptionDetails{}
		for _, subscription := range subscriptionList {
			if subscription.ProjectName == project {
				subscriptionByProject = append(subscriptionByProject, subscription)
			}
		}
		subscriptionList = subscriptionByProject
	}

	p.writeJSON(w, subscriptionList)
}

func (p *Plugin) handleSubscriptionNotifications(w http.ResponseWriter, r *http.Request) {
	body, err := serializers.SubscriptionNotificationFromJSON(r.Body)
	if err != nil {
		p.API.LogError("Error in decoding the body for creating notifications", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	channelID := r.URL.Query().Get("channelID")
	if channelID == "" {
		p.API.LogError(constants.ChannelIDRequired)
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.ChannelIDRequired})
		return
	}

	if !model.IsValidId(channelID) {
		p.API.LogError(constants.InvalidChannelID)
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.InvalidChannelID})
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
	if _, err := p.API.CreatePost(post); err != nil {
		p.API.LogError("Error in creating post", "Error", err.Error())
	}

	returnStatusOK(w)
}

func (p *Plugin) handleDeleteSubscriptions(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	body, err := serializers.DeleteSubscriptionRequestPayloadFromJSON(r.Body)
	if err != nil {
		p.API.LogError("Error in decoding the body for deleting subscriptions", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if validationErr := body.IsSubscriptionRequestPayloadValid(); validationErr != nil {
		p.API.LogDebug("Request payload is not valid", "Error", validationErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: validationErr.Error()})
		return
	}

	subscriptionList, err := p.Store.GetAllSubscriptions(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.FetchSubscriptionListError, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	subscription, isSubscriptionPresent := p.IsSubscriptionPresent(subscriptionList, &serializers.SubscriptionDetails{
		OrganizationName: body.Organization,
		ProjectName:      body.Project,
		ChannelID:        body.ChannelID,
		EventType:        body.EventType,
	})
	if !isSubscriptionPresent {
		p.API.LogError(constants.SubscriptionNotFound)
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.SubscriptionNotFound})
		return
	}

	if statusCode, err := p.Client.DeleteSubscription(body.Organization, subscription.SubscriptionID, mattermostUserID); err != nil {
		p.API.LogError(constants.DeleteSubscriptionError, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	if deleteErr := p.Store.DeleteSubscription(&serializers.SubscriptionDetails{
		MattermostUserID: mattermostUserID,
		ProjectName:      body.Project,
		OrganizationName: body.Organization,
		EventType:        body.EventType,
		ChannelID:        body.ChannelID,
	}); deleteErr != nil {
		p.API.LogError(constants.DeleteSubscriptionError, "Error", deleteErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: deleteErr.Error()})
	}

	returnStatusOK(w)
}

func (p *Plugin) getUserChannelsForTeam(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	pathParams := mux.Vars(r)
	teamID := pathParams[constants.PathParamTeamID]
	if !model.IsValidId(teamID) {
		p.API.LogError("Invalid team id")
		http.Error(w, "Invalid team id", http.StatusBadRequest)
		return
	}

	channels, channelErr := p.API.GetChannelsForTeamForUser(teamID, mattermostUserID, false)
	if channelErr != nil {
		p.API.LogError(constants.GetChannelError, "Error", channelErr.Error())
		http.Error(w, fmt.Sprintf("%s. Error: %s", constants.GetChannelError, channelErr.Error()), channelErr.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if channels == nil {
		_, _ = w.Write([]byte("[]"))
		return
	}

	var requiredChannels []*model.Channel
	for _, channel := range channels {
		if channel.Type == model.CHANNEL_PRIVATE || channel.Type == model.CHANNEL_OPEN {
			requiredChannels = append(requiredChannels, channel)
		}
	}
	if requiredChannels == nil {
		_, _ = w.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(w).Encode(requiredChannels); err != nil {
		p.API.LogError("Error while writing response", "Error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (p *Plugin) checkOAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
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

func returnStatusOK(w http.ResponseWriter) {
	m := make(map[string]string)
	w.Header().Set("Content-Type", "application/json")
	m[model.STATUS] = model.STATUS_OK
	_, _ = w.Write([]byte(model.MapToJson(m)))
}

// handleAuthRequired verifies if the provided request is performed by an authorized source.
func (p *Plugin) handleAuthRequired(handleFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
		if mattermostUserID == "" {
			p.handleError(w, r, &serializers.Error{Code: http.StatusUnauthorized, Message: constants.NotAuthorized})
			return
		}

		handleFunc(w, r)
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

// handleGetUserAccountDetails provides user details
func (p *Plugin) handleGetUserAccountDetails(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	userDetails, err := p.Store.LoadUser(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorLoadingDataFromKVStore, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	if userDetails.MattermostUserID == "" {
		p.API.LogError(constants.ConnectAccountFirst)
		p.handleError(w, r, &serializers.Error{Code: http.StatusUnauthorized, Message: constants.ConnectAccountFirst})
		return
	}

	p.API.PublishWebSocketEvent(
		constants.WSEventConnect,
		nil,
		&model.WebsocketBroadcast{UserId: mattermostUserID},
	)

	p.writeJSON(w, &userDetails)
}

func (p *Plugin) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		p.API.LogError("Failed to marshal JSON response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		p.API.LogError("Failed to write JSON response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
