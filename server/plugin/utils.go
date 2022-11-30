package plugin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var ErrNotFound = errors.New("not found")

// sendEphemeralPostForCommand sends an ephermal message
func (p *Plugin) sendEphemeralPostForCommand(args *model.CommandArgs, text string) (*model.CommandResponse, *model.AppError) {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)

	return &model.CommandResponse{}, nil
}

// DM posts a simple Direct Message to the specified user
func (p *Plugin) DM(mattermostUserID, format string, isSlackAttachment bool, args ...interface{}) (string, error) {
	channel, err := p.API.GetDirectChannel(mattermostUserID, p.botUserID)
	if err != nil {
		p.API.LogError("Couldn't get bot's DM channel", "userID", mattermostUserID, "Error", err.Error())
		return "", err
	}
	var post *model.Post
	post = &model.Post{
		ChannelId: channel.Id,
		UserId:    p.botUserID,
		Message:   fmt.Sprintf(format, args...),
	}
	if isSlackAttachment {
		post = &model.Post{
			ChannelId: channel.Id,
			UserId:    p.botUserID,
		}
		attachment := &model.SlackAttachment{
			Text: fmt.Sprintf(format, args...),
		}
		model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	}
	sentPost, err := p.API.CreatePost(post)
	if err != nil {
		p.API.LogError("Error occurred while creating post", "error", err.Error())
		return "", err
	}
	return sentPost.Id, nil
}

// Encode encodes bytes into base64 string
func (p *Plugin) Encode(encrypted []byte) string {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(encrypted)))
	base64.URLEncoding.Encode(encoded, encrypted)
	return string(encoded)
}

// Decode decodes a base64 string into bytes
func (p *Plugin) Decode(encoded string) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(encoded)))
	noOfBytes, err := base64.URLEncoding.Decode(decoded, []byte(encoded))
	if err != nil {
		return nil, err
	}
	return decoded[:noOfBytes], nil
}

// Encrypt used for generating encrypted bytes
func (p *Plugin) Encrypt(plain, secret []byte) ([]byte, error) {
	if len(secret) == 0 {
		return plain, nil
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	sealed := aesgcm.Seal(nil, nonce, plain, nil)
	return append(nonce, sealed...), nil
}

// Decrypt used for generating decrypted bytes
func (p *Plugin) Decrypt(encrypted, secret []byte) ([]byte, error) {
	if len(secret) == 0 {
		return encrypted, nil
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, errors.New("token too short")
	}

	nonce, encrypted := encrypted[:nonceSize], encrypted[nonceSize:]
	plain, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}

func (p *Plugin) GetSiteURL() string {
	return p.getConfiguration().MattermostSiteURL
}

func (p *Plugin) GetPluginURLPath() string {
	return fmt.Sprintf("/plugins/%s/api/v1", constants.PluginID)
}

func (p *Plugin) GetPluginURL() string {
	return fmt.Sprintf("%s%s", strings.TrimRight(p.GetSiteURL(), "/"), p.GetPluginURLPath())
}

func (p *Plugin) ParseAuthToken(encoded string) (string, error) {
	decodedAccessToken, err := p.Decode(encoded)
	if err != nil {
		return "", err
	}
	decryptedAccessToken, err := p.Decrypt(decodedAccessToken, []byte(p.getConfiguration().EncryptionSecret))
	if err != nil {
		return "", err
	}
	return string(decryptedAccessToken), nil
}

// AddAuthorization function to add authorization to a request.
func (p *Plugin) AddAuthorization(r *http.Request, mattermostUserID string) error {
	user, err := p.Store.LoadUser(mattermostUserID)
	if err != nil {
		return err
	}
	token, err := p.ParseAuthToken(user.AccessToken)
	if err != nil {
		return err
	}

	r.Header.Add(constants.Authorization, fmt.Sprintf("%s %s", constants.Bearer, token))
	return nil
}

func (p *Plugin) IsProjectLinked(projectList []serializers.ProjectDetails, project serializers.ProjectDetails) (*serializers.ProjectDetails, bool) {
	for _, a := range projectList {
		if a.ProjectName == project.ProjectName && a.OrganizationName == project.OrganizationName {
			return &a, true
		}
	}
	return nil, false
}

func (p *Plugin) IsSubscriptionPresent(subscriptionList []*serializers.SubscriptionDetails, subscription *serializers.SubscriptionDetails) (*serializers.SubscriptionDetails, bool) {
	for _, a := range subscriptionList {
		if a.ProjectName == subscription.ProjectName && a.OrganizationName == subscription.OrganizationName && a.ChannelID == subscription.ChannelID && a.EventType == subscription.EventType && a.Repository == subscription.Repository && a.TargetBranch == subscription.TargetBranch && a.PullRequestCreatedBy == subscription.PullRequestCreatedBy && a.PullRequestReviewersContains == subscription.PullRequestReviewersContains && a.PushedBy == subscription.PushedBy && a.MergeResult == subscription.MergeResult && a.NotificationType == subscription.NotificationType && a.AreaPath == subscription.AreaPath {
			return a, true
		}
	}
	return nil, false
}

func (p *Plugin) IsAnyProjectLinked(mattermostUserID string) (bool, error) {
	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		return false, err
	}

	if len(projectList) == 0 {
		return false, nil
	}

	return true, nil
}

