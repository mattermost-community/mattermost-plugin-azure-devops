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
	TestApi() (string, error) // TODO: remove later
	GenerateOAuthToken(formValues url.Values) (*serializers.OAuthSuccessResponse, int, error)
	CreateTask(body *serializers.TaskCreateRequestPayload, mattermostUserID string) (*serializers.TaskValue, int, error)
	GetTask(queryParams serializers.GetTaskData, mattermostUserID string) (*serializers.TaskValue, int, error)
	Link(body *serializers.LinkRequestPayload, mattermostUserID string) (*serializers.Project, int, error)
	CreateSubscription(body *serializers.CreateSubscriptionRequestPayload, project *serializers.ProjectDetails, channelID, pluginURL, mattermostUserID string) (*serializers.SubscriptionValue, int, error)
	DeleteSubscription(organization, subscriptionID, mattermostUserID string) (int, error)
}

type client struct {
	plugin     *Plugin
	httpClient *http.Client
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// TODO: remove later
func (c *client) TestApi() (string, error) {
	return "hello world", nil
}

// Function to create task for a project.
func (c *client) CreateTask(body *serializers.TaskCreateRequestPayload, mattermostUserID string) (*serializers.TaskValue, int, error) {
	taskURL := fmt.Sprintf(constants.CreateTask, body.Organization, body.Project, body.Type)

	// Create request body.
	payload := []*serializers.TaskCreateBodyPayload{}
	payload = append(payload,
		&serializers.TaskCreateBodyPayload{
			Operation: "add",
			Path:      "/fields/System.Title",
			From:      "",
			Value:     body.Fields.Title,
		})

	if body.Fields.Description != "" {
		payload = append(payload,
			&serializers.TaskCreateBodyPayload{
				Operation: "add",
				Path:      "/fields/System.Description",
				From:      "",
				Value:     body.Fields.Description,
			})
	}

	var task *serializers.TaskValue
	_, statusCode, err := c.callPatchJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskURL, http.MethodPost, mattermostUserID, payload, &task, false)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to create the Task")
	}

	return task, statusCode, nil
}

func (c *client) GenerateOAuthToken(formValues url.Values) (*serializers.OAuthSuccessResponse, int, error) {
	var oAuthSuccessResponse *serializers.OAuthSuccessResponse

	_, statusCode, err := c.callFormURLEncoded(constants.BaseOauthURL, constants.PathToken, "", http.MethodPost, &oAuthSuccessResponse, formValues, false)
	if err != nil {
		return nil, statusCode, err
	}

	return oAuthSuccessResponse, statusCode, nil
}

// Function to get the task.
func (c *client) GetTask(queryParams serializers.GetTaskData, mattermostUserID string) (*serializers.TaskValue, int, error) {
	taskURL := fmt.Sprintf(constants.GetTask, queryParams.Organization, queryParams.TaskID)

	var task *serializers.TaskValue
	_, statusCode, err := c.callJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskURL, http.MethodGet, mattermostUserID, nil, &task, false)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to get the Task")
	}

	return task, statusCode, nil
}

// Function to link a project and an organization.
func (c *client) Link(body *serializers.LinkRequestPayload, mattermostUserID string) (*serializers.Project, int, error) {
	projectURL := fmt.Sprintf(constants.GetProject, body.Organization, body.Project)
	var project *serializers.Project

	_, statusCode, err := c.callJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, projectURL, http.MethodGet, mattermostUserID, nil, &project, false)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to link Project")
	}

	return project, statusCode, nil
}

