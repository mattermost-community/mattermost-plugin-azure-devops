package store

import (
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/stretchr/testify/assert"
)

func TestStoreUser(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "StoreUser: user is stored successfully",
		},
		{
			description: "StoreUser: user is not stored successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetOAuthKey, func(string) string {
				return "mockMattermostUserID"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "StoreJSON", func(*Store, string, interface{}) error {
				return testCase.err
			})

			err := s.StoreUser(&serializers.User{})

			if testCase.err != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestLoadUser(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "LoadUser: user is loaded successfully",
		},
		{
			description: "LoadUser: user is not loaded successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetOAuthKey, func(string) string {
				return "mockMattermostUserID"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "LoadJSON", func(*Store, string, interface{}) error {
				return testCase.err
			})

			user, err := s.LoadUser("mockMattermostUserID")

			if testCase.err != nil {
				assert.NotNil(t, err)
				assert.Nil(t, user)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, user)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "DeleteUser: user is deleted successfully",
		},
		{
			description: "DeleteUser: user is not deleted successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetOAuthKey, func(string) string {
				return "mockMattermostUserID"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "Delete", func(*Store, string) error {
				return testCase.err
			})

			isDeleted, err := s.DeleteUser("mockMattermostUserID")

			if testCase.err != nil {
				assert.NotNil(t, err)
				assert.Equal(t, false, isDeleted)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, true, isDeleted)
		})
	}
}