func (p *Plugin) getConnectAccountFirstMessage() string {
	return fmt.Sprintf(constants.ConnectAccountFirst, fmt.Sprintf(constants.ConnectAccount, p.GetPluginURLPath(), constants.PathOAuthConnect))
}

func (p *Plugin) ParseSubscriptionsToCommandResponse(subscriptionsList []*serializers.SubscriptionDetails, channelID, createdBy, userID, command, teamID string) string {
	var sb strings.Builder

	filteredSubscriptionList, filteredSubscriptionErr := p.GetSubscriptionsForAccessibleChannelsOrProjects(subscriptionsList, teamID, userID)
	if filteredSubscriptionErr != nil {
		p.API.LogError(constants.FetchFilteredSubscriptionListError, "Error", filteredSubscriptionErr.Error())
		sb.WriteString(constants.GenericErrorMessage)
		return sb.String()
	}

	if len(filteredSubscriptionList) == 0 {
		sb.WriteString(fmt.Sprintf("No %s subscription exists", command))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("###### %s subscription(s)\n", cases.Title(language.Und).String(command)))
	sb.WriteString("| Subscription ID | Organization | Project | Event Type | Created By | Channel |\n")
	sb.WriteString("| :-------------- | :----------- | :------ | :--------- | :--------- | :------ |\n")

	displayEventType := map[string]string{
		constants.SubscriptionEventWorkItemCreated:                    "Work Item Created",
		constants.SubscriptionEventWorkItemUpdated:                    "Work Item Updated",
		constants.SubscriptionEventWorkItemDeleted:                    "Work Item Deleted",
		constants.SubscriptionEventWorkItemCommented:                  "Work Item Commented",
		constants.SubscriptionEventPullRequestCreated:                 "Pull Request Created",
		constants.SubscriptionEventPullRequestUpdated:                 "Pull Request Updated",
		constants.SubscriptionEventPullRequestMerged:                  "Pull Request Merge Attempted",
		constants.SubscriptionEventPullRequestCommented:               "Pull Requested Commented",
		constants.SubscriptionEventCodePushed:                         "Code Pushed",
		constants.SubscriptionEventBuildCompleted:                     "Build Completed",
		constants.SubscriptionEventReleaseAbandoned:                   "Release Abandoned",
		constants.SubscriptionEventReleaseCreated:                     "Release Created",
		constants.SubscriptionEventReleaseDeploymentApprovalCompleted: "Release Deployment Approval Completed",
		constants.SubscriptionEventReleaseDeploymentCompleted:         "Release Deployment Completed",
		constants.SubscriptionEventReleaseDeploymentEventPending:      "Release Deployment Event Pending",
		constants.SubscriptionEventReleaseDeploymentStarted:           "Release Deployment Started",
		constants.SubscriptionEventRunStageApprovalCompleted:          "Run Stage Approval Completed",
		constants.SubscriptionEventRunStageStateChanged:               "Run Stage State Changed",
		constants.SubscriptionEventRunStageWaitingForApproval:         "Run Stage Waiting For Approval",
		constants.SubscriptionEventRunStateChanged:                    "Run State Changed",
	}

	noSubscriptionFound := true
	for _, subscription := range filteredSubscriptionList {
		if channelID == "" || subscription.ChannelID == channelID {
			switch createdBy {
			case constants.FilterCreatedByMe:
				if subscription.MattermostUserID == userID && subscription.ServiceType == command {
					noSubscriptionFound = false
					sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n", subscription.SubscriptionID, subscription.OrganizationName, subscription.ProjectName, displayEventType[subscription.EventType], subscription.CreatedBy, subscription.ChannelName))
				}
			case constants.FilterCreatedByAnyone:
				if subscription.ServiceType == command {
					noSubscriptionFound = false
					sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n", subscription.SubscriptionID, subscription.OrganizationName, subscription.ProjectName, displayEventType[subscription.EventType], subscription.CreatedBy, subscription.ChannelName))
				}
			}
		}
	}

	if noSubscriptionFound {
		sb.Reset()
		sb.WriteString(fmt.Sprintf("No %s subscription exists", command))
	}

	return sb.String()
}

