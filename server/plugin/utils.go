package plugin

import (
	"errors"

	"github.com/mattermost/mattermost-server/v5/model"
)

var ErrNotFound = errors.New("not found")

func (p *Plugin) sendEphemeralPost(args *model.CommandArgs, text string) (*model.CommandResponse, *model.AppError) {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)

	return &model.CommandResponse{}, nil
}
