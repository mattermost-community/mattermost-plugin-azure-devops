package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
}

// InitAPI initializes the REST API
func (api *Api) InitRoutes(s *mux.Router) {
	// Add the custom plugin routes here
	s.HandleFunc("/test", api.test).Methods(http.MethodGet)
}

func (api *Api) test(w http.ResponseWriter, req *http.Request) {
	// if we've made it here, we're authorized.
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"is_authorized": true}`))
}
