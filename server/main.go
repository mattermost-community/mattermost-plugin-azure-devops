package main

import (
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/plugin"
	mattermostPlugin "github.com/mattermost/mattermost-server/v5/plugin"
)

func main() {
	mattermostPlugin.ClientMain(&plugin.Plugin{})
}
