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

func SubscriptionListRequestPayloadFromJSON(data io.Reader) (*SubscriptionListRequestPayload, error) {
	var body *SubscriptionListRequestPayload
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
