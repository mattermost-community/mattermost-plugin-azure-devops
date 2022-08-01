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
	GenerateOAuthToken(encodedFormValues string) (*serializers.OAuthSuccessResponse, error)
	CreateTask(body *serializers.TaskCreateRequestPayload, mattermostUserID string) (*serializers.TaskValue, error)
	GetTask(queryParams serializers.GetTaskData, mattermostUserID string) (*serializers.TaskValue, error)
	Link(body *serializers.LinkRequestPayload, mattermostUserID string) (*serializers.Project, error)
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

func (c *client) GenerateOAuthToken(encodedFormValues string) (*serializers.OAuthSuccessResponse, error) {
	var oAuthSuccessResponse *serializers.OAuthSuccessResponse

	if _, err := c.callFormURLEncoded(constants.BaseOauthURL, constants.PathToken, http.MethodPost, &oAuthSuccessResponse, encodedFormValues); err != nil {
		return nil, err
	}

	return oAuthSuccessResponse, nil
}

// Function to create task of a project.
func (azureDevops *client) CreateTask(body *serializers.TaskCreateRequestPayload, mattermostUserID string) (*serializers.TaskValue, error) {
	taskURL := fmt.Sprintf(constants.CreateTask, body.Organization, body.Project, body.Type)

	// Create payload body to send.
	payload := []serializers.TaskCreateBodyPayload{}
	payload = append(payload,
		serializers.TaskCreateBodyPayload{
			Operation: "add",
			Path:      "/fields/System.Title",
			From:      "",
			Value:     body.Feilds.Title,
		})

	if body.Feilds.Description != "" {
		payload = append(payload,
			serializers.TaskCreateBodyPayload{
				Operation: "add",
				Path:      "/fields/System.Description",
				From:      "",
				Value:     body.Feilds.Description,
			})
	}
	var task *serializers.TaskValue
	if _, err := azureDevops.callPatchJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskURL, http.MethodPost, mattermostUserID, payload, &task); err != nil {
		return nil, errors.Wrap(err, "failed to get create Task")
	}

	return task, nil
}

// Function to get the task.
func (c *client) GetTask(queryParams serializers.GetTaskData, mattermostUserID string) (*serializers.TaskValue, error) {
	taskURL := fmt.Sprintf(constants.GetTask, queryParams.Organization, queryParams.TaskID)

	var task *serializers.TaskValue
	if _, err := c.callJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskURL, http.MethodGet, mattermostUserID, nil, &task); err != nil {
		return nil, errors.Wrap(err, "failed to get the Task")
	}

	return task, nil
}

// Function to link a project and an organization.
func (c *client) Link(body *serializers.LinkRequestPayload, mattermostUserID string) (*serializers.Project, error) {
	projectURL := fmt.Sprintf(constants.GetProject, body.Organization, body.Project)
	var project *serializers.Project

	if _, err := c.callJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, projectURL, http.MethodGet, mattermostUserID, nil, &project); err != nil {
		return nil, errors.Wrap(err, "failed to link Project")
	}

	return project, nil
}

// Wrapper to make REST API requests with "application/json-patch+json" type content
func (c *client) callPatchJSON(url, path, method, mattermostUserID string, in, out interface{}) (responseData []byte, err error) {
	contentType := "application/json-patch+json"
	buf := &bytes.Buffer{}
	if err = json.NewEncoder(buf).Encode(in); err != nil {
		return nil, err
	}
	return c.call(url, method, path, contentType, mattermostUserID, buf, out, "")
}

// Wrapper to make REST API requests with "application/json" type content
func (c *client) callJSON(url, path, method, mattermostUserID string, in, out interface{}) (responseData []byte, err error) {
	contentType := "application/json"
	buf := &bytes.Buffer{}
	if err = json.NewEncoder(buf).Encode(in); err != nil {
		return nil, err
	}
	return c.call(url, method, path, contentType, mattermostUserID, buf, out, "")
}

// Wrapper to make REST API requests with "application/x-www-form-urlencoded" type content
func (c *client) callFormURLEncoded(url, path, method string, out interface{}, formValues string) (responseData []byte, err error) {
	contentType := "application/x-www-form-urlencoded"
	return c.call(url, method, path, contentType, "", nil, out, formValues)
}

// Makes HTTP request to REST APIs
func (c *client) call(basePath, method, path, contentType string, mattermostUserID string, inBody io.Reader, out interface{}, formValues string) (responseData []byte, err error) {
	errContext := fmt.Sprintf("Azure Devops: Call failed: method:%s, path:%s", method, path)
	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, errors.WithMessage(err, errContext)
	}

	if pathURL.Scheme == "" || pathURL.Host == "" {
		var baseURL *url.URL
		baseURL, err = url.Parse(basePath)
		if err != nil {
			return nil, errors.WithMessage(err, errContext)
		}
		if path[0] != '/' {
			path = "/" + path
		}
		path = baseURL.String() + path
	}

	var req *http.Request
	if formValues != "" {
		req, err = http.NewRequest(method, path, strings.NewReader(formValues))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, path, inBody)
		if err != nil {
			return nil, err
		}
	}

	if mattermostUserID != "" {
		if err = c.plugin.AddAuthorization(req, mattermostUserID); err != nil {
			return nil, err
		}
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if mattermostUserID != "" {
		if err = c.plugin.AddAuthorization(req, mattermostUserID); err != nil {
			return nil, err
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if out != nil {
			if err = json.Unmarshal(responseData, out); err != nil {
				return responseData, err
			}
		}
		return responseData, nil

	case http.StatusNoContent:
		return nil, nil

	case http.StatusNotFound:
		return nil, ErrNotFound
	}

	errResp := ErrorResponse{}
	if err = json.Unmarshal(responseData, &errResp); err != nil {
		return responseData, errors.WithMessagef(err, "status: %s", resp.Status)
	}
	return responseData, fmt.Errorf("errorMessage %s", errResp.Message)
}

func InitClient(p *Plugin) Client {
	return &client{
		plugin:     p,
		httpClient: &http.Client{},
	}
}
