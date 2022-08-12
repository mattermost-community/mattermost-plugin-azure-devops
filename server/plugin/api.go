package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/gorilla/mux"

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

	// Plugin APIs
	s.HandleFunc("/tasks", p.handleAuthRequired(p.handleCreateTask)).Methods(http.MethodPost)
	s.HandleFunc("/link", p.handleAuthRequired(p.handleLink)).Methods(http.MethodPost)
	s.HandleFunc(constants.PathGetAllLinkedProjects, p.handleAuthRequired(p.handleGetAllLinkedProjects)).Methods(http.MethodGet)
}

// API to create task of a project in an organization.
func (p *Plugin) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)

	body, err := serializers.CreateTaskRequestPayloadFromJSON(r.Body)
	if err != nil {
		p.API.LogError("Error in decoding the body for creating a task", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if err := body.IsValid(); err != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	task, statusCode, err := p.Client.CreateTask(body, mattermostUserID)
	if err != nil {
		p.handleError(w, r, &serializers.Error{Code: statusCode, Message: err.Error()})
		return
	}
	response, err := json.Marshal(task)
	if err != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	message := fmt.Sprintf(constants.CreatedTask, task.Link.HTML.Href)

	// Send message to DM.
	p.DM(mattermostUserID, message)
}

// API to link a project and an organization to a user.
func (p *Plugin) handleLink(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	var body *serializers.LinkRequestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		p.API.LogError("Error in decoding body", "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	if err := body.IsLinkPayloadValid(); err != "" {
		error := serializers.Error{Code: http.StatusBadRequest, Message: err}
		p.handleError(w, r, &error)
		return
	}

	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	if p.IsProjectLinked(projectList, serializers.ProjectDetails{OrganizationName: body.Organization, ProjectName: body.Project}) {
		p.DM(mattermostUserID, constants.AlreadyLinkedProject)
		return
	}

	response, err := p.Client.Link(body, mattermostUserID)
	if err != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	project := serializers.ProjectDetails{
		MattermostUserID: mattermostUserID,
		ProjectID:        response.ID,
		ProjectName:      response.Name,
		OrganizationName: body.Organization,
	}

	p.Store.StoreProject(&project)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}

// handleGetAllLinkedProjects returns all linked projects list
func (p *Plugin) handleGetAllLinkedProjects(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserIDAPI)
	projectList, err := p.Store.GetAllProjects(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		error := serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		p.handleError(w, r, &error)
		return
	}

	response, err := json.Marshal(projectList)
	if err != nil {
		p.API.LogError(constants.ErrorFetchProjectList, "Error", err.Error())
		error := serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		p.handleError(w, r, &error)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: err.Error()})
	}
}

// handleAuthRequired verifies if the provided request is performed by an authorized source.
func (p *Plugin) handleAuthRequired(handleFunc http.HandlerFunc) http.HandlerFunc {
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
	response, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
