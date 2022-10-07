package plugin

import (
	"fmt"
	"strings"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-server/v5/model"
)

// postTaskPreview function returns the new post containing the preview of the work item.
// (UI may change in the future)
func (p *Plugin) PostTaskPreview(linkData []string, userID, channelID string) (*model.Post, string) {
	task, _, err := p.Client.GetTask(linkData[3], linkData[7], linkData[4], userID)
	if err != nil {
		return nil, ""
	}

	assignedTo := task.Fields.AssignedTo.DisplayName
	if assignedTo == "" {
		assignedTo = "None"
	}

	description := task.Fields.Description
	if description == "" {
		description = "No description"
	}

	post := &model.Post{
		UserId:    userID,
		ChannelId: channelID,
	}
	attachment := &model.SlackAttachment{
		Title: fmt.Sprintf(constants.TaskTitle, task.Fields.Type, task.ID, task.Fields.Title, task.Link.HTML.Href),
		Color: constants.BoardsIconColor,
		Fields: []*model.SlackAttachmentField{
			{
				Title: "State",
				Value: task.Fields.State,
				Short: true,
			},
			{
				Title: "Assigned To",
				Value: assignedTo,
				Short: true,
			},
			{
				Title: "Description",
				Value: description,
			},
		},
		Footer:     linkData[4],
		FooterIcon: fmt.Sprintf("%s/plugins/%s/static/%s", p.GetSiteURL(), constants.PluginID, constants.ProjectIcon),
	}
	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})

	return post, ""
}

func (p *Plugin) getReviewersListString(reviewersList []serializers.Reviewer) string {
	reviewers := ""
	for i := 0; i < len(reviewersList); i++ {
		if i != len(reviewersList)-1 {
			reviewers += fmt.Sprintf("%s, ", reviewersList[i].DisplayName)
		} else {
			reviewers += reviewersList[i].DisplayName
		}
	}

	if reviewers == "" {
		return "None" // When no reviewers are added
	}
	return reviewers
}

func (p *Plugin) PostPullRequestPreview(linkData []string, link, userID, channelID string) (*model.Post, string) {
	pullRequest, _, err := p.Client.GetPullRequest(linkData[3], linkData[8], linkData[6], userID)
	if err != nil {
		return nil, ""
	}

	var targetBranchName, sourceBranchName string
	if len(strings.Split(pullRequest.TargetRefName, "/")) == 3 {
		targetBranchName = strings.Split(pullRequest.TargetRefName, "/")[2]
	}

	if len(strings.Split(pullRequest.SourceRefName, "/")) == 3 {
		sourceBranchName = strings.Split(pullRequest.SourceRefName, "/")[2]
	}

	post := &model.Post{
		UserId:    userID,
		ChannelId: channelID,
	}
	reviewers := p.getReviewersListString(pullRequest.Reviewers)
	attachment := &model.SlackAttachment{
		Title: fmt.Sprintf(constants.PullRequestTitle, pullRequest.PullRequestID, pullRequest.Title, link),
		Color: constants.ReposIconColor,
		Fields: []*model.SlackAttachmentField{
			{
				Title: "Target Branch",
				Value: targetBranchName,
				Short: true,
			},
			{
				Title: "Source Branch",
				Value: sourceBranchName,
				Short: true,
			},
			{
				Title: "Reviewer(s)",
				Value: reviewers,
			},
		},
		Footer:     linkData[6],
		FooterIcon: fmt.Sprintf("%s/plugins/%s/static/%s", p.GetSiteURL(), constants.PluginID, constants.ProjectIcon),
	}
	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})

	return post, ""
}
