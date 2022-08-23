package store

import (
	"encoding/json"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewProjectList(t *testing.T) {
	defer monkey.UnpatchAll()
	for _, testCase := range []struct {
		description string
	}{
		{
			description: "test NewProjectList",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			resp := NewProjectList()
			assert.NotNil(t, resp)
		})
	}
}

func TestStoreProjectAtomicModify(t *testing.T) {
	defer monkey.UnpatchAll()
	for _, testCase := range []struct {
		description         string
		projectList         *ProjectList
		marshalError        error
		projectListFromJSON error
	}{
		{
			description: "test StoreProjectAtomicModify when project is added successfully",
			projectList: NewProjectList(),
		},
		{
			description:  "test StoreProjectAtomicModify when marshaling gives error",
			projectList:  NewProjectList(),
			marshalError: errors.New("mockError"),
		},
		{
			description:         "test StoreProjectAtomicModify when ProjectListFromJSON gives error",
			projectList:         NewProjectList(),
			projectListFromJSON: errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.projectList.AddProject("mockMattermostUserId", &serializers.ProjectDetails{
				OrganizationName: "mockOrganization",
				ProjectID:        "mockProjectID",
				ProjectName:      "mockProject",
			})

			monkey.Patch(ProjectListFromJSON, func([]byte) (*ProjectList, error) {
				return testCase.projectList, testCase.projectListFromJSON
			})
			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})
			resp, err := storeProjectAtomicModify(&serializers.ProjectDetails{}, []byte{})

			if testCase.marshalError != nil || testCase.projectListFromJSON != nil {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, resp)
		})
	}
}

func TestStoreProject(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "test StoreProject when project is stored successfully",
		},
		{
			description: "test StoreProject when project is not stored successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetProjectListMapKey, func() string {
				return "mockProjectKey"
			})
			monkey.Patch(ProjectListFromJSON, func([]byte) (*ProjectList, error) {
				return &ProjectList{}, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "AtomicModify", func(*Store, string, func([]byte) ([]byte, error)) error {
				return testCase.err
			})

			err := s.StoreProject(&serializers.ProjectDetails{})

			if testCase.err != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestAddProject(t *testing.T) {
	defer monkey.UnpatchAll()
	for _, testCase := range []struct {
		description string
		projectList *ProjectList
	}{
		{
			description: "test AddProject when project is added successfully",
			projectList: NewProjectList(),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetProjectKey, func(string, string) string {
				return "mockMattermostUserID"
			})

			testCase.projectList.AddProject("mockMattermostUserId", &serializers.ProjectDetails{
				OrganizationName: "mockOrganization",
				ProjectID:        "mockProjectID",
				ProjectName:      "mockProject",
			})
		})
	}
}

func TestGetProjects(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description      string
		err              error
		projectListError error
	}{
		{
			description: "test GetProjects when projects are fetched successfully",
		},
		{
			description: "test GetProjects when load gives error",
			err:         errors.New("mockError"),
		},
		{
			description:      "test GetProjects when projects are not fetched successfully",
			projectListError: errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetProjectListMapKey, func() string {
				return "mockMattermostUserID"
			})
			monkey.Patch(ProjectListFromJSON, func([]byte) (*ProjectList, error) {
				return &ProjectList{}, testCase.projectListError
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "Load", func(*Store, string) ([]byte, error) {
				return []byte("mockState"), testCase.err
			})

			projectList, err := s.GetProject()

			if testCase.err != nil || testCase.projectListError != nil {
				assert.Nil(t, projectList)
				assert.NotNil(t, err)
				return
			}

			assert.NotNil(t, projectList)
			assert.Nil(t, err)
		})
	}
}

func TestGetAllProjects(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "test GetAllProjects when projects are fetched successfully",
		},
		{
			description: "test GetAllProjects when projects are not fetched successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "GetProject", func(*Store) (*ProjectList, error) {
				return &ProjectList{}, testCase.err
			})

			projectList, err := s.GetAllProjects("mockMattermostUserID")

			if testCase.err != nil {
				assert.Nil(t, projectList)
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestDeleteProjectAtomicModify(t *testing.T) {
	defer monkey.UnpatchAll()
	for _, testCase := range []struct {
		description         string
		projectList         *ProjectList
		marshalError        error
		projectListFromJSON error
	}{
		{
			description: "test DeleteProjectAtomicModify when project is added successfully",
			projectList: NewProjectList(),
		},
		{
			description:  "test DeleteProjectAtomicModify when marshaling gives error",
			projectList:  NewProjectList(),
			marshalError: errors.New("mockError"),
		},
		{
			description:         "test DeleteProjectAtomicModify when ProjectListFromJSON gives error",
			projectList:         NewProjectList(),
			projectListFromJSON: errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.projectList.AddProject("mockMattermostUserId", &serializers.ProjectDetails{
				OrganizationName: "mockOrganization",
				ProjectID:        "mockProjectID",
				ProjectName:      "mockProject",
			})

			monkey.Patch(GetProjectKey, func(string, string) string {
				return "mockProjectKey"
			})
			monkey.Patch(GetProjectListMapKey, func() string {
				return "mockMattermostUserID"
			})
			monkey.Patch(ProjectListFromJSON, func([]byte) (*ProjectList, error) {
				return testCase.projectList, testCase.projectListFromJSON
			})
			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})
			resp, err := deleteProjectAtomicModify(&serializers.ProjectDetails{}, []byte{})

			if testCase.marshalError != nil || testCase.projectListFromJSON != nil {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, resp)
		})
	}
}

func TestDeleteProject(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "test DeleteProject when project is deleted successfully",
		},
		{
			description: "test DeleteProject when project is not deleted successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetProjectListMapKey, func() string {
				return "mockProjectKey"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "AtomicModify", func(*Store, string, func([]byte) ([]byte, error)) error {
				return testCase.err
			})

			err := s.DeleteProject(&serializers.ProjectDetails{})

			if testCase.err != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestDeleteProjectByKey(t *testing.T) {
	defer monkey.UnpatchAll()
	for _, testCase := range []struct {
		description string
		projectList *ProjectList
	}{
		{
			description: "test DeleteProjectByKey when project is deleted successfully",
			projectList: NewProjectList(),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.projectList.AddProject("mockMattermostUserId", &serializers.ProjectDetails{
				ProjectID: "mockProjectID",
			})
			testCase.projectList.DeleteProjectByKey("mockMattermostUserId", "mockProjectID_mockMattermostUserId")
		})
	}
}

func TestProjectListFromJSON(t *testing.T) {
	defer monkey.UnpatchAll()
	for _, testCase := range []struct {
		description string
		bytes       []byte
		err         error
	}{
		{
			description: "test NewProjectList",
			bytes:       make([]byte, 0),
		},
		{
			description: "test NewProjectList when unmarshaling gives error",
			bytes:       make([]byte, 10),
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(json.Unmarshal, func([]byte, interface{}) error {
				return testCase.err
			})

			resp, err := ProjectListFromJSON(testCase.bytes)

			if testCase.err != nil {
				assert.Nil(t, resp)
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
