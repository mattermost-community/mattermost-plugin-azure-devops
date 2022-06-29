package main

import (
	mattermostPlugin "github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/plugin"
)

func main() {
	mattermostPlugin.ClientMain(&plugin.Plugin{})
}
