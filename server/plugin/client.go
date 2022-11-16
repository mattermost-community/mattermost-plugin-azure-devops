package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/pkg/errors"
)

type Client interface {
	GenerateOAuthToken(encodedFormValues url.Values) (*serializers.OAuthSuccessResponse, int, error)
	CreateTask(body *serializers.CreateTaskRequestPayload, mattermostUserID string) (*serializers.TaskValue, int, error)
	GetTask(organization, taskID, projectName, mattermostUserID string) (*serializers.TaskValue, int, error)
	GetPullRequest(organization, pullRequestID, projectName, mattermostUserID string) (*serializers.PullRequest, int, error)
	Link(body *serializers.LinkRequestPayload, mattermostUserID string) (*serializers.Project, int, error)
	CreateSubscription(body *serializers.CreateSubscriptionRequestPayload, project *serializers.ProjectDetails, channelID, pluginURL, mattermostUserID string) (*serializers.SubscriptionValue, int, error)
	DeleteSubscription(organization, subscriptionID, mattermostUserID string) (int, error)
	CheckIfUserIsProjectAdmin(organizationName, projectID, pluginURL, mattermostUserID string) (int, error)
	GetGitRepositories(organization, projectName, mattermostUserID string) (*serializers.GitRepositoriesResponse, int, error)
	GetGitRepositoryBranches(organization, projectName, repository, mattermostUserID string) (*serializers.GitBranchesResponse, int, error)
	GetSubscriptionFilterPossibleValues(request *serializers.GetSubscriptionFilterPossibleValuesRequestPayload, mattermostUserID string) (*serializers.SubscriptionFilterPossibleValuesResponseFromClient, int, error)
}

type client struct {
	plugin     *Plugin
	httpClient *http.Client
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (c *client) GenerateOAuthToken(encodedFormValues url.Values) (*serializers.OAuthSuccessResponse, int, error) {
	var oAuthSuccessResponse *serializers.OAuthSuccessResponse

	_, statusCode, err := c.callFormURLEncoded(constants.BaseOauthURL, constants.PathToken, http.MethodPost, &oAuthSuccessResponse, encodedFormValues)
	if err != nil {
		return nil, statusCode, err
	}

	return oAuthSuccessResponse, statusCode, nil
}

// Function to create task for a project.
func (c *client) CreateTask(body *serializers.CreateTaskRequestPayload, mattermostUserID string) (*serializers.TaskValue, int, error) {
	taskURL := fmt.Sprintf(constants.CreateTask, body.Organization, body.Project, body.Type)

	// Create request body.
	payload := []*serializers.CreateTaskBodyPayload{}
	payload = append(payload,
		&serializers.CreateTaskBodyPayload{
			Operation: "add",
			Path:      "/fields/System.Title",
			From:      "",
			Value:     body.Fields.Title,
		})

	if body.Fields.Description != "" {
		payload = append(payload,
			&serializers.CreateTaskBodyPayload{
				Operation: "add",
				Path:      "/fields/System.Description",
				From:      "",
				Value:     body.Fields.Description,
			})
	}
	if body.Fields.AreaPath != "" {
		payload = append(payload,
			&serializers.CreateTaskBodyPayload{
				Operation: "add",
				Path:      "/fields/System.AreaPath",
				From:      "",
				Value:     body.Fields.AreaPath,
			})
	}

	var task *serializers.TaskValue
	_, statusCode, err := c.CallPatchJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskURL, http.MethodPost, mattermostUserID, &payload, &task, nil)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to create task")
	}

	return task, statusCode, nil
}

// Function to get the task.
func (c *client) GetTask(organization, taskID, projectName, mattermostUserID string) (*serializers.TaskValue, int, error) {
	taskURL := fmt.Sprintf(constants.GetTask, organization, projectName, taskID)

	var task *serializers.TaskValue
	_, statusCode, err := c.CallJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskURL, http.MethodGet, mattermostUserID, nil, &task, nil)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to get the Task")
	}

	return task, statusCode, nil
}

// Function to get the pull request.
func (c *client) GetPullRequest(organization, pullRequestID, projectName, mattermostUserID string) (*serializers.PullRequest, int, error) {
	pullRequestURL := fmt.Sprintf(constants.GetPullRequest, organization, projectName, pullRequestID)

	var pullRequest *serializers.PullRequest
	_, statusCode, err := c.CallJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, pullRequestURL, http.MethodGet, mattermostUserID, nil, &pullRequest, nil)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to get the pull request")
	}

	return pullRequest, statusCode, nil
}

