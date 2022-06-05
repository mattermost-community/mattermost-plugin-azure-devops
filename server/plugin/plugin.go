package plugin

import (
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
	router        *mux.Router
	// api           *api.Api
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

// InitAPI initializes the REST API
func (p *Plugin) InitAPI() *mux.Router {

	r := mux.NewRouter()
	r.Use(p.WithRecovery)

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/test", test).Methods(http.MethodGet)
	// p.api.InitRoutes(s)
	// 404 handler
	r.Handle("{anything:.*}", http.NotFoundHandler())
	return r
}

func (p *Plugin) WithRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				p.API.LogError("Recovered from a panic",
					"url", r.URL.String(),
					"error", x,
					"stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func test(w http.ResponseWriter, req *http.Request) {
	// if we've made it here, we're authorized.
	w.WriteHeader(http.StatusOK)
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
