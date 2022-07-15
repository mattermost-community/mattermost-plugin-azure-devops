package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/pkg/errors"
)

type Client interface {
	TestApi() (string, error)
	// TODO: Remove later if not needed.
	// GetProjectList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.ProjectList, error)
	GetTaskList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.TaskList, error)
}

type client struct {
	plugin     *Plugin
	HTTPClient *http.Client
}

func (azureDevops *client) TestApi() (string, error) {

	return "hello world", nil
}

// TODO: Remove later if not needed.
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
func (azureDevops *client) GetTaskList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.TaskList, error) {
	page := queryParams[constants.Page].(int)

	// Query params of URL.
	params := url.Values{}
	params.Add(constants.PageQueryParam, fmt.Sprint(page*constants.TaskLimit))
	params.Add(constants.APIVersionQueryParam, constants.TasksIDAPIVersion)

	// Query to fetch the tasks IDs list.
	query := fmt.Sprintf(constants.TaskQuery, queryParams[constants.Project])

	// Add filters to the query.
	if queryParams[constants.Status] != "" {
		query += fmt.Sprintf(constants.TaskQueryStatusFilter, queryParams[constants.Status])
	}
	if queryParams[constants.AssignedTo] == "me" {
		query += constants.TaskQueryAssignedToFilter
	}

	// Query payload.
	taskQuery := map[string]string{
		"query": query,
	}
	// URL to fetch tasks IDs list.
	taskIDs := fmt.Sprintf(constants.GetTasksID, queryParams[constants.Organization])

	var taskIDList *serializers.TaskIDList
	if _, err := azureDevops.callJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskIDs, http.MethodPost, mattermostUserID, taskQuery, &taskIDList, params); err != nil {
		return nil, errors.Wrap(err, "failed to get Task ID list")
	}

	// Check if new task ID are present for current page.
	if page*constants.TaskLimit >= len(taskIDList.TaskList)+constants.TaskLimit {
		return nil, errors.Errorf(constants.NoResultPresent)
	}

	var IDs string
	for i := 0; i < len(taskIDList.TaskList); i++ {
		IDs += fmt.Sprint(strconv.Itoa(taskIDList.TaskList[i].ID), ",")
	}

	params = url.Values{}
	params.Add(constants.IDsQueryParam, strings.TrimSuffix(IDs, ","))
	params.Add(constants.APIVersionQueryParam, constants.TasksAPIVersion)

	// URL to fetch tasks list.
	task := fmt.Sprintf(constants.GetTasks, queryParams[constants.Organization])

	var taskList *serializers.TaskList
	if _, err := azureDevops.callJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, task, http.MethodGet, mattermostUserID, nil, &taskList, params); err != nil {
		return nil, errors.Wrap(err, "failed to get Task list")
	}

	return taskList, nil
}

// Wrapper to make REST API requests with "application/json" type content
func (azureDevops *client) callJSON(url, path, method, mattermostUserID string, in, out interface{}, params url.Values) (responseData []byte, err error) {
	contentType := "application/json"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(in)
	if err != nil {
		return nil, err
	}
	return azureDevops.call(url, method, path, contentType, mattermostUserID, buf, out, params)
}

// Makes HTTP request to REST APIs
func (azureDevops *client) call(basePath, method, path, contentType, mamattermostUserID string, inBody io.Reader, out interface{}, params url.Values) (responseData []byte, err error) {
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

	req, err := http.NewRequest(method, path, inBody)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if mamattermostUserID != "" {
		if err = azureDevops.plugin.AddAuthorization(req, mamattermostUserID); err != nil {
			return nil, err
		}
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	resp, err := azureDevops.HTTPClient.Do(req)
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
		return nil, errors.Errorf("not found")
	}

	type ErrorResponse struct {
		Message string `json:"message"`
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
		HTTPClient: &http.Client{},
	}
}
