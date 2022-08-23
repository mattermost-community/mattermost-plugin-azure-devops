package store

import (
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestStoreOAuthState(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "test StoreOAuthState when oAuth state is stored successfully",
		},
		{
			description: "test StoreOAuthState when oAuth is not stored successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetOAuthKey, func(string) string {
				return "mockMattermostUserID"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "StoreTTL", func(*Store, string, []byte, int64) error {
				return testCase.err
			})

			err := s.StoreOAuthState("mockMattermostUserID", "mockState")

			if testCase.err != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestVerifyOAuthState(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "test VerifyOAuthState when oAuth is verified successfully",
		},
		{
			description: "test VerifyOAuthState when oAuth is not verified successfully",
			err:         errors.New("mockError"),
		},
		{
			description: "test VerifyOAuthState when oAuth is not found",
			err:         ErrNotFound,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetOAuthKey, func(string) string {
				return "mockMattermostUserID"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "Load", func(*Store, string) ([]byte, error) {
				return []byte("mockState"), testCase.err
			})

			err := s.VerifyOAuthState("mockMattermostUserID", "mockState")

			if testCase.err != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
