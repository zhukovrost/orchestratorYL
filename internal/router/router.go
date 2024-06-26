package router

import (
	"github.com/gorilla/mux"

	"orchestrator/internal/handlers"
	"orchestrator/internal/service"
)

// SetupRouter настраивает router
func SetupRouter(srv *service.Service) *mux.Router {
	srv.Logger.Debug("Setting up router...")

	r := mux.NewRouter()
	h := handlers.NewHandler(srv)

	r.HandleFunc("/api/v1/calculate", h.AddExpressionHandler).Methods("POST")
	r.HandleFunc("/api/v1/expressions", h.GetExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", h.GetExpressionByIdHandler).Methods("GET")
	r.HandleFunc("/internal/task", h.GetTaskHandler).Methods("GET")
	r.HandleFunc("/internal/task", h.SetResultHandler).Methods("POST")

	return r
}
