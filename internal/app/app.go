package app

import (
	"net/http"
	"orchestrator/internal/router"
	"orchestrator/internal/service"
)

func Run(srv *service.Service) {
	// Настройка маршрутизатора
	r := router.SetupRouter(srv)

	srv.Logger.Infof("Starting server on %s...", srv.Cfg.GetAddress())
	// Запуск сервера
	if err := http.ListenAndServe(srv.Cfg.GetAddress(), r); err != nil {
		srv.Logger.Fatalf("Could not start server: %s\n", err.Error())
		return
	}
}
