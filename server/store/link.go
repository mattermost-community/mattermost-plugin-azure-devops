package store

import (
	"fmt"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/pkg/errors"
)

type ProjectList struct {
	Project []Project
}

type Project struct {
	MattermostUserID string
	ProjectID        string
	ProjectName      string
	OrganizationName string
}

func (s *Store) StoreProject(project *Project) error {
	projectKey := fmt.Sprintf(constants.ProjectListPrefix, project.MattermostUserID)
	prevProject, err := s.LoadProject(project.MattermostUserID)
	if err != nil {
		return err
	}

	// Check if a project is already linked with a user.
	for _, value := range prevProject.Project {
		if value.ProjectName == project.ProjectName {
			return nil
		}
	}

	prevProject.Project = append(prevProject.Project, *project)
	if err := s.StoreJSON(projectKey, prevProject); err != nil {
		return err
	}

	return nil
}

func (s *Store) LoadProject(mattermostUserID string) (*ProjectList, error) {
	projectKey := fmt.Sprintf(constants.ProjectListPrefix, mattermostUserID)
	project := ProjectList{}
	if err := s.LoadJSON(projectKey, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *Store) DeleteProject(project *Project) bool {
	projectKey := fmt.Sprintf(constants.ProjectListPrefix, project.MattermostUserID)
	projectList, err := s.LoadProject(project.MattermostUserID)
	if err != nil {
		return false
	}
	newProjectList := ProjectList{}

	for _, value := range projectList.Project {
		if value.ProjectName != project.ProjectName {
			newProjectList.Project = append(newProjectList.Project, value)
		}
	}
	if err := s.Delete(projectKey); err != nil {
		errors.Wrap(err, err.Error())
		return false
	}
	if err := s.StoreJSON(projectKey, newProjectList); err != nil {
		return false
	}
	return true
}
