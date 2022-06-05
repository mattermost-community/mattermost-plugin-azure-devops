package routes

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/controllers"
	"github.com/gorilla/mux"
)

func InitRoutes(muxRouter *mux.Router) {
	controllers := controllers.InitControllers()
	// Add the custom plugin routes here
	muxRouter.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		controllers.Auth().SignIn(w, req)
	}).Methods(http.MethodGet)
}