// Function to link a project and an organization.
func (c *client) Link(body *serializers.LinkRequestPayload, mattermostUserID string) (*serializers.Project, int, error) {
	projectURL := fmt.Sprintf(constants.GetProject, body.Organization, body.Project)
	var project *serializers.Project

	_, statusCode, err := c.CallJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, projectURL, http.MethodGet, mattermostUserID, nil, &project, nil)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to link Project")
	}
	return project, statusCode, nil
}

// Wrapper to make REST API requests with "application/x-www-form-urlencoded" type content
func (c *client) callFormURLEncoded(url, path, method string, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
	contentType := "application/x-www-form-urlencoded"
	return c.Call(url, method, path, contentType, "", nil, out, formValues)
}

var publisherID = map[string]string{
	constants.SubscriptionEventPullRequestCreated:                 constants.PublisherIDTFS,
	constants.SubscriptionEventPullRequestUpdated:                 constants.PublisherIDTFS,
	constants.SubscriptionEventPullRequestCommented:               constants.PublisherIDTFS,
	constants.SubscriptionEventPullRequestMerged:                  constants.PublisherIDTFS,
	constants.SubscriptionEventCodePushed:                         constants.PublisherIDTFS,
	constants.SubscriptionEventWorkItemCreated:                    constants.PublisherIDTFS,
	constants.SubscriptionEventWorkItemUpdated:                    constants.PublisherIDTFS,
	constants.SubscriptionEventWorkItemDeleted:                    constants.PublisherIDTFS,
	constants.SubscriptionEventWorkItemCommented:                  constants.PublisherIDTFS,
	constants.SubscriptionEventBuildCompleted:                     constants.PublisherIDTFS,
	constants.SubscriptionEventReleaseAbandoned:                   constants.PublisherIDRM,
	constants.SubscriptionEventReleaseCreated:                     constants.PublisherIDRM,
	constants.SubscriptionEventReleaseDeploymentApprovalCompleted: constants.PublisherIDRM,
	constants.SubscriptionEventReleaseDeploymentEventPending:      constants.PublisherIDRM,
	constants.SubscriptionEventReleaseDeploymentCompleted:         constants.PublisherIDRM,
	constants.SubscriptionEventReleaseDeploymentStarted:           constants.PublisherIDRM,
	constants.SubscriptionEventRunStageApprovalCompleted:          constants.PublisherIDPipelines,
	constants.SubscriptionEventRunStageStateChanged:               constants.PublisherIDPipelines,
	constants.SubscriptionEventRunStageWaitingForApproval:         constants.PublisherIDPipelines,
	constants.SubscriptionEventRunStateChanged:                    constants.PublisherIDPipelines,
}

func (c *client) CreateSubscription(body *serializers.CreateSubscriptionRequestPayload, project *serializers.ProjectDetails, channelID, pluginURL, mattermostUserID string) (*serializers.SubscriptionValue, int, error) {
	subscriptionURL := fmt.Sprintf(constants.CreateSubscription, body.Organization)

	consumerInputs := serializers.ConsumerInputs{
		URL: fmt.Sprintf("%s%s?channelID=%s", strings.TrimRight(pluginURL, "/"), constants.PathSubscriptionNotifications, channelID),
	}

	payload := serializers.CreateSubscriptionBodyPayload{
		PublisherID:      publisherID[body.EventType],
		EventType:        body.EventType,
		ConsumerID:       constants.ConsumerID,
		ConsumerActionID: constants.ConsumerActionID,
		ConsumerInputs:   consumerInputs,
		PublisherInputs: serializers.PublisherInputsGeneric{
			ProjectID:                    project.ProjectID,
			AreaPath:                     body.AreaPath,
			Repository:                   body.Repository,
			Branch:                       body.TargetBranch,
			PushedBy:                     body.PushedBy,
			MergeResult:                  body.MergeResult,
			PullRequestCreatedBy:         body.PullRequestCreatedBy,
			PullRequestReviewersContains: body.PullRequestReviewersContains,
			NotificationType:             body.NotificationType,
		},
	}

	baseURL := c.plugin.getConfiguration().AzureDevopsAPIBaseURL
	if strings.Contains(body.EventType, "release") {
		baseURL = strings.Replace(baseURL, "://", "://vsrm.", 1)
	}

	var subscription *serializers.SubscriptionValue
	_, statusCode, err := c.CallJSON(baseURL, subscriptionURL, http.MethodPost, mattermostUserID, payload, &subscription, nil)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to create subscription")
	}

	return subscription, statusCode, nil
}

