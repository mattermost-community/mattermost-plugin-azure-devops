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
	GetProjectsList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.ProjectsList, error)
	GetTasksList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.TasksList, error)
}

type client struct {
	plugin     *Plugin
	HTTPClient *http.Client
}

func (azureDevops *client) TestApi() (string, error) {

	return "hello world", nil
}

// Function to get the list of projects.
func (azureDevops *client) GetProjectsList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.ProjectsList, error) {
	var projectsList *serializers.ProjectsList
	page := queryParams["page"].(int)

	// Url to fetch projects list.
	project := fmt.Sprintf(constants.GetProject, queryParams["organization"], page * constants.MaxProjectsPerPage)
	_, err := azureDevops.callJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, project, http.MethodGet, mattermostUserID, "", &projectsList)

	// Check if new projects are present for current page.
	if page * constants.MaxProjectsPerPage >= projectsList.Count + constants.MaxProjectsPerPage {
		return nil, errors.Errorf(constants.NoResultPresent)
	}
	if err != nil {
		errors.Wrap(err, "failed to get Projects list")
		return nil, err
	}

	return projectsList, nil
}

// Function to get the list of tasks.
func (azureDevops *client) GetTasksList(queryParams map[string]interface{}, mattermostUserID string) (*serializers.TasksList, error) {
	var tasksIDList *serializers.TasksIDList
	var tasksList *serializers.TasksList
	page := queryParams["page"].(int)

	// Url to fetch tasks IDs list.
	taskIDs := fmt.Sprintf(constants.GetTasksID, queryParams["organization"], page * constants.MaxTasksPerPage)

	// Query to fetch the tasks IDs list.
	query := fmt.Sprintf(constants.TaskQuery, queryParams["project"])

	// Add filters to the query.
	if queryParams["status"] != "" {
		query += fmt.Sprintf(constants.TaskQueryStatusFilter, queryParams["status"])
	}
	if queryParams["assignedTo"] == "me" {
		query += constants.TaskQueryAssignedToFilter
	}

	// Query payload.
	taskQuery := map[string]string{
		"query": query,
	}
	_, err := azureDevops.callJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, taskIDs, http.MethodPost, mattermostUserID, taskQuery, &tasksIDList)

	if err != nil {
		errors.Wrap(err, "failed to get Tasks Id list")
		return nil, err
	}

	// Check if new task Id are present for current page.
	if page * constants.MaxTasksPerPage >= len(tasksIDList.TaskList) + constants.MaxTasksPerPage {
		return nil, errors.Errorf(constants.NoResultPresent)
	}

	Ids := ""
	for i := 0; i < len(tasksIDList.TaskList); i++ {
		Ids += strconv.Itoa(tasksIDList.TaskList[i].Id) + ","
	}

	// Url to fetch tasks list.
	task := fmt.Sprintf(constants.GetTasks, queryParams["organization"], strings.TrimSuffix(Ids, ","))
	_, err = azureDevops.callJSON(azureDevops.plugin.getConfiguration().AzureDevopsAPIBaseURL, task, http.MethodGet, mattermostUserID, "", &tasksList)

	if err != nil {
		errors.Wrap(err, "failed to get Tasks list")
		return nil, err
	}

	return tasksList, nil
}

// Wrapper to make REST API requests with "application/json" type content
func (azureDevops *client) callJSON(url, path, method, mattermostUserID string, in, out interface{}) (responseData []byte, err error) {
	contentType := "application/json"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(in)
	if err != nil {
		return nil, err
	}
	return azureDevops.call(url, method, path, contentType, mattermostUserID, buf, out)
}

// Makes HTTP request to REST APIs
func (azureDevops *client) call(basePath, method, path, contentType, mamattermostUserID string, inBody io.Reader, out interface{}) (responseData []byte, err error) {
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

	err = azureDevops.plugin.AddAuthorization(req, mamattermostUserID)
	if err != nil {
		return nil, err
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