func (p *Plugin) GetOffsetAndLimitFromQueryParams(r *http.Request) (offset, limit int) {
	query := r.URL.Query()
	var page int
	if val, err := strconv.Atoi(query.Get(constants.QueryParamPage)); err != nil || val < 0 {
		p.API.LogError(constants.InvalidPaginationQueryParam, "Error", err.Error())
		page = constants.DefaultPage
	} else {
		page = val
	}

	val, err := strconv.Atoi(query.Get(constants.QueryParamPerPage))
	switch {
	case err != nil || val < 0 || val > constants.DefaultPerPageLimit: // We can keep max limit per page and default limit per page same
		p.API.LogError(constants.InvalidPaginationQueryParam, "Error", err.Error())
		limit = constants.DefaultPerPageLimit
	default:
		limit = val
	}

	return page * limit, limit
}

func (p *Plugin) GetSubscriptionsForAccessibleChannelsOrProjects(subscriptionList []*serializers.SubscriptionDetails, teamID, mattermostUserID string) ([]*serializers.SubscriptionDetails, error) {
	channels, channelErr := p.API.GetChannelsForTeamForUser(teamID, mattermostUserID, false)
	if channelErr != nil {
		p.API.LogError(constants.GetChannelError, "Error", channelErr.Error())
		return nil, channelErr
	}

	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		return nil, err
	}

	var filteredSubscriptionList []*serializers.SubscriptionDetails
	for _, subscription := range subscriptionList {
		// For each subscription on a project check if a user is an admin or member of the MM channel to return subscriptions
		if projectDetails, isProjectLinked := p.IsProjectLinked(projectList, serializers.ProjectDetails{OrganizationName: subscription.OrganizationName, ProjectName: subscription.ProjectName}); isProjectLinked {
			if projectDetails.IsAdmin {
				filteredSubscriptionList = append(filteredSubscriptionList, subscription)
			} else {
				for _, channel := range channels {
					if subscription.ChannelID == channel.Id {
						filteredSubscriptionList = append(filteredSubscriptionList, subscription)
						break
					}
				}
			}
		}
	}

	return filteredSubscriptionList, nil
}

// TODO: use this function at all the places where baseURL need to be updated this way
func (p *Plugin) updateBaseURLForReleaseEventTypes(url, eventType string) string {
	if strings.Contains(eventType, "release") {
		url = strings.Replace(url, "://", "://vsrm.", 1)
	}

	return url
}
