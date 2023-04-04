package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
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
	s.HandleFunc(constants.PathOAuthConnect, p.handleAuthRequired(p.OAuthConnect)).Methods(http.MethodGet)
	s.HandleFunc(constants.PathOAuthCallback, p.handleAuthRequired(p.OAuthComplete)).Methods(http.MethodGet)
	// Plugin APIs
	s.HandleFunc(constants.PathCreateTasks, p.handleAuthRequired(p.checkOAuth(p.handleCreateTask))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathLinkProject, p.handleAuthRequired(p.checkOAuth(p.handleLink))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathGetAllLinkedProjects, p.handleAuthRequired(p.checkOAuth(p.handleGetAllLinkedProjects))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathUnlinkProject, p.handleAuthRequired(p.checkOAuth(p.handleUnlinkProject))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathUser, p.handleAuthRequired(p.checkOAuth(p.handleGetUserAccountDetails))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathSubscriptions, p.handleAuthRequired(p.checkOAuth(p.handleCreateSubscription))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathGetSubscriptions, p.handleAuthRequired(p.checkOAuth(p.handleGetSubscriptions))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathSubscriptionNotifications, p.handleSubscriptionNotifications).Methods(http.MethodPost)
	s.HandleFunc(constants.PathSubscriptions, p.handleAuthRequired(p.checkOAuth(p.handleDeleteSubscriptions))).Methods(http.MethodDelete)
	s.HandleFunc(constants.PathPipelineReleaseRequest, p.handleAuthRequired(p.checkOAuth(p.handlePipelineApproveOrRejectReleaseRequest))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathPipelineRunRequest, p.handleAuthRequired(p.checkOAuth(p.handlePipelineApproveOrRejectRunRequest))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathPipelineCommentModal, p.handleAuthRequired(p.checkOAuth(p.handlePipelineCommentModal))).Methods(http.MethodPost)
	s.HandleFunc(constants.PathGetSubscriptionFilterPossibleValues, p.handleAuthRequired(p.checkOAuth(p.handleGetSubscriptionFilterPossibleValues))).Methods(http.MethodPost)
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
	message := fmt.Sprintf(constants.CreatedTask, task.ID, task.Fields.Title, task.Link.HTML.Href, task.Fields.Type, task.Fields.CreatedBy.DisplayName)

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

	if _, isProjectLinked := p.IsProjectLinked(projectList, serializers.ProjectDetails{OrganizationName: strings.ToLower(body.Organization), ProjectName: cases.Title(language.Und).String(body.Project)}); isProjectLinked {
		returnStatusWithMessage(w, http.StatusOK, constants.AlreadyLinkedProject)
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
		ProjectName:      cases.Title(language.Und).String(body.Project),
		OrganizationName: strings.ToLower(body.Organization),
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

	if project.DeleteSubscriptions {
		if statusCode, err := p.handleDeleteAllSubscriptions(mattermostUserID, project.ProjectID); err != nil {
			p.API.LogError("Error deleting the project subscriptions", "Error", err.Error())
			p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
			return
		}
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

func (p *Plugin) handleDeleteAllSubscriptions(mattermostUserID, projectID string) (int, error) {
	subscriptionList, err := p.Store.GetAllSubscriptions(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.FetchSubscriptionListError, "Error", err.Error())
		return http.StatusInternalServerError, err
	}

	for _, subscription := range subscriptionList {
		if subscription.ProjectID == projectID {
			statusCode, deleteErr := p.deleteSubscription(subscription, mattermostUserID)
			if deleteErr != nil {
				p.API.LogError(constants.DeleteSubscriptionError, "Error", deleteErr.Error())
				return statusCode, deleteErr
			}
		}
	}

	return http.StatusOK, nil
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

	if statusCode, channelAccessErr := p.CheckValidChannelForSubscription(body.ChannelID, mattermostUserID); channelAccessErr != nil {
		p.API.LogError(constants.ErrorCreateSubscription, "Error", channelAccessErr.Error())

		message := channelAccessErr.Error()
		responseStatusCode := statusCode
		if statusCode == http.StatusNotFound {
			message = "you are not allowed to create subscription for the provided channel"
			responseStatusCode = http.StatusForbidden
		}

		p.handleError(w, r, &serializers.Error{Code: responseStatusCode, Message: message})
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

	if _, isSubscriptionPresent := p.IsSubscriptionPresent(subscriptionList, &serializers.SubscriptionDetails{
		OrganizationName: body.Organization,
		ProjectName:      body.Project,
		ChannelID:        body.ChannelID,
		EventType:        body.EventType,
		// Below all are filters that could be present on different categories of subscriptions from Boards, Repos and Pipelines
		Repository:                   body.Repository,
		TargetBranch:                 body.TargetBranch,
		PullRequestCreatedBy:         body.PullRequestCreatedBy,
		PullRequestReviewersContains: body.PullRequestReviewersContains,
		PushedBy:                     body.PushedBy,
		MergeResult:                  body.MergeResult,
		NotificationType:             body.NotificationType,
		AreaPath:                     body.AreaPath,
		BuildStatus:                  body.BuildStatus,
		BuildPipeline:                body.BuildPipeline,
		StageName:                    body.StageName,
		ReleasePipeline:              body.ReleasePipeline,
		ReleaseStatus:                body.ReleaseStatus,
		ApprovalType:                 body.ApprovalType,
		ApprovalStatus:               body.ApprovalStatus,
		RunPipeline:                  body.RunPipeline,
		RunStageName:                 body.RunStageName,
		RunEnvironmentName:           body.RunEnvironmentName,
		RunStageNameID:               body.RunStageNameID,
		RunStageStateID:              body.RunStageStateID,
		RunStageResultID:             body.RunStageResultID,
		RunStateID:                   body.RunStateID,
		RunResultID:                  body.RunResultID,
	}); isSubscriptionPresent {
		p.API.LogError(constants.SubscriptionAlreadyPresent, "Error")
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.SubscriptionAlreadyPresent})
		return
	}

	uniqueWebhookSecret := uuid.New().String()
	subscription, statusCode, err := p.Client.CreateSubscription(body, project, body.ChannelID, p.GetPluginURL(), mattermostUserID, uniqueWebhookSecret)
	if err != nil {
		p.API.LogError(constants.CreateSubscriptionError, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	if err := p.Store.StoreSubscriptionAndChannelIDMap(subscription.ID, uniqueWebhookSecret, body.ChannelID); err != nil {
		p.API.LogError("Error storing channel ID for subscription", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	channel, channelErr := p.API.GetChannel(body.ChannelID)
	if channelErr != nil {
		p.API.LogError(constants.GetChannelError, "Error", channelErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: constants.GetChannelError})
		return
	}

	user, userErr := p.API.GetUser(mattermostUserID)
	if userErr != nil {
		p.API.LogError(constants.GetUserError, "Error", userErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: constants.GetUserError})
		return
	}

	createdByDisplayName := user.Username

	showFullName := p.API.GetConfig().PrivacySettings.ShowFullName
	// If "PrivacySettings.ShowFullName" is true then show the user's first/last name
	// If the user's first/last name doesn't exist then show the username as fallback
	if showFullName != nil && *showFullName && (user.FirstName != "" || user.LastName != "") {
		createdByDisplayName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	if storeErr := p.Store.StoreSubscription(&serializers.SubscriptionDetails{
		MattermostUserID: mattermostUserID,
		ProjectName:      body.Project,
		ProjectID:        project.ProjectID,
		OrganizationName: body.Organization,
		EventType:        body.EventType,
		ServiceType:      body.ServiceType,
		ChannelID:        body.ChannelID,
		SubscriptionID:   subscription.ID,
		ChannelName:      channel.DisplayName,
		ChannelType:      channel.Type,
		CreatedBy:        strings.TrimSpace(createdByDisplayName),
		// Below all are filters that could be present on different categories of subscriptions from Boards, Repos and Pipelines
		Repository:                       body.Repository,
		TargetBranch:                     body.TargetBranch,
		RepositoryName:                   body.RepositoryName,
		PullRequestCreatedBy:             body.PullRequestCreatedBy,
		PullRequestReviewersContains:     body.PullRequestReviewersContains,
		PullRequestCreatedByName:         body.PullRequestCreatedByName,
		PullRequestReviewersContainsName: body.PullRequestReviewersContainsName,
		PushedBy:                         body.PushedBy,
		PushedByName:                     body.PushedByName,
		MergeResult:                      body.MergeResult,
		MergeResultName:                  body.MergeResultName,
		NotificationType:                 body.NotificationType,
		NotificationTypeName:             body.NotificationTypeName,
		AreaPath:                         body.AreaPath,
		BuildStatus:                      body.BuildStatus,
		BuildPipeline:                    body.BuildPipeline,
		StageName:                        body.StageName,
		ReleasePipeline:                  body.ReleasePipeline,
		ReleaseStatus:                    body.ReleaseStatus,
		ApprovalType:                     body.ApprovalType,
		ApprovalStatus:                   body.ApprovalStatus,
		BuildStatusName:                  body.BuildStatusName,
		StageNameValue:                   body.StageNameValue,
		ReleasePipelineName:              body.ReleasePipelineName,
		ReleaseStatusName:                body.ReleaseStatusName,
		ApprovalTypeName:                 body.ApprovalTypeName,
		ApprovalStatusName:               body.ApprovalStatusName,
		RunPipeline:                      body.RunPipeline,
		RunPipelineName:                  body.RunPipelineName,
		RunStageName:                     body.RunStageName,
		RunEnvironmentName:               body.RunEnvironmentName,
		RunStageNameID:                   body.RunStageNameID,
		RunStageStateID:                  body.RunStageStateID,
		RunStageStateIDName:              body.RunStageStateIDName,
		RunStageResultID:                 body.RunStageResultID,
		RunStateID:                       body.RunStateID,
		RunStateIDName:                   body.RunStateIDName,
		RunResultID:                      body.RunResultID,
	}); storeErr != nil {
		p.API.LogError("Error in creating a subscription", "Error", storeErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: storeErr.Error()})
	}

	p.writeJSON(w, subscription)
}

func (p *Plugin) handleGetSubscriptions(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)

	pathParams := mux.Vars(r)
	teamID := pathParams[constants.PathParamTeamID]
	if !model.IsValidId(teamID) {
		p.API.LogWarn("Invalid team id")
		http.Error(w, "Invalid team id", http.StatusBadRequest)
		return
	}

	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogWarn(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	organization := pathParams[constants.PathParamOrganization]
	project := pathParams[constants.PathParamProject]
	organizationName := strings.ToLower(organization)
	projectName := cases.Title(language.Und).String(project)
	if _, isProjectLinked := p.IsProjectLinked(projectList, serializers.ProjectDetails{
		OrganizationName: organizationName,
		ProjectName:      projectName,
	}); !isProjectLinked {
		p.API.LogWarn(fmt.Sprintf("Project %s is not linked", project))
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: "requested project is not linked"})
		return
	}

	var subscriptionList []*serializers.SubscriptionDetails
	var subscriptionErr error
	createdBy := r.URL.Query().Get(constants.QueryParamCreatedBy)
	switch createdBy {
	case constants.FilterCreatedByMe, "":
		subscriptionList, subscriptionErr = p.Store.GetAllSubscriptions(mattermostUserID)
	case constants.FilterCreatedByAnyone:
		subscriptionList, subscriptionErr = p.Store.GetAllSubscriptions("")
	}
	if subscriptionErr != nil {
		p.API.LogWarn(constants.FetchSubscriptionListError, "Error", subscriptionErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: subscriptionErr.Error()})
		return
	}

	offset, limit := p.GetOffsetAndLimitFromQueryParams(r)
	channelID := r.URL.Query().Get(constants.QueryParamChannelID)
	serviceType := r.URL.Query().Get(constants.QueryParamServiceType)
	eventType := r.URL.Query().Get(constants.QueryParamEventType)

	subscriptionByProject := []*serializers.SubscriptionDetails{}
	for _, subscription := range subscriptionList {
		if subscription.ProjectName == project {
			if channelID == "" || subscription.ChannelID == channelID {
				switch serviceType {
				case "", constants.FilterAll:
					subscriptionByProject = append(subscriptionByProject, subscription)
				case constants.FilterBoards:
					switch eventType {
					case "", constants.FilterAll:
						if constants.ValidSubscriptionEventsForBoards[subscription.EventType] {
							subscriptionByProject = append(subscriptionByProject, subscription)
						}
					default:
						if subscription.EventType == eventType {
							subscriptionByProject = append(subscriptionByProject, subscription)
						}
					}
				case constants.FilterRepos:
					switch eventType {
					case "", constants.FilterAll:
						if constants.ValidSubscriptionEventsForRepos[subscription.EventType] {
							subscriptionByProject = append(subscriptionByProject, subscription)
						}
					default:
						if subscription.EventType == eventType {
							subscriptionByProject = append(subscriptionByProject, subscription)
						}
					}
				case constants.FilterPipelines:
					switch eventType {
					case "", constants.FilterAll:
						if constants.ValidSubscriptionEventsForPipelines[subscription.EventType] {
							subscriptionByProject = append(subscriptionByProject, subscription)
						}
					default:
						if subscription.EventType == eventType {
							subscriptionByProject = append(subscriptionByProject, subscription)
						}
					}
				}
			}
		}
	}

	sort.Slice(subscriptionByProject, func(i, j int) bool {
		return subscriptionByProject[i].CreatedAt.After(subscriptionByProject[j].CreatedAt)
	})

	filteredSubscriptionList, filteredSubscriptionErr := p.GetSubscriptionsForAccessibleChannelsOrProjects(subscriptionByProject, teamID, mattermostUserID, constants.FilterCreatedByAnyone)
	if filteredSubscriptionErr != nil {
		p.API.LogWarn(constants.FetchFilteredSubscriptionListError, "Error", filteredSubscriptionErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: filteredSubscriptionErr.Error()})
		return
	}

	paginatedSubscriptions := []*serializers.SubscriptionDetails{}
	for index, subscription := range filteredSubscriptionList {
		if len(paginatedSubscriptions) == limit {
			break
		}
		if index >= offset {
			paginatedSubscriptions = append(paginatedSubscriptions, subscription)
		}
	}

	p.writeJSON(w, paginatedSubscriptions)
}

