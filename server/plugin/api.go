package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"time"

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
	s.HandleFunc(constants.PathOAuthConnect, p.OAuthConnect).Methods(http.MethodGet)
	s.HandleFunc(constants.PathOAuthCallback, p.OAuthComplete).Methods(http.MethodGet)
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
	s.HandleFunc(constants.PathGetUserChannelsForTeam, p.handleAuthRequired(p.getUserChannelsForTeam)).Methods(http.MethodGet)
	s.HandleFunc(constants.PathGetGitRepositories, p.handleAuthRequired(p.checkOAuth(p.handleGetGitRepositories))).Methods(http.MethodGet)
	s.HandleFunc(constants.PathGetGitRepositoryBranches, p.handleAuthRequired(p.checkOAuth(p.handleGetGitRepositoryBranches))).Methods(http.MethodGet)
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

	if _, isProjectLinked := p.IsProjectLinked(projectList, serializers.ProjectDetails{OrganizationName: strings.ToLower(body.Organization), ProjectName: cases.Title(language.Und).String(body.Project)}); isProjectLinked {
		returnStatusWithMessage(w, http.StatusOK, constants.AlreadyLinkedProject)
		return
	}

	response, statusCode, err := p.Client.Link(body, mattermostUserID)
	if err != nil {
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	isAdmin := false
	subscriptionStatusCode, subscriptionErr := p.Client.CheckIfUserIsProjectAdmin(body.Organization, response.ID, p.GetPluginURL(), mattermostUserID)
	if subscriptionErr != nil {
		switch {
		case subscriptionStatusCode == http.StatusBadRequest && strings.Contains(subscriptionErr.Error(), fmt.Sprintf(constants.ErrorMessageForAdmin, constants.SubscriptionEventTypeDummy)):
			isAdmin = true
		case subscriptionStatusCode == http.StatusForbidden && strings.Contains(subscriptionErr.Error(), constants.AccessDenied):
			isAdmin = false
		default:
			p.API.LogError(fmt.Sprintf(constants.ErrorCheckingProjectAdmin, body.Project), "Error", subscriptionErr.Error())
			p.handleError(w, r, &serializers.Error{Code: subscriptionStatusCode, Message: constants.ErrorLinkProject})
			return
		}
	}

	project := serializers.ProjectDetails{
		MattermostUserID: mattermostUserID,
		ProjectID:        response.ID,
		ProjectName:      cases.Title(language.Und).String(body.Project),
		OrganizationName: strings.ToLower(body.Organization),
		IsAdmin:          isAdmin,
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
	}); isSubscriptionPresent {
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
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: constants.GetChannelError})
		return
	}

	user, userErr := p.API.GetUser(mattermostUserID)
	if userErr != nil {
		p.API.LogError(constants.GetUserError, "Error", userErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: constants.GetUserError})
		return
	}

	createdByDisplayName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	if len(strings.TrimSpace(createdByDisplayName)) == 0 {
		createdByDisplayName = user.Username // If user's first/last name doesn't exist then show username as fallback
	}
	if storeErr := p.Store.StoreSubscription(&serializers.SubscriptionDetails{
		MattermostUserID:                 mattermostUserID,
		ProjectName:                      body.Project,
		ProjectID:                        project.ProjectID,
		OrganizationName:                 body.Organization,
		EventType:                        body.EventType,
		ServiceType:                      body.ServiceType,
		ChannelID:                        body.ChannelID,
		SubscriptionID:                   subscription.ID,
		ChannelName:                      channel.DisplayName,
		ChannelType:                      channel.Type,
		CreatedBy:                        createdByDisplayName,
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
		p.API.LogError("Invalid team id")
		http.Error(w, "Invalid team id", http.StatusBadRequest)
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
		p.API.LogError(constants.FetchSubscriptionListError, "Error", subscriptionErr.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: subscriptionErr.Error()})
		return
	}

	offset, limit := p.GetOffsetAndLimitFromQueryParams(r)
	channelID := r.URL.Query().Get(constants.QueryParamChannelID)
	serviceType := r.URL.Query().Get(constants.QueryParamServiceType)
	eventType := r.URL.Query().Get(constants.QueryParamEventType)
	project := r.URL.Query().Get(constants.QueryParamProject)
	if project != "" {
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
					}
				}
			}
		}

		sort.Slice(subscriptionByProject, func(i, j int) bool {
			return subscriptionByProject[i].ChannelName+
				subscriptionByProject[i].EventType+
				subscriptionByProject[i].TargetBranch+
				subscriptionByProject[i].PullRequestCreatedByName+
				subscriptionByProject[i].PullRequestReviewersContainsName+
				subscriptionByProject[i].PushedByName+
				subscriptionByProject[i].MergeResultName+
				subscriptionByProject[i].NotificationTypeName+
				subscriptionByProject[i].AreaPath <
				subscriptionByProject[j].ChannelName+
					subscriptionByProject[j].EventType+
					subscriptionByProject[j].TargetBranch+
					subscriptionByProject[j].PullRequestCreatedByName+
					subscriptionByProject[j].PullRequestReviewersContainsName+
					subscriptionByProject[j].PushedByName+
					subscriptionByProject[j].MergeResultName+
					subscriptionByProject[j].NotificationTypeName+
					subscriptionByProject[j].AreaPath
		})

		filteredSubscriptionList, filteredSubscriptionErr := p.GetSubscriptionsForAccessibleChannelsOrProjects(subscriptionByProject, teamID, mattermostUserID)
		if filteredSubscriptionErr != nil {
			p.API.LogError(constants.FetchFilteredSubscriptionListError, "Error", filteredSubscriptionErr.Error())
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

		subscriptionList = paginatedSubscriptions
	}

	p.writeJSON(w, subscriptionList)
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

	var attachment *model.SlackAttachment
	switch body.EventType {
	case constants.SubscriptionEventWorkItemCreated, constants.SubscriptionEventWorkItemUpdated, constants.SubscriptionEventWorkItemDeleted, constants.SubscriptionEventWorkItemCommented:
		attachment = &model.SlackAttachment{
			AuthorName: constants.SlackAttachmentAuthorNameBoards,
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameBoardsIcon),
			Color:      constants.IconColorBoards,
			Text:       body.DetailedMessage.Markdown,
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
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameReposIcon),
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
			FooterIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
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
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameReposIcon),
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
			FooterIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
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
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameReposIcon),
			Color:      constants.IconColorRepos,
			Title:      "Commit(s)",
			Text:       commits,
			Footer:     fmt.Sprintf("%s | %s", strings.Split(body.Resource.RefUpdates[0].Name, "/")[2], body.Resource.Repository.Name),
			FooterIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameGitBranchIcon),
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
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
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
			FooterIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
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
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
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
			FooterIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
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
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
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
			FooterIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventReleaseDeploymentStarted:
		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
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
			FooterIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventReleaseDeploymentCompleted:
		comment := body.Resource.Comment.(string)
		if comment == "" {
			comment = "No comments"
		}

		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
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
			FooterIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
		}
	case constants.SubscriptionEventRunStageStateChanged:
		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Pipeline.Name, body.Resource.Stage.Links.PipelineWeb.Href),
					Short: true,
				},
			},
		}
	case constants.SubscriptionEventRunStateChanged:
		attachment = &model.SlackAttachment{
			Pretext:    body.Message.Markdown,
			AuthorName: constants.SlackAttachmentAuthorNamePipelines,
			AuthorIcon: fmt.Sprintf(constants.StaticFiles, p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon),
			Color:      constants.IconColorPipelines,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Pipeline",
					Value: fmt.Sprintf("[%s](%s)", body.Resource.Pipeline.Name, body.Resource.Run.Links.PipelineWeb.Href),
					Short: true,
				},
			},
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

