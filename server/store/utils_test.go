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

func TestEnsure(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		loadError   error
		storeError  error
	}{
		{
			description: "test Ensure",
		},
		{
			description: "test Ensure when load gives error",
			loadError:   errors.New("mockError"),
		},
		{
			description: "test Ensure when store gives error",
			loadError:   ErrNotFound,
			storeError:  errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "Load", func(*Store, string) ([]byte, error) {
				return []byte{}, testCase.loadError
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "Store", func(*Store, string, []byte) error {
				return testCase.storeError
			})

			resp, err := s.Ensure("mockKey", []byte("mockData"))

			if testCase.loadError != nil || testCase.storeError != nil {
				assert.Nil(t, resp)
				assert.NotNil(t, err)
				return
			}

			assert.NotNil(t, resp)
			assert.Nil(t, err)
		})
	}
}

func TestLoadJSON(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		data        []byte
		loadError   error
	}{
		{
			description: "test LoadJSON",
			data:        []byte{},
		},
		{
			description: "test LoadJSON when load gives error",
			loadError:   errors.New("mockError"),
		},
		{
			description: "test LoadJSON when data is nil",
			data:        nil,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "Load", func(*Store, string) ([]byte, error) {
				return testCase.data, testCase.loadError
			})
			monkey.Patch(json.Unmarshal, func([]byte, interface{}) error {
				return nil
			})

			user := &serializers.User{}
			err := s.LoadJSON("mockKey", user)

			if testCase.loadError != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestStoreJSON(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description  string
		storeError   error
		marshalError error
	}{
		{
			description: "test StoreJSON",
		},
		{
			description:  "test StoreJSON when marshaling gives error",
			marshalError: errors.New("mockError"),
		},
		{
			description: "test StoreJSON when store gives error",
			storeError:  errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "Store", func(*Store, string, []byte) error {
				return testCase.storeError
			})
			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})

			user := &serializers.User{}
			err := s.StoreJSON("mockKey", user)

			if testCase.marshalError != nil || testCase.storeError != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
