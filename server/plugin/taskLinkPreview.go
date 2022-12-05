package plugin

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
)

// postTaskPreview function returns the new post containing the preview of the work item.
// (UI may change in the future)
func (p *Plugin) PostTaskPreview(linkData []string, userID, channelID string) (*model.Post, string) {
	task, _, err := p.Client.GetTask(linkData[3], linkData[7], linkData[4], userID)
	if err != nil {
		p.API.LogDebug("Error in getting task details from Azure", "Error", err.Error())
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
		AuthorName: "Azure Boards",
		AuthorIcon: fmt.Sprintf("%s/plugins/%s/static/%s", p.GetSiteURL(), constants.PluginID, constants.FileNameBoardsIcon),
		Title:      fmt.Sprintf(constants.TaskTitle, task.Fields.Type, task.ID, task.Fields.Title, task.Link.HTML.Href),
		Color:      constants.IconColorBoards,
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
		FooterIcon: fmt.Sprintf("%s/plugins/%s/static/%s", p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
	}
	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	return post, ""
}

func (p *Plugin) PostPullRequestPreview(linkData []string, link, userID, channelID string) (*model.Post, string) {
	pullRequest, _, err := p.Client.GetPullRequest(linkData[3], linkData[8], linkData[6], userID)
	if err != nil {
		p.API.LogDebug("Error in getting pull request details from Azure", "Error", err.Error())
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
		AuthorName: "Azure Repos",
		AuthorIcon: fmt.Sprintf("%s/plugins/%s/static/%s", p.GetSiteURL(), constants.PluginID, constants.FileNameReposIcon),
		Title:      fmt.Sprintf(constants.PullRequestTitle, pullRequest.PullRequestID, pullRequest.Title, link),
		Color:      constants.IconColorRepos,
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
		FooterIcon: fmt.Sprintf("%s/plugins/%s/static/%s", p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
	}
	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})

	return post, ""
}

func (p *Plugin) PostBuildDetailsPreview(linkData []string, link, userID, channelID string) (*model.Post, string) {
	organization := linkData[3]
	project := linkData[4]
	buildID := strings.Split(linkData[6], "&")[0][16:]
	buildDetails, _, err := p.Client.GetBuildDetails(organization, project, buildID, userID)
	if err != nil {
		p.API.LogDebug("Error in getting build details from Azure", "Error", err.Error())
		return nil, ""
	}

	post := &model.Post{
		UserId:    userID,
		ChannelId: channelID,
	}

	attachment := &model.SlackAttachment{
		AuthorName: "Azure Pipelines",
		AuthorIcon: fmt.Sprintf("%s/plugins/%s/static/%s", p.GetSiteURL(), constants.PluginID, constants.FileNamePipelinesIcon), // TODO: update icon file
		Title:      fmt.Sprintf(constants.BuildDetailsTitle, buildDetails.BuildNumber, buildDetails.Link.Web.Href, buildDetails.Definition.Name),
		Color:      constants.IconColorPipelines,
		Fields: []*model.SlackAttachmentField{
			{
				Title: "Repository",
				Value: buildDetails.Repository.Name,
				Short: true,
			},
			{
				Title: "Source Branch",
				Value: buildDetails.SourceBranch,
				Short: true,
			},
			{
				Title: "Requested By",
				Value: buildDetails.RequestedBy.DisplayName,
				Short: true,
			},
			{
				Title: "Status",
				Value: buildDetails.Status,
				Short: true,
			},
		},
		Footer:     project,
		FooterIcon: fmt.Sprintf("%s/plugins/%s/static/%s", p.GetSiteURL(), constants.PluginID, constants.FileNameProjectIcon),
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	return post, ""
}
