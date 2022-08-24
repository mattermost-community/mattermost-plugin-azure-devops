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

func TestNewSubscriptionListt(t *testing.T) {
	defer monkey.UnpatchAll()
	for _, testCase := range []struct {
		description string
	}{
		{
			description: "test NewSubscriptionList",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			resp := NewSubscriptionList()
			assert.NotNil(t, resp)
		})
	}
}

func TestStoreSubscriptionAtomicModify(t *testing.T) {
	defer monkey.UnpatchAll()
	subscriptionList := NewSubscriptionList()
	subscriptionList.AddSubscription("mockMattermostUserId", &serializers.SubscriptionDetails{
		OrganizationName: "mockOrganization",
		ProjectID:        "mockProjectID",
		ProjectName:      "mockProject",
		EventType:        "mockEventType",
		ChannelID:        "mockChannelID",
		ChannelName:      "mockChannelName",
		SubscriptionID:   "mockSubscriptionID",
	})
	for _, testCase := range []struct {
		description              string
		subscriptionList         *SubscriptionList
		marshalError             error
		subscriptionListFromJSON error
	}{
		{
			description:      "test StoreSubscriptionAtomicModify when subscription is added successfully",
			subscriptionList: subscriptionList,
		},
		{
			description:      "test StoreSubscriptionAtomicModify when marshaling gives error",
			subscriptionList: subscriptionList,
			marshalError:     errors.New("mockError"),
		},
		{
			description:              "test StoreSubscriptionAtomicModify when SubscriptionListFromJSON gives error",
			subscriptionList:         subscriptionList,
			subscriptionListFromJSON: errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(SubscriptionListFromJSON, func([]byte) (*SubscriptionList, error) {
				return testCase.subscriptionList, testCase.subscriptionListFromJSON
			})
			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})
			resp, err := storeSubscriptionAtomicModify(&serializers.SubscriptionDetails{}, []byte{})

			if testCase.marshalError != nil || testCase.subscriptionListFromJSON != nil {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, resp)
		})
	}
}

func TestStoreSubscription(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "test StoreSubscription when subscription is stored successfully",
		},
		{
			description: "test StoreSubscription when subscription is not stored successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetSubscriptionListMapKey, func() string {
				return "mockSubscriptionKey"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "AtomicModify", func(*Store, string, func([]byte) ([]byte, error)) error {
				return testCase.err
			})

			err := s.StoreSubscription(&serializers.SubscriptionDetails{})

			if testCase.err != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestAddSubscription(t *testing.T) {
	defer monkey.UnpatchAll()
	subscriptionList := NewSubscriptionList()
	for _, testCase := range []struct {
		description      string
		subscriptionList *SubscriptionList
	}{
		{
			description:      "test AddSubscription when subscription is added successfully",
			subscriptionList: subscriptionList,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetProjectKey, func(string, string) string {
				return "mockMattermostUserID"
			})

			testCase.subscriptionList.AddSubscription("mockMattermostUserId", &serializers.SubscriptionDetails{
				OrganizationName: "mockOrganization",
				ProjectID:        "mockProjectID",
				ProjectName:      "mockProject",
				EventType:        "mockEventType",
				ChannelID:        "mockChannelID",
				ChannelName:      "mockChannelName",
				SubscriptionID:   "mockSubscriptionID",
			})
		})
	}
}

func TestGetSubscription(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description           string
		err                   error
		subscriptionListError error
	}{
		{
			description: "test GetSubscription when subscriptions are fetched successfully",
		},
		{
			description: "test GetSubscription when 'Load' gives error",
			err:         errors.New("mockError"),
		},
		{
			description:           "test GetSubscription when subscriptions are not fetched successfully",
			subscriptionListError: errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetSubscriptionListMapKey, func() string {
				return "mockMattermostUserID"
			})
			monkey.Patch(SubscriptionListFromJSON, func([]byte) (*SubscriptionList, error) {
				return &SubscriptionList{}, testCase.subscriptionListError
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "Load", func(*Store, string) ([]byte, error) {
				return []byte("mockState"), testCase.err
			})

			subscriptionList, err := s.GetSubscription()

			if testCase.err != nil || testCase.subscriptionListError != nil {
				assert.Nil(t, subscriptionList)
				assert.NotNil(t, err)
				return
			}

			assert.NotNil(t, subscriptionList)
			assert.Nil(t, err)
		})
	}
}