// We are passing an invalid request payload for creating a subscription to check if the user who is linking this project is an admin here or not
// If the user is an admin then we will get a response status code 400 with a message of invalid payload
// and if the user is not an admin we will get status code 403 with a message saying "Access Denied"
func (c *client) CheckIfUserIsProjectAdmin(organizationName, projectID, pluginURL, mattermostUserID string) (int, error) {
	subscriptionURL := fmt.Sprintf(constants.CreateSubscription, organizationName)

	publisherInputs := serializers.PublisherInputsGeneric{
		ProjectID: projectID,
	}

	consumerInputs := serializers.ConsumerInputs{
		URL: fmt.Sprintf("%s%s?channelID=%s", strings.TrimRight(pluginURL, "/"), constants.PathSubscriptionNotifications, constants.SubscriptionEventTypeDummy),
	}

	payload := serializers.CreateSubscriptionBodyPayload{
		PublisherID:      constants.PublisherIDTFS,
		EventType:        constants.SubscriptionEventTypeDummy,
		ConsumerID:       constants.ConsumerID,
		ConsumerActionID: constants.ConsumerActionID,
		PublisherInputs:  publisherInputs,
		ConsumerInputs:   consumerInputs,
	}

	_, statusCode, err := c.CallJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, subscriptionURL, http.MethodPost, mattermostUserID, payload, nil, nil)
	if err != nil {
		return statusCode, errors.Wrap(err, "failed to create subscription")
	}

	return statusCode, nil
}

func (c *client) DeleteSubscription(organization, subscriptionID, mattermostUserID string) (int, error) {
	subscriptionURL := fmt.Sprintf(constants.DeleteSubscription, organization, subscriptionID)
	_, statusCode, err := c.CallJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, subscriptionURL, http.MethodDelete, mattermostUserID, nil, nil, nil)
	if err != nil {
		return statusCode, errors.Wrap(err, "failed to delete subscription")
	}

	return statusCode, nil
}

func (c *client) GetGitRepositories(organization, projectName, mattermostUserID string) (*serializers.GitRepositoriesResponse, int, error) {
	getGitRepositoriesURL := fmt.Sprintf(constants.GetGitRepositories, organization, projectName)

	var gitRepositories *serializers.GitRepositoriesResponse
	_, statusCode, err := c.CallJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, getGitRepositoriesURL, http.MethodGet, mattermostUserID, nil, &gitRepositories, nil)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to get the git repositories for a project")
	}

	return gitRepositories, statusCode, nil
}

func (c *client) GetGitRepositoryBranches(organization, projectName, repository, mattermostUserID string) (*serializers.GitBranchesResponse, int, error) {
	getGitRepositoriesURL := fmt.Sprintf(constants.GetGitRepositoryBranches, organization, projectName, repository)

	var gitBranchesResponse *serializers.GitBranchesResponse
	_, statusCode, err := c.CallJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, getGitRepositoriesURL, http.MethodGet, mattermostUserID, nil, &gitBranchesResponse, nil)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to get the git repository branches for a project")
	}

	return gitBranchesResponse, statusCode, nil
}

