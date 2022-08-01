package plugin

import (
	"fmt"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-server/v5/model"
)

// UI may change in the future.
// Function to return the new post of the work item.
func (p *Plugin) postTaskPreview(linkData []string, msg, userID, channelID string) (*model.Post, string) {
	taskData := serializers.GetTaskData{
		Organization: linkData[3],
		Project:      linkData[4],
		TaskID:       linkData[7],
	}
	task, _, err := p.Client.GetTask(taskData, userID)
	if err != nil {
		return nil, ""
	}

	assignedTo := task.Fields.AssignedTo.DisplayName
	if assignedTo == "" {
		assignedTo = "none"
	}

	description := task.Fields.Description
	if description == "" {
		description = "no description"
	}

	taskTitle := fmt.Sprintf(constants.TaskTitle, task.Fields.Type, task.ID, task.Fields.Title, task.Link.Html.Href)
	TaskPreviewMessage := fmt.Sprintf(constants.TaskPreviewMessage, task.Fields.State, assignedTo, description)
	message := fmt.Sprintf("%s\n%s\n```\n%s\n```", msg, taskTitle, TaskPreviewMessage)
	post := &model.Post{
		UserId:    userID,
		ChannelId: channelID,
		Message:   message,
	}

	return post, ""
}