func TestGetAllSubscriptions(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "test GetAllSubscriptions when subscriptions are fetched successfully",
		},
		{
			description: "test GetAllSubscriptions when subscriptions are not fetched successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "GetSubscription", func(*Store) (*SubscriptionList, error) {
				return &SubscriptionList{}, testCase.err
			})

			subscriptionList, err := s.GetAllSubscriptions("mockMattermostUserID")

			if testCase.err != nil {
				assert.Nil(t, subscriptionList)
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestDeleteSubscriptionAtomicModify(t *testing.T) {
	defer monkey.UnpatchAll()
	subscriptionList := NewSubscriptionList()
	subscriptionList.AddSubscription("mockMattermostUserId", &serializers.SubscriptionDetails{
		OrganizationName: "mockOrganization",
		ProjectID:        "mockProjectID",
		ProjectName:      "mockProject",
		EventType:        "mockEventType",
		ChannelID:        "mockChannelID",
		ChannelName:      "mockChannelName",
		SubscriptionID:   "mockSubscriptionID",
	})
	for _, testCase := range []struct {
		description              string
		subscriptionList         *SubscriptionList
		marshalError             error
		subscriptionListFromJSON error
	}{
		{
			description:      "test DeleteSubscriptionAtomicModify when subscription is added successfully",
			subscriptionList: subscriptionList,
		},
		{
			description:      "test DeleteSubscriptionAtomicModify when marshaling gives error",
			subscriptionList: subscriptionList,
			marshalError:     errors.New("mockError"),
		},
		{
			description:              "test DeleteSubscriptionAtomicModify when SubscriptionListFromJSON gives error",
			subscriptionList:         subscriptionList,
			subscriptionListFromJSON: errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetSubscriptionKey, func(string, string, string, string) string {
				return "mockSubscriptionKey"
			})
			monkey.Patch(SubscriptionListFromJSON, func([]byte) (*SubscriptionList, error) {
				return testCase.subscriptionList, testCase.subscriptionListFromJSON
			})
			monkey.Patch(json.Marshal, func(interface{}) ([]byte, error) {
				return []byte{}, testCase.marshalError
			})
			resp, err := deleteSubscriptionAtomicModify(&serializers.SubscriptionDetails{}, []byte{})

			if testCase.marshalError != nil || testCase.subscriptionListFromJSON != nil {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, resp)
		})
	}
}

func TestDeleteSubscription(t *testing.T) {
	defer monkey.UnpatchAll()
	s := Store{}
	for _, testCase := range []struct {
		description string
		err         error
	}{
		{
			description: "test DeleteSubscription when subscription is deleted successfully",
		},
		{
			description: "test DeleteSubscription when subscription is not deleted successfully",
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(GetSubscriptionListMapKey, func() string {
				return "mockSubscriotioKey"
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&s), "AtomicModify", func(*Store, string, func([]byte) ([]byte, error)) error {
				return testCase.err
			})

			err := s.DeleteSubscription(&serializers.SubscriptionDetails{})

			if testCase.err != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestDeleteSubscriptionByKey(t *testing.T) {
	defer monkey.UnpatchAll()
	subscriptionList := NewSubscriptionList()
	subscriptionList.AddSubscription("mockMattermostUserID", &serializers.SubscriptionDetails{
		ProjectName: "mockProjectName",
		ChannelID:   "mockChannelID",
		EventType:   "mockEventType",
	})
	for _, testCase := range []struct {
		description      string
		subscriptionList *SubscriptionList
	}{
		{
			description:      "test DeleteSubscriptionByKey when subscription is deleted successfully",
			subscriptionList: subscriptionList,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.subscriptionList.DeleteSubscriptionByKey("mockMattermostUserID", "mockMattermostUserID_mockProjectName_mockChannelID_mockEventType")
		})
	}
}

func TestSubscriptionListFromJSON(t *testing.T) {
	defer monkey.UnpatchAll()
	for _, testCase := range []struct {
		description string
		bytes       []byte
		err         error
	}{
		{
			description: "test SubscriptionListFromJSON",
			bytes:       make([]byte, 0),
		},
		{
			description: "test SubscriptionListFromJSON when unmarshaling gives error",
			bytes:       make([]byte, 10),
			err:         errors.New("mockError"),
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			monkey.Patch(json.Unmarshal, func([]byte, interface{}) error {
				return testCase.err
			})

			resp, err := SubscriptionListFromJSON(testCase.bytes)

			if testCase.err != nil {
				assert.Nil(t, resp)
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
