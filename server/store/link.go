package store

import (
	"encoding/json"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/pkg/errors"
)

type ProjectListMap map[string]Project

type Project struct {
	MattermostUserID string
	ProjectID        string
	ProjectName      string
	OrganizationName string
}

type ProjectList struct {
	ByMattermostUserID map[string]ProjectListMap
}

func NewProjectList() *ProjectList {
	return &ProjectList{
		ByMattermostUserID: map[string]ProjectListMap{},
	}
}

func (s *Store) StoreProject(project *Project) error {
	key := GetProjectListMapKey()
	if err := s.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
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
	}); err != nil {
		return err
	}

	return nil
}

func (projectList *ProjectList) AddProject(userID string, project *Project) {
	if _, valid := projectList.ByMattermostUserID[userID]; !valid {
		projectList.ByMattermostUserID[userID] = make(ProjectListMap)
	}
	projectKey := GetProjectKey(project.ProjectID, userID)
	projectListValue := Project{
		MattermostUserID: userID,
		ProjectID:        project.ProjectID,
		ProjectName:      project.ProjectName,
		OrganizationName: project.OrganizationName,
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

func (s *Store) GetAllProjects(userID string) ([]Project, error) {
	projects, err := s.GetProject()
	if err != nil {
		return nil, err
	}
	var projectList []Project
	for _, project := range projects.ByMattermostUserID[userID] {
		projectList = append(projectList, project)
	}
	return projectList, nil
}

func (s *Store) DeleteProject(project *Project) error {
	key := GetProjectListMapKey()
	if err := s.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
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


