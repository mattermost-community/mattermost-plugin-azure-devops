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
			description: "Ensure: valid",
		},
		{
			description: "Ensure: load gives error",
			loadError:   errors.New("mockError"),
		},
		{
			description: "Ensure: store gives error",
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
			description: "LoadJSON: valid",
			data:        []byte{},
		},
		{
			description: "LoadJSON: load gives error",
			loadError:   errors.New("mockError"),
		},
		{
			description: "LoadJSON: data is nil",
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
			description: "StoreJSON: valid",
		},
		{
			description:  "StoreJSON: marshaling gives error",
			marshalError: errors.New("mockError"),
		},
		{
			description: "StoreJSON: store gives error",
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
