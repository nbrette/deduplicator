package router

import (
	"deduplicator/internal/controllers"

	"github.com/gorilla/mux"
)

func InitializeRouter(controller *controllers.DeduplicatorController) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", controller.HandlePost).Methods("POST")
	return router
}
