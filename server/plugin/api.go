package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
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
	// TODO: WIP.
	// s.HandleFunc("/projects", p.handleAuthRequired(p.handleGetProjects)).Methods(http.MethodGet)
	// s.HandleFunc("/tasks", p.handleAuthRequired(p.handleGetTasks)).Methods(http.MethodGet)
	s.HandleFunc("/tasks", p.handleAuthRequired(p.handleCreateTask)).Methods(http.MethodPost)
}

// handleAuthRequired verifies if provided request is performed by an authorized source.
func (p *Plugin) handleAuthRequired(handleFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
		if mattermostUserID == "" {
			error := serializers.Error{Code: http.StatusUnauthorized, Message: constants.NotAuthorized}
			p.handleError(w, r, &error)
			return
		}

		handleFunc(w, r)
	}
}

func (p *Plugin) handleError(w http.ResponseWriter, r *http.Request, error *serializers.Error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(error.Code)
	message := map[string]string{constants.Error: error.Message}
	response, _ := json.Marshal(message)
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TODO: WIP.
// API to get projects in an organization.
// func (p *Plugin) handleGetProjects(w http.ResponseWriter, r *http.Request) {
// 	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)

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
// 		if _, err = w.Write([]byte(err.Error())); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	response, err := json.Marshal(boards)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		if _, err = w.Write([]byte(err.Error())); err != nil {
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
// func (p *Plugin) handleGetTasks(w http.ResponseWriter, r *http.Request) {
// 	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
// 	statusData := map[string]string{
// 		constants.Doing: "doing",
// 		constants.Todo:  "To Do",
// 		constants.Done:  "done",
// 	}
// 	organization := r.URL.Query().Get(constants.Organization)
// 	if organization == "" {
// 		error := serializers.Error{Code: http.StatusBadRequest, Message: constants.OrganizationRequired}
// 		p.handleError(w, r, &error)
// 		return
// 	}
// 	project := r.URL.Query().Get(constants.Project)
// 	if project == "" {
// 		error := serializers.Error{Code: http.StatusBadRequest, Message: constants.ProjectRequired}
// 		p.handleError(w, r, &error)
// 		return
// 	}
// 	status := r.URL.Query().Get(constants.Status)
// 	if status != "" && statusData[status] == "" {
// 		error := serializers.Error{Code: http.StatusBadRequest, Message: constants.InvalidStatus}
// 		p.handleError(w, r, &error)
// 		return
// 	}
// 	assignedTo := r.URL.Query().Get(constants.AssignedTo)
// 	if assignedTo != "" && assignedTo != "me" {
// 		error := serializers.Error{Code: http.StatusBadRequest, Message: constants.InvalidAssignedTo}
// 		p.handleError(w, r, &error)
// 		return
// 	}
// 	page := StringToInt(r.URL.Query().Get(constants.Page))
// 	if page <= 0 {
// 		error := serializers.Error{Code: http.StatusBadRequest, Message: constants.InvalidPageNumber}
// 		p.handleError(w, r, &error)
// 		return
// 	}

// 	// Wrap all query params.
// 	queryParams := map[string]interface{}{
// 		constants.Organization: organization,
// 		constants.Project:      project,
// 		constants.Status:       statusData[status],
// 		constants.AssignedTo:   assignedTo,
// 		constants.Page:         page,
// 	}

// 	tasks, err := p.Client.GetTaskList(queryParams, mattermostUserID)
// 	if err != nil {
// 		error := serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()}
// 		p.handleError(w, r, &error)
// 		return
// 	}

// 	response, err := json.Marshal(tasks)
// 	if err != nil {
// 		error := serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()}
// 		p.handleError(w, r, &error)
// 		return
// 	}
// 	w.Header().Add("Content-Type", "application/json")
// 	if _, err := w.Write(response); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// API to create task of a project in an organization.
func (p *Plugin) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	var body *serializers.TaskCreateRequestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		p.API.LogError("Error in decoding body", "Error", err.Error())
		error := serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		p.handleError(w, r, &error)
		return
	}

	if err := body.IsValid(); err != "" {
		error := serializers.Error{Code: http.StatusBadRequest, Message: err}
		p.handleError(w, r, &error)
		return
	}

	task, err := p.Client.CreateTask(body, mattermostUserID)
	if err != nil {
		error := serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		p.handleError(w, r, &error)
		return
	}
	response, err := json.Marshal(task)
	if err != nil {
		error := serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		p.handleError(w, r, &error)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err = w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	message := fmt.Sprintf(constants.CreatedTask, task.Link.HTML.Href)

	// Send message to DM.
	_, err = p.DM(mattermostUserID, message)
	_ = errors.Wrap(err, "failed to DM the created Task")
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
