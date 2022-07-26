package plugin

import (
	"fmt"
	"strings"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/mattermost/mattermost-server/v5/model"
)

// UI may change in the future.
// Function to return the new post of the work item.
func (p *Plugin) postTaskPreview(msg, userID, channelID string) (*model.Post, string) {
	link := strings.Split(msg, "/")
	data := serializers.GetTaskData{
		Organization: link[3],
		Project:      link[4],
		TaskID:       link[7],
	}
	task, err := p.Client.GetTask(data, userID)
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