func (c *client) GetSubscriptionFilterPossibleValues(request *serializers.GetSubscriptionFilterPossibleValuesRequestPayload, mattermostUserID string) (*serializers.SubscriptionFilterPossibleValuesResponseFromClient, int, error) {
	getSubscriptionFilterValuesURL := fmt.Sprintf(constants.GetSubscriptionFilterPossibleValues, request.Organization)

	var subscriptionFilters []*serializers.SubscriptionFilter
	for _, filter := range request.Filters {
		subscriptionFilters = append(subscriptionFilters, &serializers.SubscriptionFilter{InputID: filter})
	}

	subscriptionFiltersRequest := &serializers.GetSubscriptionFilterValuesRequestPayloadFromClient{
		Subscription: &serializers.CreateSubscriptionBodyPayload{
			PublisherID:      constants.PublisherIDTFS,
			ConsumerID:       constants.ConsumerID,
			ConsumerActionID: constants.ConsumerActionID,
			EventType:        request.EventType,
			PublisherInputs: serializers.PublisherInputsGeneric{
				ProjectID: request.ProjectID,
			},
		},
		InputValues: subscriptionFilters,
		Scope:       10, // TODO: This is a required field for Azure DevOps and must have value 10, it's use or role is not documented anywhere in the Azure DevOps API docs so, it can be investigated further for more details
	}

	if constants.ValidSubscriptionEventsForRepos[request.EventType] {
		subscriptionFiltersRequest.Subscription.PublisherInputs = serializers.PublisherInputsGeneric{
			ProjectID:  request.ProjectID,
			Repository: request.RepositoryID,
		}
	}

	var subscriptionFiltersResponse *serializers.SubscriptionFilterPossibleValuesResponseFromClient
	_, statusCode, err := c.CallJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, getSubscriptionFilterValuesURL, http.MethodPost, mattermostUserID, &subscriptionFiltersRequest, &subscriptionFiltersResponse, nil)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to get the subscription filter values")
	}

	return subscriptionFiltersResponse, statusCode, nil
}

// Wrapper to make REST API requests with "application/json-patch+json" type content
func (c *client) CallPatchJSON(url, path, method, mattermostUserID string, in, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
	contentType := "application/json-patch+json"
	buf := &bytes.Buffer{}
	if err = json.NewEncoder(buf).Encode(in); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return c.Call(url, method, path, contentType, mattermostUserID, buf, out, formValues)
}

// Wrapper to make REST API requests with "application/json" type content
func (c *client) CallJSON(url, path, method, mattermostUserID string, in, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
	contentType := "application/json"
	buf := &bytes.Buffer{}
	if err = json.NewEncoder(buf).Encode(in); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return c.Call(url, method, path, contentType, mattermostUserID, buf, out, formValues)
}

// Makes HTTP request to REST APIs
func (c *client) Call(basePath, method, path, contentType string, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, statusCode int, err error) {
	errContext := fmt.Sprintf("Azure DevOps: Call failed: method:%s, path:%s", method, path)
	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.WithMessage(err, errContext)
	}

	if pathURL.Scheme == "" || pathURL.Host == "" {
		var baseURL *url.URL
		baseURL, err = url.Parse(basePath)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.WithMessage(err, errContext)
		}
		if path[0] != '/' {
			path = "/" + path
		}
		path = baseURL.String() + path
	}

	// Check refresh token only for APIs other than OAuth
	if basePath != constants.BaseOauthURL {
		if isAccessTokenExpired, refreshToken := c.plugin.IsAccessTokenExpired(mattermostUserID); isAccessTokenExpired {
			if errRefreshingToken := c.plugin.RefreshOAuthToken(mattermostUserID, refreshToken); errRefreshingToken != nil {
				return nil, http.StatusInternalServerError, errRefreshingToken
			}
		}
	}

	var req *http.Request
	if formValues != nil {
		req, err = http.NewRequest(method, path, strings.NewReader(formValues.Encode()))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	} else {
		req, err = http.NewRequest(method, path, inBody)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if mattermostUserID != "" {
		if err = c.plugin.AddAuthorization(req, mattermostUserID); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if resp.Body == nil {
		return nil, resp.StatusCode, nil
	}
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if out != nil {
			if err = json.Unmarshal(responseData, out); err != nil {
				return responseData, http.StatusInternalServerError, err
			}
		}
		return responseData, resp.StatusCode, nil

	case http.StatusNoContent:
		return nil, resp.StatusCode, nil

	case http.StatusNotFound:
		return nil, resp.StatusCode, ErrNotFound
	}

	errResp := ErrorResponse{}
	if err = json.Unmarshal(responseData, &errResp); err != nil {
		return responseData, http.StatusInternalServerError, errors.WithMessagef(err, "status: %s", resp.Status)
	}
	return responseData, resp.StatusCode, fmt.Errorf("errorMessage %s", errResp.Message)
}

func InitClient(p *Plugin) Client {
	return &client{
		plugin:     p,
		httpClient: &http.Client{},
	}
}