func (p *Plugin) getReviewersListString(reviewersList []serializers.Reviewer) string {
	reviewers := ""
	for i := 0; i < len(reviewersList); i++ {
		if i != len(reviewersList)-1 {
			reviewers += fmt.Sprintf("%s, ", reviewersList[i].DisplayName)
		} else {
			reviewers += reviewersList[i].DisplayName
		}
	}

	if reviewers == "" {
		return "None" // When no reviewers are added
	}
	return reviewers
}

func (p *Plugin) getPipelineReleaseEnvironmentList(environments []*serializers.Environment) string {
	envs := ""
	for index, env := range environments {
		envs += env.Name
		if index != (len(environments) - 1) {
			envs += " | "
		}
	}

	if envs == "" {
		return "None"
	}

	return envs
}

func (p *Plugin) handleSubscriptionNotifications(w http.ResponseWriter, r *http.Request) {
	body, err := serializers.SubscriptionNotificationFromJSON(r.Body)
	if err != nil {
		p.API.LogError("Error in decoding the body for listening notifications", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	webhookSecret := r.URL.Query().Get(constants.AzureDevopsQueryParamWebhookSecret)
	if webhookSecret == "" {
		p.API.LogError(constants.ErrorUnauthorisedSubscriptionsWebhookRequest)
		p.handleError(w, r, &serializers.Error{Code: http.StatusUnauthorized, Message: constants.ErrorUnauthorisedSubscriptionsWebhookRequest})
		return
	}

	channelID, status, err := p.VerifySubscriptionWebhookSecretAndGetChannelID(body.SubscriptionID, webhookSecret)
	if err != nil {
		p.API.LogError("Unable to verify webhook secret for subscription", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: status, Message: err.Error()})
		return
	}

	var attachment *model.SlackAttachment
	switch body.EventType {
	case constants.SubscriptionEventWorkItemCreated, constants.SubscriptionEventWorkItemDeleted:
		attachment = &model.SlackAttachment{
			AuthorName: constants.SlackAttachmentAuthorNameBoards,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameBoardsIcon),
			Color:      constants.IconColorBoards,
			Pretext:    body.Message.Markdown,
			Title:      body.Resource.Fields.Title.(string),
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Area Path",
					Value: body.Resource.Fields.AreaPath,
					Short: true,
				},
				{
					Title: "State",
					Value: body.Resource.Fields.State,
					Short: true,
				},
				{
					Title: "Workitem Type",
					Value: body.Resource.Fields.WorkItemType,
				},
			},
			Footer:     body.Resource.Fields.ProjectName.(string),
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventWorkItemCommented:
		reg := regexp.MustCompile(constants.WorkItemCommentedOnMarkdownRegex)
		comment := reg.Split(body.DetailedMessage.Markdown, -1)

		attachment = &model.SlackAttachment{
			AuthorName: constants.SlackAttachmentAuthorNameBoards,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameBoardsIcon),
			Color:      constants.IconColorBoards,
			Pretext:    body.Message.Markdown,
			Title:      "Comment",
			Text:       strings.TrimSpace(comment[len(comment)-1]),
			Footer:     body.Resource.Fields.ProjectName.(string),
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventWorkItemUpdated:
		attachment = &model.SlackAttachment{
			AuthorName: constants.SlackAttachmentAuthorNameBoards,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameBoardsIcon),
			Color:      constants.IconColorBoards,
			Pretext:    body.Message.Markdown,
			Title:      body.Resource.Revision.Fields.Title.(string),
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Area Path",
					Value: body.Resource.Revision.Fields.AreaPath,
					Short: true,
				},
				{
					Title: "State",
					Value: body.Resource.Revision.Fields.State,
					Short: true,
				},
				{
					Title: "Workitem Type",
					Value: body.Resource.Revision.Fields.WorkItemType,
				},
			},
			Footer:     body.Resource.Revision.Fields.ProjectName.(string),
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventPullRequestCreated, constants.SubscriptionEventPullRequestUpdated, constants.SubscriptionEventPullRequestMerged:
		reviewers := p.getReviewersListString(body.Resource.Reviewers)

		var targetBranchName, sourceBranchName string
		if len(strings.Split(body.Resource.TargetRefName, "/")) == 3 {
			targetBranchName = strings.Split(body.Resource.TargetRefName, "/")[2]
		}

		if len(strings.Split(body.Resource.SourceRefName, "/")) == 3 {
			sourceBranchName = strings.Split(body.Resource.SourceRefName, "/")[2]
		}

		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNameRepos,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameReposIcon),
			Color:      constants.IconColorRepos,
			Title:      fmt.Sprintf("%d: %s", body.Resource.PullRequestID, body.Resource.Title),
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Target Branch",
					Value: targetBranchName,
					Short: true,
				},
				{
					Title: "Source Branch",
					Value: sourceBranchName,
					Short: true,
				},
				{
					Title: "Reviewer(s)",
					Value: reviewers,
				},
			},
			Footer:     body.Resource.Repository.Name,
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventPullRequestCommented:
		reviewers := p.getReviewersListString(body.Resource.PullRequest.Reviewers)

		var targetBranchName, sourceBranchName string
		if len(strings.Split(body.Resource.PullRequest.TargetRefName, "/")) == 3 {
			targetBranchName = strings.Split(body.Resource.PullRequest.TargetRefName, "/")[2]
		}

		if len(strings.Split(body.Resource.PullRequest.SourceRefName, "/")) == 3 {
			sourceBranchName = strings.Split(body.Resource.PullRequest.SourceRefName, "/")[2]
		}

		// Convert map to json string
		jsonBytes, err := json.Marshal(body.Resource.Comment)
		if err != nil {
			p.API.LogError(err.Error())
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		// Convert json string to struct
		var comment *serializers.Comment
		if err := json.Unmarshal(jsonBytes, &comment); err != nil {
			p.API.LogError(err.Error())
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNameRepos,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameReposIcon),
			Color:      constants.IconColorRepos,
			Title:      fmt.Sprintf("%d: %s", body.Resource.PullRequest.PullRequestID, body.Resource.PullRequest.Title),
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Target Branch",
					Value: targetBranchName,
					Short: true,
				},
				{
					Title: "Source Branch",
					Value: sourceBranchName,
					Short: true,
				},
				{
					Title: "Reviewer(s)",
					Value: reviewers,
				},
				{
					Title: "Comment",
					Value: comment.Content,
				},
			},
			Footer:     body.Resource.PullRequest.Repository.Name,
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventCodePushed:
		commits := ""
		for i := 0; i < len(body.Resource.Commits); i++ {
			commits += fmt.Sprintf("\n[%s](%s): **%s**", body.Resource.Commits[i].CommitID[0:8], body.Resource.Commits[i].URL, body.Resource.Commits[i].Comment)
		}

		if commits == "" {
			commits = "None" // When no commits are present
		}

		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNameRepos,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameReposIcon),
			Color:      constants.IconColorRepos,
			Title:      "Commit(s)",
			Text:       commits,
			Footer:     fmt.Sprintf("%s | %s", strings.Split(body.Resource.RefUpdates[0].Name, "/")[2], body.Resource.Repository.Name),
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameGitBranchIcon),
		}
	case constants.SubscriptionEventBuildCompleted:
		startTime, err := time.Parse(constants.DateTimeLayout, strings.Split(body.Resource.StartTime, ".")[0])
		if err != nil {
			p.API.LogError(err.Error())
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		finishTime, err := time.Parse(constants.DateTimeLayout, strings.Split(body.Resource.FinishTime, ".")[0])
		if err != nil {
			p.API.LogError(err.Error())
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Build pipeline",
					Value: body.Resource.Definition.Name,
					Short: true,
				},
				{
					Title: "Branch",
					Value: body.Resource.SourceBranch,
					Short: true,
				},
				{
					Title: "Requested for",
					Value: body.Resource.RequestedFor.Name,
					Short: true,
				},
				{
					Title: "Duration",
					Value: time.Time{}.Add(finishTime.Sub(startTime)).Format(constants.TimeLayout),
					Short: true,
				},
			},
			Footer:     body.Resource.Project.Name,
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventReleaseCreated:
		artifacts := ""
		for i := 0; i < len(body.Resource.Release.Artifacts); i++ {
			artifacts += body.Resource.Release.Artifacts[i].Name
			if i != len(body.Resource.Release.Artifacts)-1 {
				artifacts += ", "
			}
		}

		if artifacts == "" {
			artifacts = "No artifacts"
		}

		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Release pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Release.ReleaseDefinition.Name, body.Resource.Release.ReleaseDefinition.Links.Web.Href),
					Short: true,
				},
				{
					Title: "Created by",
					Value: body.Resource.Release.CreatedBy.DisplayName,
					Short: true,
				},
				{
					Title: "Trigger reason",
					Value: cases.Title(language.Und).String(body.Resource.Release.Reason),
					Short: true,
				},
				{
					Title: "Artifacts",
					Value: artifacts,
					Short: true,
				},
			},
			Footer:     body.Resource.Project.Name,
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventReleaseAbandoned:
		abandonTime, err := time.Parse(constants.DateTimeLayout, strings.Split(body.Resource.Release.ModifiedOn, ".")[0])
		if err != nil {
			p.API.LogError(err.Error())
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Release pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Release.ReleaseDefinition.Name, body.Resource.Release.ReleaseDefinition.Links.Web.Href),
					Short: true,
				},
				{
					Title: "Abandoned by",
					Value: body.Resource.Release.ModifiedBy.DisplayName,
					Short: true,
				},
				{
					Title: "Abandoned on",
					Value: abandonTime.Format(constants.DateTimeFormat),
				},
			},
			Footer:     body.Resource.Project.Name,
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventReleaseDeploymentStarted:
		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Release pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Release.ReleaseDefinition.Name, body.Resource.Release.ReleaseDefinition.Links.Web.Href),
					Short: true,
				},
				{
					Title: "Release",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Release.Name, body.Resource.Release.Links.Web.Href),
					Short: true,
				},
			},
			Footer:     body.Resource.Project.Name,
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventReleaseDeploymentCompleted:
		comment := body.Resource.Comment.(string)
		if comment == "" {
			comment = "No comments"
		}

		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Release pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Environment.ReleaseDefinition.Name, body.Resource.Environment.ReleaseDefinition.Links.Web.Href),
					Short: true,
				},
				{
					Title: "Release",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Environment.Release.Name, body.Resource.Environment.Release.Links.Web.Href),
					Short: true,
				},
				{
					Title: "Comment",
					Value: comment,
				},
			},
			Footer:     body.Resource.Project.Name,
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventRunStageStateChanged:
		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Pipeline.Name, body.Resource.Stage.Links.PipelineWeb.Href),
					Short: true,
				},
			},
		}
	case constants.SubscriptionEventRunStageWaitingForApproval:
		organization := ""
		webLinkPaths := strings.Split(body.Resource.Pipeline.Links.Web.Href, "/")
		if len(webLinkPaths) >= 4 {
			organization = webLinkPaths[3]
		}

		approverTitle := "Approver(s)"
		if body.Resource.Approval.ExecutionOrder == "inSequence" {
			approverTitle = "Approver(s) in sequence"
		} else if body.Resource.Approval.MinRequiredApprovers > 0 && len(body.Resource.Approval.Steps) > body.Resource.Approval.MinRequiredApprovers {
			approverTitle = fmt.Sprintf("Approvers (any %d)", body.Resource.Approval.MinRequiredApprovers)
		}

		approvers := ""
		for _, approvalStep := range body.Resource.Approval.Steps {
			approvers += approvalStep.AssignedApprover.DisplayName + "\n"
		}

		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Run pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Pipeline.Name, body.Resource.Pipeline.Links.Web.Href),
					Short: true,
				},
				{
					Title: "Stage",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Stage.Name, body.Resource.Stage.Links.Web.Href),
					Short: true,
				},
				{
					Title: approverTitle,
					Value: approvers,
				},
			},
			Actions: []*model.PostAction{
				{
					Id:    constants.PipelineRequestIDApproved,
					Type:  model.POST_ACTION_TYPE_BUTTON,
					Name:  "Approve",
					Style: "primary",
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("%s%s", p.GetPluginURL(), constants.PathPipelineCommentModal),
						Context: map[string]interface{}{
							constants.PipelineRequestContextRequestName:  constants.PipelineRequestNameRun,
							constants.PipelineRequestContextApprovalID:   body.Resource.Approval.ID,
							constants.PipelineRequestContextOrganization: organization,
							constants.PipelineRequestContextRequestType:  constants.PipelineRequestIDApproved,
							constants.PipelineRequestContextProjectID:    body.Resource.ProjectID,
						},
					},
				},
				{
					Id:    constants.PipelineRequestIDRejected,
					Type:  model.POST_ACTION_TYPE_BUTTON,
					Name:  "Reject",
					Style: "danger",
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("%s%s", p.GetPluginURL(), constants.PathPipelineCommentModal),
						Context: map[string]interface{}{
							constants.PipelineRequestContextRequestName:  constants.PipelineRequestNameRun,
							constants.PipelineRequestContextApprovalID:   body.Resource.Approval.ID,
							constants.PipelineRequestContextOrganization: organization,
							constants.PipelineRequestContextRequestType:  constants.PipelineRequestIDRejected,
							constants.PipelineRequestContextProjectID:    body.Resource.ProjectID,
						},
					},
				},
			},
		}
	case constants.SubscriptionEventReleaseDeploymentEventPending:
		artifacts := ""
		for i, artifact := range body.Resource.Release.Artifacts {
			artifacts += artifact.Name
			if i != len(body.Resource.Release.Artifacts)-1 {
				artifacts += ", "
			}
		}

		if artifacts == "" {
			artifacts = "No artifacts"
		}

		organization := ""
		webLinkPaths := strings.Split(body.Resource.Release.ReleaseDefinition.Links.Web.Href, "/")
		if len(webLinkPaths) >= 4 {
			organization = webLinkPaths[3]
		}
		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Release pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Release.Name, body.Resource.Release.ReleaseDefinition.Links.Web.Href),
					Short: true,
				},
				{
					Title: "Artifacts",
					Value: artifacts,
					Short: true,
				},
				{
					Title: "Approver(s)",
					Value: body.Resource.Approval.Approver.DisplayName,
				},
			},
			Actions: []*model.PostAction{
				{
					Id:    constants.PipelineRequestIDApproved,
					Type:  model.POST_ACTION_TYPE_BUTTON,
					Name:  "Approve",
					Style: "primary",
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("%s%s", p.GetPluginURL(), constants.PathPipelineCommentModal),
						Context: map[string]interface{}{
							constants.PipelineRequestContextRequestName:  constants.PipelineRequestNameRelease,
							constants.PipelineRequestContextApprovalID:   body.Resource.Approval.ID,
							constants.PipelineRequestContextOrganization: organization,
							constants.PipelineRequestContextProjectName:  body.Resource.Project.Name,
							constants.PipelineRequestContextRequestType:  constants.PipelineRequestIDApproved,
						},
					},
				},
				{
					Id:    constants.PipelineRequestIDRejected,
					Type:  model.POST_ACTION_TYPE_BUTTON,
					Name:  "Reject",
					Style: "danger",
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("%s%s", p.GetPluginURL(), constants.PathPipelineCommentModal),
						Context: map[string]interface{}{
							constants.PipelineRequestContextRequestName:  constants.PipelineRequestNameRelease,
							constants.PipelineRequestContextApprovalID:   body.Resource.Approval.ID,
							constants.PipelineRequestContextOrganization: organization,
							constants.PipelineRequestContextProjectName:  body.Resource.Project.Name,
							constants.PipelineRequestContextRequestType:  constants.PipelineRequestIDRejected,
						},
					},
				},
			},
		}
	case constants.SubscriptionEventReleaseDeploymentApprovalCompleted:
		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Release pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Release.Name, body.Resource.Release.Links.Web.Href),
					Short: true,
				},
			},
			Footer:     body.Resource.Project.Name,
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventRunStateChanged:
		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Pipeline.Name, body.Resource.Run.Links.PipelineWeb.Href),
					Short: true,
				},
			},
		}
	case constants.SubscriptionEventRunStageApprovalCompleted:
		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Pipeline.Name, body.Resource.Pipeline.Links.Web.Href),
					Short: true,
				},
			},
			Footer:     body.Resource.Project.Name,
			FooterIcon: fmt.Sprintf(constants.PublicFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
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

func (p *Plugin) handlePipelineCommentModal(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	decoder := json.NewDecoder(r.Body)
	postActionIntegrationRequest := &model.PostActionIntegrationRequest{}
	if err := decoder.Decode(&postActionIntegrationRequest); err != nil {
		// TODO: prevent posting any error messages except oAuth in DM for now and use dialog for all such cases
		p.handlePipelineApprovalRequestUpdateError("Error decoding PostActionIntegrationRequest param: ", mattermostUserID, err)
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	var project, approval, url string
	if postActionIntegrationRequest.Context[constants.PipelineRequestContextRequestName].(string) == constants.PipelineRequestNameRelease {
		url = fmt.Sprintf("%s%s", p.GetPluginURL(), constants.PathPipelineReleaseRequest)
		approval = fmt.Sprintf("%f", postActionIntegrationRequest.Context[constants.PipelineRequestContextApprovalID].(float64))
		project = postActionIntegrationRequest.Context[constants.PipelineRequestContextProjectName].(string)
	} else {
		url = fmt.Sprintf("%s%s", p.GetPluginURL(), constants.PathPipelineRunRequest)
		approval = postActionIntegrationRequest.Context[constants.PipelineRequestContextApprovalID].(string)
		project = postActionIntegrationRequest.Context[constants.PipelineRequestContextProjectID].(string)
	}

	organization := postActionIntegrationRequest.Context[constants.PipelineRequestContextOrganization].(string)
	requestType := postActionIntegrationRequest.Context[constants.PipelineRequestContextRequestType].(string)

	elements := []model.DialogElement{
		{
			DisplayName: "Comment:",
			Name:        constants.DialogFieldNameComment,
			Type:        "text",
			Optional:    true,
		},
	}

	dialogTitle := "Confirm Approval"
	if requestType == constants.PipelineRequestIDRejected {
		dialogTitle = "Confirm Rejection"
	}

	requestBody := model.OpenDialogRequest{
		TriggerId: postActionIntegrationRequest.TriggerId,
		URL:       url,
		Dialog: model.Dialog{
			Title:       dialogTitle,
			CallbackId:  postActionIntegrationRequest.PostId,
			SubmitLabel: "Submit",
			Elements:    elements,
			State:       fmt.Sprintf("%s$%s$%v$%s", organization, project, approval, requestType),
		},
	}

	if statusCode, err := p.Client.OpenDialogRequest(&requestBody, mattermostUserID); err != nil {
		p.handlePipelineApprovalRequestUpdateError("Error opening the comment dialog for user: ", mattermostUserID, err)
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	p.returnPostActionIntegrationResponse(w, &model.PostActionIntegrationResponse{})
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

	subscriptionList, err := p.Store.GetAllSubscriptions(body.MMUserID)
	if err != nil {
		p.API.LogError(constants.FetchSubscriptionListError, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	subscription, isSubscriptionPresent := p.IsSubscriptionPresent(subscriptionList, &serializers.SubscriptionDetails{
		OrganizationName:             body.Organization,
		ProjectName:                  body.Project,
		ChannelID:                    body.ChannelID,
		EventType:                    body.EventType,
		Repository:                   body.Repository,
		TargetBranch:                 body.TargetBranch,
		PullRequestCreatedBy:         body.PullRequestCreatedBy,
		PullRequestReviewersContains: body.PullRequestReviewersContains,
		PushedBy:                     body.PushedBy,
		MergeResult:                  body.MergeResult,
		NotificationType:             body.NotificationType,
		AreaPath:                     body.AreaPath,
		BuildStatus:                  body.BuildStatus,
		BuildPipeline:                body.BuildPipeline,
		StageName:                    body.StageName,
		ReleasePipeline:              body.ReleasePipeline,
		ReleaseStatus:                body.ReleaseStatus,
		ApprovalType:                 body.ApprovalType,
		ApprovalStatus:               body.ApprovalStatus,
		RunPipeline:                  body.RunPipeline,
		RunStageName:                 body.RunStageName,
		RunEnvironmentName:           body.RunEnvironmentName,
		RunStageNameID:               body.RunStageNameID,
		RunStageStateID:              body.RunStageStateID,
		RunStageResultID:             body.RunStageResultID,
		RunStateID:                   body.RunStateID,
		RunResultID:                  body.RunResultID,
	})
	if !isSubscriptionPresent {
		p.API.LogError(constants.SubscriptionNotFound)
		p.handleError(w, r, &serializers.Error{Code: http.StatusNotFound, Message: constants.SubscriptionNotFound})
		return
	}

	statusCode, deleteErr := p.deleteSubscription(subscription, mattermostUserID)
	if deleteErr != nil {
		p.API.LogError(constants.DeleteSubscriptionError, "Error", deleteErr.Error())
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: deleteErr.Error()})
		return
	}

	returnStatusOK(w)
}

func (p *Plugin) checkOAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
		azureDevopsUserID, err := p.Store.LoadAzureDevopsUserIDFromMattermostUser(mattermostUserID)
		if err != nil {
			p.API.LogError(constants.ErrorLoadingUserData, "Error", err.Error())
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: constants.GenericErrorMessage})
			return
		}

		user, err := p.Store.LoadAzureDevopsUserDetails(azureDevopsUserID)
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

func returnStatusWithMessage(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	m := map[string]string{"message": message}
	response, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
	azureDevopsUserID, err := p.Store.LoadAzureDevopsUserIDFromMattermostUser(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorLoadingDataFromKVStore, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	userDetails, err := p.Store.LoadAzureDevopsUserDetails(azureDevopsUserID)
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

func (p *Plugin) handlePipelineApproveOrRejectReleaseRequest(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	decoder := json.NewDecoder(r.Body)
	submitRequest := &model.SubmitDialogRequest{}
	if err := decoder.Decode(&submitRequest); err != nil {
		// TODO: prevent posting any error messages except oAuth in DM for now and use dialog for all such cases
		p.handlePipelineApprovalRequestUpdateError("Error decoding SubmitDialogRequest param: ", mattermostUserID, err)
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	alreadyUpdatedInformationPost := &model.Post{
		UserId:    p.botUserID,
		ChannelId: submitRequest.ChannelId,
		Message:   constants.PipelinesRequestBeingProcessed,
	}
	loaderMessagePost := p.API.SendEphemeralPost(mattermostUserID, alreadyUpdatedInformationPost)

	values := strings.Split(submitRequest.State, "$")
	var organization, projectName, requestType string
	var err error
	var approvalID float64
	if len(values) == 4 {
		organization = values[0]
		projectName = values[1]
		requestType = values[3]
		approvalID, err = strconv.ParseFloat(values[2], 64)
		if err != nil {
			p.handlePipelineApprovalRequestUpdateError(constants.GenericErrorMessage, mattermostUserID, err)
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
	}

	comments := ""
	if submitRequest.Submission[constants.DialogFieldNameComment] != nil {
		comments = submitRequest.Submission[constants.DialogFieldNameComment].(string)
	}

	pipelineApproveRequestPayload := &serializers.PipelineApproveRequest{
		Status:   requestType,
		Comments: comments,
	}

	statusCode, updatePipelineApprovalRequestErr := p.Client.UpdatePipelineApprovalRequest(pipelineApproveRequestPayload, organization, projectName, mattermostUserID, int(approvalID))
	switch statusCode {
	case http.StatusOK:
		if err := p.UpdatePipelineReleaseApprovalPost(requestType, submitRequest.CallbackId, mattermostUserID); err != nil {
			p.handlePipelineApprovalRequestUpdateError(constants.GenericErrorMessage, mattermostUserID, err)
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
	case http.StatusBadRequest:
		pipelineApprovalDetails, statusCode, err := p.Client.GetApprovalDetails(organization, projectName, mattermostUserID, int(approvalID))
		if err != nil {
			p.handlePipelineApprovalRequestUpdateError(constants.ErrorUpdatingPipelineApprovalRequest, mattermostUserID, err)
			p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
			return
		}

		if err := p.UpdatePipelineReleaseApprovalPost(pipelineApprovalDetails.Status, submitRequest.CallbackId, mattermostUserID); err != nil {
			p.handlePipelineApprovalRequestUpdateError(constants.ErrorUpdatingPipelineApprovalRequest, mattermostUserID, err)
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		alreadyUpdatedInformationPost := &model.Post{
			UserId:    p.botUserID,
			ChannelId: submitRequest.ChannelId,
			Message:   "This deployment approval pending request has already been processed.",
		}
		_ = p.API.SendEphemeralPost(mattermostUserID, alreadyUpdatedInformationPost)

	default:
		p.handlePipelineApprovalRequestUpdateError(constants.GenericErrorMessage, mattermostUserID, updatePipelineApprovalRequestErr)
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: updatePipelineApprovalRequestErr.Error()})
		return
	}

	loaderMessagePost.Message = constants.PipelinesRequestProcessed
	_ = p.API.UpdateEphemeralPost(mattermostUserID, loaderMessagePost)
	response := &model.PostActionIntegrationResponse{}
	p.returnPostActionIntegrationResponse(w, response)
}

func (p *Plugin) handlePipelineApproveOrRejectRunRequest(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	decoder := json.NewDecoder(r.Body)
	submitRequest := &model.SubmitDialogRequest{}
	if err := decoder.Decode(&submitRequest); err != nil {
		// TODO: prevent posting any error messages except oAuth in DM for now and use dialog for all such cases
		p.handlePipelineApprovalRequestUpdateError("Error decoding SubmitDialogRequest param: ", mattermostUserID, err)
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	alreadyUpdatedInformationPost := &model.Post{
		UserId:    p.botUserID,
		ChannelId: submitRequest.ChannelId,
		Message:   constants.PipelinesRequestBeingProcessed,
	}
	loaderMessagePost := p.API.SendEphemeralPost(mattermostUserID, alreadyUpdatedInformationPost)

	values := strings.Split(submitRequest.State, "$")
	var organization, projectID, approvalID, requestType string
	if len(values) == 4 {
		organization = values[0]
		projectID = values[1]
		approvalID = values[2]
		requestType = values[3]
	}

	comment := ""
	if submitRequest.Submission[constants.DialogFieldNameComment] != nil {
		comment = submitRequest.Submission[constants.DialogFieldNameComment].(string)
	}

	pipelineApproveRequestPayload := []*serializers.PipelineApproveRequest{
		{
			Status:     requestType,
			Comment:    comment,
			ApprovalID: approvalID,
		},
	}

	pipelineRunApproveResponse, statusCode, updatePipelineApprovalRequestErr := p.Client.UpdatePipelineRunApprovalRequest(pipelineApproveRequestPayload, organization, projectID, mattermostUserID)
	switch statusCode {
	case http.StatusOK:
		if err := p.UpdatePipelineRunApprovalPost(pipelineRunApproveResponse.Value[0].ApprovalSteps, pipelineRunApproveResponse.Value[0].MinRequiredApprovers, pipelineRunApproveResponse.Value[0].Status, submitRequest.CallbackId, mattermostUserID); err != nil {
			p.handlePipelineApprovalRequestUpdateError(constants.GenericErrorMessage, mattermostUserID, err)
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

	case http.StatusInternalServerError, http.StatusConflict:
		if strings.Contains(updatePipelineApprovalRequestErr.Error(), "not permitted to complete approval") || strings.Contains(updatePipelineApprovalRequestErr.Error(), "Approval is already in completed state.") {
			pipelineApprovalDetails, getApprovalDetailsStatusCode, err := p.Client.GetRunApprovalDetails(organization, projectID, mattermostUserID, approvalID)
			if err != nil {
				p.handlePipelineApprovalRequestUpdateError(constants.ErrorUpdatingPipelineApprovalRequest, mattermostUserID, err)
				p.handleError(w, r, &serializers.Error{Code: getApprovalDetailsStatusCode, Message: err.Error()})
				return
			}

			if err := p.UpdatePipelineRunApprovalPost(pipelineApprovalDetails.ApprovalSteps, pipelineApprovalDetails.MinRequiredApprovers, pipelineApprovalDetails.Status, submitRequest.CallbackId, mattermostUserID); err != nil {
				p.handlePipelineApprovalRequestUpdateError(constants.ErrorUpdatingPipelineApprovalRequest, mattermostUserID, err)
				p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
				return
			}

			alreadyUpdatedInformationPost := &model.Post{
				UserId:    p.botUserID,
				ChannelId: submitRequest.ChannelId,
				Message:   "Looks like you do not have any pending approvals or have insufficient permissions for this resource.",
			}
			if statusCode == http.StatusConflict {
				alreadyUpdatedInformationPost.Message = "Approval is already in completed state."
			}
			_ = p.API.SendEphemeralPost(mattermostUserID, alreadyUpdatedInformationPost)
		}

	default:
		p.handlePipelineApprovalRequestUpdateError(constants.GenericErrorMessage, mattermostUserID, updatePipelineApprovalRequestErr)
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: updatePipelineApprovalRequestErr.Error()})
		return
	}

	loaderMessagePost.Message = constants.PipelinesRequestProcessed
	_ = p.API.UpdateEphemeralPost(mattermostUserID, loaderMessagePost)
	response := &model.PostActionIntegrationResponse{}
	p.returnPostActionIntegrationResponse(w, response)
}

func (p *Plugin) returnPostActionIntegrationResponse(w http.ResponseWriter, res *model.PostActionIntegrationResponse) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(res.ToJson()); err != nil {
		p.API.LogWarn("Failed to write PostActionIntegrationResponse", "Error", err.Error())
	}
}

func (p *Plugin) handleGetSubscriptionFilterPossibleValues(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)

	body, err := serializers.GetSubscriptionFilterPossibleValuesRequestPayloadFromJSON(r.Body)
	if err != nil {
		p.API.LogError("Error in decoding the body for fetching subscription filter possible values", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if validationErr := body.IsSubscriptionRequestPayloadValid(); validationErr != nil {
		p.API.LogError("Request payload for fetching subscription filter possible values is not valid", "Error", validationErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: validationErr.Error()})
		return
	}

	subscriptionFilterValues, statusCode, err := p.Client.GetSubscriptionFilterPossibleValues(body, mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchSubscriptionFilterPossibleValues, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}
	filterwiseResponse := make(map[string][]*serializers.PossibleValues)
	for _, filter := range subscriptionFilterValues.InputValues {
		filterwiseResponse[filter.InputID] = filter.PossibleValues
	}

	p.writeJSON(w, filterwiseResponse)
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