func (p *Plugin) deleteSubscription(subscription *serializers.SubscriptionDetails, mattermostUserID string) (int, error) {
	if statusCode, err := p.Client.DeleteSubscription(subscription.OrganizationName, subscription.SubscriptionID, mattermostUserID); err != nil {
		return statusCode, err
	}

	if deleteErr := p.Store.DeleteSubscription(subscription); deleteErr != nil {
		return http.StatusInternalServerError, deleteErr
	}

	return http.StatusOK, nil
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

func (p *Plugin) handleGetGitRepositories(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)

	pathParams := mux.Vars(r)
	organization := pathParams[constants.PathParamOrganization]
	project := pathParams[constants.PathParamProject]

	if len(strings.TrimSpace(organization)) == 0 || len(strings.TrimSpace(project)) == 0 {
		p.API.LogError(constants.ErrorOrganizationOrProjectQueryParam)
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.ErrorOrganizationOrProjectQueryParam})
		return
	}

	response, statusCode, err := p.Client.GetGitRepositories(organization, project, mattermostUserID)
	if err != nil {
		p.API.LogError("Error in fetching git repositories", err.Error())
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	p.writeJSON(w, response.Value)
}

func (p *Plugin) handleGetGitRepositoryBranches(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)

	pathParams := mux.Vars(r)
	organization := strings.TrimSpace(pathParams[constants.PathParamOrganization])
	project := strings.TrimSpace(pathParams[constants.PathParamProject])
	repository := strings.TrimSpace(pathParams[constants.PathParamRepository])

	if len(organization) == 0 || len(project) == 0 || len(repository) == 0 {
		p.API.LogError(constants.ErrorRepositoryPathParam)
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.ErrorRepositoryPathParam})
		return
	}

	response, statusCode, err := p.Client.GetGitRepositoryBranches(organization, project, repository, mattermostUserID)
	if err != nil {
		p.API.LogError("Error in fetching git repository branches", err.Error())
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}

	// Azure DevOps returns branch name as "refs/heads/<branch-name>", but we need to use only "<branch-name>" so, remove unused part from the name
	for _, value := range response.Value {
		if strings.Contains(value.Name, "refs/heads/") && len(value.Name) > 11 {
			value.Name = value.Name[11:]
		}
	}

	p.writeJSON(w, response.Value)
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
