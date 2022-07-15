package plugin

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/gorilla/mux"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

// Initializes the plugin REST API
func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.Use(p.WithRecovery)

	// 404 handler
	r.Handle(constants.WildRoute, http.NotFoundHandler())
	return r
}

// Add custom routes and corresponding handlers here
func (p *Plugin) InitRoutes() {
	p.Client = InitClient(p)

	s := p.router.PathPrefix(constants.APIPrefix).Subrouter()

	// OAuth
	s.HandleFunc(constants.PathOAuthConnect, p.OAuthConnect).Methods(http.MethodGet)
	s.HandleFunc(constants.PathOAuthCallback, p.OAuthComplete).Methods(http.MethodGet)
	// s.HandleFunc("/projects", p.handleGetProjects).Methods(http.MethodGet)
	s.HandleFunc("/tasks", p.handleGetTasks).Methods(http.MethodGet)

	// TODO: for testing purpose, remove later
	s.HandleFunc("/test", p.testAPI).Methods(http.MethodGet)
}

// Todo later.
// API to get projects in an organization.
// func (p *Plugin) handleGetProjects(w http.ResponseWriter, r *http.Request) {
// 	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
// 	if mattermostUserID == "" {
// 		http.Error(w, constants.NotAuthorized, http.StatusUnauthorized)
// 		return
// 	}

// 	organization := r.URL.Query().Get("organization")
// 	if organization == "" {
// 		http.Error(w, constants.OrganizationRequired, http.StatusBadRequest)
// 		return
// 	}

// 	page := StringToInt(r.URL.Query().Get("page"))
// 	if page <= 0 {
// 		http.Error(w, constants.InvalidPageNumber, http.StatusBadRequest)
// 		return
// 	}

// 	// Wrap all query params.
// 	queryParams := map[string]interface{}{
// 		"organization": organization,
// 		"page":         page,
// 	}

// 	boards, err := p.Client.GetProjectList(queryParams, mattermostUserID)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		if _, err := w.Write([]byte(err.Error())); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	response, err := json.Marshal(boards)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		if _, err := w.Write([]byte(err.Error())); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		return
// 	}
// 	w.Header().Add("Content-Type", "application/json")
// 	if _, err := w.Write(response); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// API to get tasks of a projects in an organization.
func (p *Plugin) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	if mattermostUserID == "" {
		http.Error(w, constants.NotAuthorized, http.StatusUnauthorized)
		return
	}
	statusData := map[string]string{
		"doing": "doing",
		"to-do": "To Do",
		"done":  "done",
	}
	organization := r.URL.Query().Get("organization")
	if organization == "" {
		http.Error(w, constants.OrganizationRequired, http.StatusBadRequest)
		return
	}
	project := r.URL.Query().Get("project")
	if project == "" {
		http.Error(w, constants.ProjectRequired, http.StatusBadRequest)
		return
	}
	status := r.URL.Query().Get("status")
	if status != "" && statusData[status] == "" {
		http.Error(w, constants.InvalidStatus, http.StatusBadRequest)
		return
	}
	assignedTo := r.URL.Query().Get("assigned_to")
	if assignedTo != "" && assignedTo != "me" {
		http.Error(w, constants.InvalidAssignedTo, http.StatusBadRequest)
		return
	}
	page := StringToInt(r.URL.Query().Get("page"))
	if page <= 0 {
		http.Error(w, constants.InvalidPageNumber, http.StatusBadRequest)
		return
	}

	// Wrap all query params.
	queryParams := map[string]interface{}{
		"organization": organization,
		"project":      project,
		"status":       statusData[status],
		"assignedTo":   assignedTo,
		"page":         page,
	}

	tasks, err := p.Client.GetTaskList(queryParams, mattermostUserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

// TODO: for testing purpose, remove later
func (p *Plugin) testAPI(w http.ResponseWriter, r *http.Request) {
	// TODO: remove later
	response, err := p.Client.TestApi()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

// Handles the static files under the assets directory.
func (p *Plugin) HandleStaticFiles() {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		p.API.LogWarn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	// This will serve static files from the 'assets' directory under '/static/<filename>'
	p.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(bundlePath, "assets")))))
}
