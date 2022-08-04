package serializers

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

type UserID struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
}

type PublisherInputs struct {
	ProjectID string `json:"projectId"`
}

type ConsumerInputs struct {
	URL string `json:"url"`
}

type SubscriptionValue struct {
	ID               string          `json:"id"`
	URL              string          `json:"url"`
	EventType        string          `json:"eventType"`
	ConsumerID       string          `json:"consumerId"`
	ConsumerActionID string          `json:"consumerActionId"`
	CreatedBy        UserID          `json:"createdBy"`
	ModifiedBy       UserID          `json:"modifiedBy"`
	PublisherInputs  PublisherInputs `json:"publisherInputs"`
	ConsumerInputs   ConsumerInputs  `json:"consumerInputs"`
}

type SubscriptionList struct {
	Count             int                 `json:"count"`
	SubscriptionValue []SubscriptionValue `json:"value"`
}

type SubscriptionListRequestPayload struct {
	Organization string `json:"organization"`
}

type CreateSubscriptionRequestPayload struct {
	Organization string `json:"organization"`
	Project      string `json:"project"`
	EventType    string `json:"eventType"`
	ChannelName  string `json:"channelName"`
}

type CreateSubscriptionBodyPayload struct {
	PublisherID      string          `json:"publisherId"`
	EventType        string          `json:"eventType"`
	ConsumerId       string          `json:"consumerId"`
	ConsumerActionId string          `json:"consumerActionId"`
	PublisherInputs  PublisherInputs `json:"publisherInputs"`
	ConsumerInputs   ConsumerInputs  `json:"consumerInputs"`
}

type SubscriptionDetails struct {
	MattermostUserID string `json:"mattermostUserID"`
	ProjectName      string `json:"projectName"`
	ProjectID        string `json:"projectId"`
	OrganizationName string `json:"organizationName"`
	EventType        string `json:"eventType"`
	ChannelID        string `json:"channelID"`
	SubscriptionID   string `json:"subscriptionID"`
}

type DetailedMessage struct {
	Markdown string `json:"markdown"`
}

type SubscriptionNotification struct {
	DetailedMessage DetailedMessage `json:"detailedMessage"`
}

func SubscriptionListRequestPayloadFromJSON(data io.Reader) (*SubscriptionListRequestPayload, error) {
	var body *SubscriptionListRequestPayload
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}

func CreateSubscriptionRequestPayloadFromJSON(data io.Reader) (*CreateSubscriptionRequestPayload, error) {
	var body *CreateSubscriptionRequestPayload
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}

func SubscriptionNotificationFromJSON(data io.Reader) (*SubscriptionNotification, error) {
	var body *SubscriptionNotification
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}

func (t *SubscriptionListRequestPayload) IsSubscriptionRequestPayloadValid() error {
	if t.Organization == "" {
		return errors.New(constants.OrganizationRequired)
	}
	return nil
}

func (t *CreateSubscriptionRequestPayload) IsSubscriptionRequestPayloadValid() error {
	if t.Organization == "" {
		return errors.New(constants.OrganizationRequired)
	}
	if t.Project == "" {
		return errors.New(constants.ProjectRequired)
	}
	if t.EventType == "" {
		return errors.New(constants.EventTypeRequired)
	}
	if t.ChannelName == "" {
		return errors.New(constants.ChannelNameRequired)
	}
	return nil
}
