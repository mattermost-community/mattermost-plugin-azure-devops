package store

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
)

type SubscriptionStore interface {
	StoreSubscription(subscription *serializers.SubscriptionDetails) error
	GetSubscriptionList() (*SubscriptionList, error)
	GetAllSubscriptions(userID string) ([]*serializers.SubscriptionDetails, error)
	DeleteSubscription(subscription *serializers.SubscriptionDetails) error
	StoreSubscriptionAndChannelIDMap(subscriptionID, webhookSecret, channelID string) error
	GetSubscriptionAndChannelIDMap(subscriptionID string) (*SubscriptionWebhookSecretAndChannelMap, error)
	DeleteSubscriptionAndChannelIDMap(subscriptionID string) error
}

type SubscriptionListMap map[string]serializers.SubscriptionDetails

type SubscriptionWebhookSecretAndChannelMap map[string]string

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

	subscriptionListValue := serializers.SubscriptionDetails{
		SubscriptionID:                   subscription.SubscriptionID,
		MattermostUserID:                 userID,
		ProjectName:                      subscription.ProjectName,
		ProjectID:                        subscription.ProjectID,
		OrganizationName:                 subscription.OrganizationName,
		ChannelID:                        subscription.ChannelID,
		EventType:                        subscription.EventType,
		ServiceType:                      subscription.ServiceType,
		ChannelName:                      subscription.ChannelName,
		ChannelType:                      subscription.ChannelType,
		CreatedBy:                        subscription.CreatedBy,
		CreatedAt:                        time.Now().UTC(),
		Repository:                       subscription.Repository,
		TargetBranch:                     subscription.TargetBranch,
		RepositoryName:                   subscription.RepositoryName,
		PullRequestCreatedBy:             subscription.PullRequestCreatedBy,
		PullRequestReviewersContains:     subscription.PullRequestReviewersContains,
		PullRequestCreatedByName:         subscription.PullRequestCreatedByName,
		PullRequestReviewersContainsName: subscription.PullRequestReviewersContainsName,
		PushedBy:                         subscription.PushedBy,
		PushedByName:                     subscription.PushedByName,
		MergeResult:                      subscription.MergeResult,
		MergeResultName:                  subscription.MergeResultName,
		NotificationType:                 subscription.NotificationType,
		NotificationTypeName:             subscription.NotificationTypeName,
		AreaPath:                         subscription.AreaPath,
		BuildStatus:                      subscription.BuildStatus,
		BuildPipeline:                    subscription.BuildPipeline,
		StageName:                        subscription.StageName,
		ReleasePipeline:                  subscription.ReleasePipeline,
		ReleaseStatus:                    subscription.ReleaseStatus,
		ApprovalType:                     subscription.ApprovalType,
		ApprovalStatus:                   subscription.ApprovalStatus,
		BuildStatusName:                  subscription.BuildStatusName,
		StageNameValue:                   subscription.StageNameValue,
		ReleasePipelineName:              subscription.ReleasePipelineName,
		ReleaseStatusName:                subscription.ReleaseStatusName,
		ApprovalTypeName:                 subscription.ApprovalTypeName,
		ApprovalStatusName:               subscription.ApprovalStatusName,
		RunPipeline:                      subscription.RunPipeline,
		RunPipelineName:                  subscription.RunPipelineName,
		RunStageName:                     subscription.RunStageName,
		RunEnvironmentName:               subscription.RunEnvironmentName,
		RunStageNameID:                   subscription.RunStageNameID,
		RunStageStateID:                  subscription.RunStageStateID,
		RunStageStateIDName:              subscription.RunStageStateIDName,
		RunStageResultID:                 subscription.RunStageResultID,
		RunStateID:                       subscription.RunStateID,
		RunStateIDName:                   subscription.RunStateIDName,
		RunResultID:                      subscription.RunResultID,
	}
	subscriptionList.ByMattermostUserID[userID][subscription.SubscriptionID] = subscriptionListValue
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

	subscriptionList.DeleteSubscriptionByKey(subscription.MattermostUserID, subscription.SubscriptionID)
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

func (s *Store) StoreSubscriptionAndChannelIDMap(subscriptionID, webhookSecret, channelID string) error {
	if err := s.StoreJSON(subscriptionID, SubscriptionWebhookSecretAndChannelMap{
		webhookSecret: channelID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSubscriptionAndChannelIDMap(subscriptionID string) (*SubscriptionWebhookSecretAndChannelMap, error) {
	var storedWebhookSecret SubscriptionWebhookSecretAndChannelMap
	if err := s.LoadJSON(subscriptionID, &storedWebhookSecret); err != nil {
		return nil, err
	}

	return &storedWebhookSecret, nil
}

func (s *Store) DeleteSubscriptionAndChannelIDMap(subscriptionID string) error {
	if err := s.Delete(subscriptionID); err != nil {
		return err
	}

	return nil
}
