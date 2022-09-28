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

type CreateSubscriptionRequestPayload struct {
	Organization string `json:"organization"`
	Project      string `json:"project"`
	EventType    string `json:"eventType"`
	ChannelID    string `json:"channelID"`
}

type CreateSubscriptionBodyPayload struct {
	PublisherID      string          `json:"publisherId"`
	EventType        string          `json:"eventType"`
	ConsumerID       string          `json:"consumerId"`
	ConsumerActionID string          `json:"consumerActionId"`
	PublisherInputs  PublisherInputs `json:"publisherInputs"`
	ConsumerInputs   ConsumerInputs  `json:"consumerInputs"`
}

type SubscriptionDetails struct {
	MattermostUserID string `json:"mattermostUserID"`
	ProjectName      string `json:"projectName"`
	ProjectID        string `json:"projectID"`
	OrganizationName string `json:"organizationName"`
	EventType        string `json:"eventType"`
	ChannelID        string `json:"channelID"`
	ChannelName      string `json:"channelName"`
	ChannelType      string `json:"channelType"`
	SubscriptionID   string `json:"subscriptionID"`
	CreatedBy        string `json:"createdBy"`
}

type DetailedMessage struct {
	Markdown string `json:"markdown"`
}

type SubscriptionNotification struct {
	DetailedMessage DetailedMessage `json:"detailedMessage"`
}

type DeleteSubscriptionRequestPayload struct {
	Organization string `json:"organization"`
	Project      string `json:"project"`
	EventType    string `json:"eventType"`
	ChannelID    string `json:"channelID"`
	MMUserID     string `json:"mmUserID"`
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

func DeleteSubscriptionRequestPayloadFromJSON(data io.Reader) (*DeleteSubscriptionRequestPayload, error) {
	var body *DeleteSubscriptionRequestPayload
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
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
	if t.ChannelID == "" {
		return errors.New(constants.ChannelIDRequired)
	}
	return nil
}

func (t *DeleteSubscriptionRequestPayload) IsSubscriptionRequestPayloadValid() error {
	if t.Organization == "" {
		return errors.New(constants.OrganizationRequired)
	}
	if t.Project == "" {
		return errors.New(constants.ProjectRequired)
	}
	if t.EventType == "" {
		return errors.New(constants.EventTypeRequired)
	}
	if t.ChannelID == "" {
		return errors.New(constants.ChannelIDRequired)
	}
	if t.MMUserID == "" {
		return errors.New(constants.MMUserIDRequired)
	}
	return nil
}
