package store

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
)

type SubscriptionStore interface {
	StoreSubscription(subscription *serializers.SubscriptionDetails) error
	GetSubscriptionList() (*SubscriptionList, error)
	GetAllSubscriptions(userID string) ([]*serializers.SubscriptionDetails, error)
	DeleteSubscription(subscription *serializers.SubscriptionDetails) error
}

type SubscriptionListMap map[string]serializers.SubscriptionDetails

type SubscriptionList struct {
	ByMattermostUserID map[string]SubscriptionListMap
}

func NewSubscriptionList() *SubscriptionList {
	return &SubscriptionList{
		ByMattermostUserID: map[string]SubscriptionListMap{},
	}
}

func storeSubscriptionAtomicModify(subscription *serializers.SubscriptionDetails, initialBytes []byte) ([]byte, error) {
	subscriptionList, err := SubscriptionListFromJSON(initialBytes)
	if err != nil {
		return nil, err
	}

	subscriptionList.AddSubscription(subscription.MattermostUserID, subscription)
	modifiedBytes, marshalErr := json.Marshal(subscriptionList)
	if marshalErr != nil {
		return nil, marshalErr
	}
	return modifiedBytes, nil
}

func (s *Store) StoreSubscription(subscription *serializers.SubscriptionDetails) error {
	key := GetSubscriptionListMapKey()
	if err := s.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
		return storeSubscriptionAtomicModify(subscription, initialBytes)
	}); err != nil {
		return err
	}

	return nil
}

func (subscriptionList *SubscriptionList) AddSubscription(userID string, subscription *serializers.SubscriptionDetails) {
	if _, valid := subscriptionList.ByMattermostUserID[userID]; !valid {
		subscriptionList.ByMattermostUserID[userID] = make(SubscriptionListMap)
	}

	subscriptionKey := GetSubscriptionKey(userID, subscription.ProjectName, subscription.ChannelID, subscription.EventType)
	subscriptionListValue := serializers.SubscriptionDetails{
		MattermostUserID: userID,
		ProjectName:      subscription.ProjectName,
		ProjectID:        subscription.ProjectID,
		OrganizationName: subscription.OrganizationName,
		ChannelID:        subscription.ChannelID,
		EventType:        subscription.EventType,
		SubscriptionID:   subscription.SubscriptionID,
		ChannelName:      subscription.ChannelName,
		ChannelType:      subscription.ChannelType,
		CreatedBy:        subscription.CreatedBy,
	}
	subscriptionList.ByMattermostUserID[userID][subscriptionKey] = subscriptionListValue
}

func (s *Store) GetSubscriptionList() (*SubscriptionList, error) {
	key := GetSubscriptionListMapKey()
	initialBytes, appErr := s.Load(key)
	if appErr != nil {
		return nil, errors.New(constants.GetSubscriptionListError)
	}

	subscriptions, err := SubscriptionListFromJSON(initialBytes)
	if err != nil {
		return nil, errors.New(constants.GetSubscriptionListError)
	}

	return subscriptions, nil
}

func (s *Store) GetAllSubscriptions(userID string) ([]*serializers.SubscriptionDetails, error) {
	subscriptions, err := s.GetSubscriptionList()
	if err != nil {
		return nil, err
	}

	var subscriptionList []*serializers.SubscriptionDetails
	if userID == "" {
		for mmUserID := range subscriptions.ByMattermostUserID {
			for _, subscription := range subscriptions.ByMattermostUserID[mmUserID] {
				subscription := subscription // we need to do this to prevent implicit memory aliasing in for loop
				subscriptionList = append(subscriptionList, &subscription)
			}
		}
	} else {
		for _, subscription := range subscriptions.ByMattermostUserID[userID] {
			subscription := subscription
			subscriptionList = append(subscriptionList, &subscription)
		}
	}

	return subscriptionList, nil
}

func deleteSubscriptionAtomicModify(subscription *serializers.SubscriptionDetails, initialBytes []byte) ([]byte, error) {
	subscriptionList, err := SubscriptionListFromJSON(initialBytes)
	if err != nil {
		return nil, err
	}
	subscriptionKey := GetSubscriptionKey(subscription.MattermostUserID, subscription.ProjectName, subscription.ChannelID, subscription.EventType)
	subscriptionList.DeleteSubscriptionByKey(subscription.MattermostUserID, subscriptionKey)
	modifiedBytes, marshalErr := json.Marshal(subscriptionList)
	if marshalErr != nil {
		return nil, marshalErr
	}
	return modifiedBytes, nil
}

func (s *Store) DeleteSubscription(subscription *serializers.SubscriptionDetails) error {
	key := GetSubscriptionListMapKey()
	if err := s.AtomicModify(key, func(initialBytes []byte) ([]byte, error) {
		return deleteSubscriptionAtomicModify(subscription, initialBytes)
	}); err != nil {
		return err
	}

	return nil
}

func (subscriptionList *SubscriptionList) DeleteSubscriptionByKey(userID, subscriptionKey string) {
	for key := range subscriptionList.ByMattermostUserID[userID] {
		if key == subscriptionKey {
			delete(subscriptionList.ByMattermostUserID[userID], key)
		}
	}
}

func SubscriptionListFromJSON(bytes []byte) (*SubscriptionList, error) {
	var subscriptionList *SubscriptionList
	if len(bytes) != 0 {
		unmarshalErr := json.Unmarshal(bytes, &subscriptionList)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
	} else {
		subscriptionList = NewSubscriptionList()
	}
	return subscriptionList, nil
}
