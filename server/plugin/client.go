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
	GenerateOAuthToken(encodedFormValues url.Values) (*serializers.OAuthSuccessResponse, error)
	// TODO: WIP.
	// GetProjectList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.ProjectList, error)
	// GetTaskList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.TaskList, error)
	CreateTask(body *serializers.TaskCreateRequestPayload, mattermostUserID string) (*serializers.TaskValue, error)
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

func (c *client) GenerateOAuthToken(encodedFormValues url.Values) (*serializers.OAuthSuccessResponse, error) {
	var oAuthSuccessResponse *serializers.OAuthSuccessResponse

	_, err := c.callFormURLEncoded(constants.BaseOauthURL, constants.PathToken, http.MethodPost, nil, &oAuthSuccessResponse, encodedFormValues)
	if err != nil {
		return nil, err
	}

	return oAuthSuccessResponse, nil
}

// TODO: WIP.
// Function to get the list of projects.
// func (azureDevops *client) GetProjectList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.ProjectList, error) {
// 	var projectList *serializers.ProjectList
// 	page := queryParams["page"].(int)

// 	// Query params of URL.
// 	params := url.Values{}
// 	params.Add(constants.PageQueryParam, fmt.Sprint(page*constants.ProjectLimit))
// 	params.Add(constants.APIVersionQueryParam, constants.ProjectAPIVersion)

// 	// URL to fetch projects list.
// 	project := fmt.Sprintf(constants.GetProjects, queryParams["organization"])
// 	if _, err := azureDevops.callJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, project, http.MethodGet, mattermostUserID, nil, &projectList, params); err != nil {
// 		return nil, errors.Wrap(err, "failed to get Projects list")
// 	}

// 	// Check if new projects are present for current page.
// 	if page*constants.ProjectLimit >= projectList.Count+constants.ProjectLimit {
// 		return nil, errors.Errorf(constants.NoResultPresent)
// 	}
// 	return projectList, nil
// }

// Function to get the list of tasks.
// func (azureDevops *client) GetTaskList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.TaskList, error) {
// 	contentType := "application/json"
// 	page := queryParams[constants.Page].(int)

// 	// Query params of URL.
// 	params := url.Values{}
// 	params.Add(constants.PageQueryParam, fmt.Sprint(page*constants.TaskLimit))
// 	params.Add(constants.APIVersionQueryParam, constants.TasksIDAPIVersion)

// 	// Query to fetch the tasks IDs list.
// 	query := fmt.Sprintf(constants.TaskQuery, queryParams[constants.Project])

// 	// Add filters to the query.
// 	if queryParams[constants.Status] != "" {
// 		query += fmt.Sprintf(constants.TaskQueryStatusFilter, queryParams[constants.Status])
// 	}
// 	if queryParams[constants.AssignedTo] == "me" {
// 		query += constants.TaskQueryAssignedToFilter
// 	}

// 	// Query payload.
// 	taskQuery := map[string]string{
// 		"query": query,
// 	}
// 	// URL to fetch tasks IDs list.
// 	taskIDs := fmt.Sprintf(constants.GetTasksID, queryParams[constants.Organization])

// 	var taskIDList *serializers.TaskIDList
// 	if _, err := azureDevops.callJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskIDs, http.MethodPost, mattermostUserID, taskQuery, &taskIDList, params, contentType); err != nil {
// 		return nil, errors.Wrap(err, "failed to get Task ID list")
// 	}

// 	// Check if new task ID are present for current page.
// 	if page*constants.TaskLimit >= len(taskIDList.TaskList)+constants.TaskLimit {
// 		return nil, errors.Errorf(constants.NoResultPresent)
// 	}

// 	var IDs string
// 	for i := 0; i < len(taskIDList.TaskList); i++ {
// 		IDs += fmt.Sprint(strconv.Itoa(taskIDList.TaskList[i].ID), ",")
// 	}

// 	params = url.Values{}
// 	params.Add(constants.IDsQueryParam, strings.TrimSuffix(IDs, ","))
// 	params.Add(constants.APIVersionQueryParam, constants.TasksAPIVersion)

// 	// URL to fetch tasks list.
// 	task := fmt.Sprintf(constants.GetTasks, queryParams[constants.Organization])

// 	var taskList *serializers.TaskList
// 	if _, err := azureDevops.callJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, task, http.MethodGet, mattermostUserID, nil, &taskList, params, contentType); err != nil {
// 		return nil, errors.Wrap(err, "failed to get Task list")
// 	}

// 	return taskList, nil
// }

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
	if _, err := azureDevops.callPatchJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskURL, http.MethodPost, mattermostUserID, payload, &task, nil); err != nil {
		return nil, errors.Wrap(err, "failed to get create Task")
	}

	return task, nil
}

// Function to link a project and an organization.
func (c *client) Link(body *serializers.LinkRequestPayload, mattermostUserID string) (*serializers.Project, error) {
	projectURL := fmt.Sprintf(constants.GetProject, body.Organization, body.Project)
	var project *serializers.Project

	if _, err := c.callJSON(c.plugin.getConfiguration().AzureDevopsAPIBaseURL, projectURL, http.MethodGet, mattermostUserID, nil, &project, nil); err != nil {
		return nil, errors.Wrap(err, "failed to link Project")
	}

	return project, nil
}

// Wrapper to make REST API requests with "application/json-patch+json" type content
func (c *client) callPatchJSON(url, path, method, mattermostUserID string, in, out interface{}, formValues url.Values) (responseData []byte, err error) {
	contentType := "application/json-patch+json"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(in)
	if err != nil {
		return nil, err
	}
	return c.call(url, method, path, contentType, mattermostUserID, buf, out, formValues)
}

// Wrapper to make REST API requests with "application/json" type content
func (c *client) callJSON(url, path, method, mattermostUserID string, in, out interface{}, formValues url.Values) (responseData []byte, err error) {
	contentType := "application/json"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(in)
	if err != nil {
		return nil, err
	}
	return c.call(url, method, path, contentType, mattermostUserID, buf, out, formValues)
}

// Wrapper to make REST API requests with "application/x-www-form-urlencoded" type content
func (c *client) callFormURLEncoded(url, path, method string, in, out interface{}, formValues url.Values) (responseData []byte, err error) {
	contentType := "application/x-www-form-urlencoded"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(in)
	if err != nil {
		return nil, err
	}
	return c.call(url, method, path, contentType, "", buf, out, formValues)
}

// Makes HTTP request to REST APIs
func (c *client) call(basePath, method, path, contentType string, mamattermostUserID string, inBody io.Reader, out interface{}, formValues url.Values) (responseData []byte, err error) {
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
	if formValues != nil {
		req, err = http.NewRequest(method, path, strings.NewReader(formValues.Encode()))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, path, inBody)
		if err != nil {
			return nil, err
		}
	}

	if mamattermostUserID != "" {
		if err = c.plugin.AddAuthorization(req, mamattermostUserID); err != nil {
			return nil, err
		}
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
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
			err = json.Unmarshal(responseData, out)
			if err != nil {
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
	err = json.Unmarshal(responseData, &errResp)
	if err != nil {
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
