package plugin

import (
	"fmt"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-server/v5/model"
)

// postTaskPreview function returns the new post containing the preview of the work item.
// (UI may change in the future)
func (p *Plugin) postTaskPreview(linkData []string, userID, channelID string) (*model.Post, string) {
	task, _, err := p.Client.GetTask(linkData[3], linkData[7], userID)
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

	taskTitle := fmt.Sprintf(constants.TaskTitle, task.Fields.Type, task.ID, task.Fields.Title, task.Link.Html.Href)
	TaskPreviewMessage := fmt.Sprintf(constants.TaskPreviewMessage, task.Fields.State, assignedTo, description)
	post := &model.Post{
		UserId:    userID,
		ChannelId: channelID,
	}
	attachment := &model.SlackAttachment{
		Text: fmt.Sprintf("%s\n%s\n", taskTitle, TaskPreviewMessage),
	}
	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})

	return post, ""
}
