package main

import (
	"log"
	"net/http"

	"gowidget/internal/adapter/api"
	"gowidget/internal/adapter/handler"
	"gowidget/internal/config"
	"gowidget/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	githubClient := api.NewGitHubClient(cfg.GitHubToken)
	dashboardUseCase := usecase.NewDashboard(githubClient)
	httpHandler := handler.NewHTTPHandler(dashboardUseCase)

	r := mux.NewRouter()
	r.HandleFunc("/", httpHandler.Dashboard).Methods("GET")
	r.HandleFunc("/api/dashboard", httpHandler.DashboardJSON).Methods("GET")

	log.Printf("🚀 Server starting on port %s", cfg.Port)
	log.Printf("📊 Dashboard: http://localhost:%s", cfg.Port)
	log.Printf("🔗 API: http://localhost:%s/api/dashboard", cfg.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