func (c *client) CreateSubscription(body *serializers.CreateSubscriptionRequestPayload, project *serializers.ProjectDetails, channelID, pluginURL, mattermostUserID string) (*serializers.SubscriptionValue, int, error) {
	subscriptionURL := fmt.Sprintf(constants.CreateSubscription, body.Organization)

	publisherInputs := serializers.PublisherInputs{
		ProjectID: project.ProjectID,
	}

	consumerInputs := serializers.ConsumerInputs{
		URL: fmt.Sprintf("%s%s?channelID=%s", strings.TrimRight(pluginURL, "/"), constants.PathNotificationSubscriptions, channelID),
	}

	StatusData := map[string]string{
		constants.Create: "workitem.created",
		constants.Update: "workitem.updated",
		constants.Delete: "workitem.deleted",
	}

	payload := serializers.CreateSubscriptionBodyPayload{
		PublisherID:      constants.PublisherID,
		EventType:        StatusData[body.EventType],
		ConsumerId:       constants.ConsumerId,
		ConsumerActionId: constants.ConsumerActionId,
		PublisherInputs:  publisherInputs,
		ConsumerInputs:   consumerInputs,
	}
	var subscription *serializers.SubscriptionValue
	_, statusCode, err := c.callJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, subscriptionURL, http.MethodPost, mattermostUserID, payload, &subscription, true)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to create subscription")
	}

	return subscription, statusCode, nil
}

func (c *client) DeleteSubscription(organization, subscriptionID, mattermostUserID string) (int, error) {
	subscriptionURL := fmt.Sprintf(constants.DeleteSubscription, organization, subscriptionID)
	_, statusCode, err := c.callJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, subscriptionURL, http.MethodDelete, mattermostUserID, nil, nil, true)
	if err != nil {
		return statusCode, errors.Wrap(err, "failed to delete subscription")
	}

	return statusCode, nil
}

// Wrapper to make REST API requests with "application/json-patch+json" type content
func (c *client) callPatchJSON(url, path, method, mattermostUserID string, in, out interface{}, basicAuth bool) (responseData []byte, statusCode int, err error) {
	contentType := "application/json-patch+json"
	buf := &bytes.Buffer{}
	if err = json.NewEncoder(buf).Encode(in); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return c.call(url, method, path, contentType, mattermostUserID, buf, out, nil, basicAuth)
}

// Wrapper to make REST API requests with "application/json" type content
func (c *client) callJSON(url, path, method string, mattermostUserID string, in, out interface{}, basicAuth bool) (responseData []byte, statusCode int, err error) {
	contentType := "application/json"
	buf := &bytes.Buffer{}
	if err = json.NewEncoder(buf).Encode(in); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return c.call(url, method, path, contentType, mattermostUserID, buf, out, nil, basicAuth)
}

// Wrapper to make REST API requests with "application/x-www-form-urlencoded" type content
func (c *client) callFormURLEncoded(url, path, mattermostUserID string, method string, out interface{}, formValues url.Values, basicAuth bool) (responseData []byte, statusCode int, err error) {
	contentType := "application/x-www-form-urlencoded"
	return c.call(url, method, path, contentType, mattermostUserID, nil, out, formValues, basicAuth)
}

// Makes HTTP request to REST APIs
func (c *client) call(basePath, method, path, contentType string, mattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values, basicAuth bool) (responseData []byte, statusCode int, err error) {
	errContext := fmt.Sprintf("Azure Devops: Call failed: method:%s, path:%s", method, path)
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

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if mattermostUserID != "" {
		if basicAuth {
			if err = c.plugin.AddBasicAuthorization(req, mattermostUserID); err != nil {
				return nil, http.StatusInternalServerError, err
			}
		} else {
			if err = c.plugin.AddAuthorization(req, mattermostUserID); err != nil {
				return nil, http.StatusInternalServerError, err
			}
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	if resp.Body == nil {
		return nil, resp.StatusCode, nil
	}
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusNonAuthoritativeInfo:
		if err := c.plugin.RefreshOAuthToken(mattermostUserID); err != nil {
			return nil, http.StatusUnauthorized, err
		}
		_, statusCode, err := c.call(basePath, method, path, contentType, mattermostUserID, inBody, out, formValues, basicAuth)
		return nil, statusCode, err
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
