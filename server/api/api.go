package api

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Add custom routes and corresponding handlers here
func InitRoutes(muxRouter *mux.Router) {
}

// Handles the static files under the assets directory.
func HandleStaticFiles(pluginApi plugin.API, r *mux.Router) {
	bundlePath, err := pluginApi.GetBundlePath()
	if err != nil {
		pluginApi.LogWarn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	// This will serve static files from the 'assets' directory under '/static/<filename>'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(bundlePath, "assets")))))
}
