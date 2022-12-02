package store

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
)

type LinkStore interface {
	StoreProject(project *serializers.ProjectDetails) error
	GetProject() (*ProjectList, error)
	GetAllProjects(userID string) ([]serializers.ProjectDetails, error)
	DeleteProject(project *serializers.ProjectDetails) error
}

type ProjectListMap map[string]serializers.ProjectDetails

type ProjectList struct {
	ByMattermostUserID map[string]ProjectListMap
}

func NewProjectList() *ProjectList {
	return &ProjectList{
		ByMattermostUserID: map[string]ProjectListMap{},
	}
}

func storeProjectAtomicModify(project *serializers.ProjectDetails, initialBytes []byte) ([]byte, error) {
	projectList, err := ProjectListFromJSON(initialBytes)
	if err != nil {
		return nil, err
	}
	projectList.AddProject(project.MattermostUserID, project)
	modifiedBytes, marshalErr := json.Marshal(projectList)
	if marshalErr != nil {
		return nil, marshalErr
	}
	return modifiedBytes, nil
}

func (s *Store) StoreProject(project *serializers.ProjectDetails) error {
	key := GetProjectListMapKey()
	if err := s.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
		return storeProjectAtomicModify(project, initialBytes)
	}); err != nil {
		return err
	}

	return nil
}

func (projectList *ProjectList) AddProject(userID string, project *serializers.ProjectDetails) {
	if _, valid := projectList.ByMattermostUserID[userID]; !valid {
		projectList.ByMattermostUserID[userID] = make(ProjectListMap)
	}
	projectKey := GetProjectKey(project.ProjectID, userID)
	projectListValue := serializers.ProjectDetails{
		MattermostUserID: userID,
		ProjectID:        project.ProjectID,
		ProjectName:      project.ProjectName,
		OrganizationName: project.OrganizationName,
		IsAdmin:          project.IsAdmin,
	}
	projectList.ByMattermostUserID[userID][projectKey] = projectListValue
}

func (s *Store) GetProject() (*ProjectList, error) {
	key := GetProjectListMapKey()
	initialBytes, appErr := s.Load(key)
	if appErr != nil {
		return nil, errors.New(constants.GetProjectListError)
	}
	projects, err := ProjectListFromJSON(initialBytes)
	if err != nil {
		return nil, errors.New(constants.GetProjectListError)
	}
	return projects, nil
}

func (s *Store) GetAllProjects(userID string) ([]serializers.ProjectDetails, error) {
	projects, err := s.GetProject()
	if err != nil {
		return nil, err
	}
	var projectList []serializers.ProjectDetails
	for _, project := range projects.ByMattermostUserID[userID] {
		projectList = append(projectList, project)
	}
	return projectList, nil
}

func deleteProjectAtomicModify(project *serializers.ProjectDetails, initialBytes []byte) ([]byte, error) {
	projectList, err := ProjectListFromJSON(initialBytes)
	if err != nil {
		return nil, err
	}
	projectKey := GetProjectKey(project.ProjectID, project.MattermostUserID)
	projectList.DeleteProjectByKey(project.MattermostUserID, projectKey)
	modifiedBytes, marshalErr := json.Marshal(projectList)
	if marshalErr != nil {
		return nil, marshalErr
	}
	return modifiedBytes, nil
}

func (s *Store) DeleteProject(project *serializers.ProjectDetails) error {
	key := GetProjectListMapKey()
	if err := s.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
		return deleteProjectAtomicModify(project, initialBytes)
	}); err != nil {
		return err
	}

	return nil
}

func (projectList *ProjectList) DeleteProjectByKey(userID, projectKey string) {
	for key := range projectList.ByMattermostUserID[userID] {
		if key == projectKey {
			delete(projectList.ByMattermostUserID[userID], key)
		}
	}
}

func ProjectListFromJSON(bytes []byte) (*ProjectList, error) {
	var projectList *ProjectList
	if len(bytes) != 0 {
		unmarshalErr := json.Unmarshal(bytes, &projectList)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
	} else {
		projectList = NewProjectList()
	}
	return projectList, nil
}
