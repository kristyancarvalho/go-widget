package handler

import (
	"encoding/json"
	"fmt"
	"gowidget/internal/usecase"
	"html/template"
	"net/http"
	"os"
)

type HTTPHandler struct {
	dashboardUseCase *usecase.Dashboard
}

func NewHTTPHandler(dashboardUseCase *usecase.Dashboard) *HTTPHandler {
	return &HTTPHandler{
		dashboardUseCase: dashboardUseCase,
	}
}

func (h *HTTPHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		http.Error(w, "GITHUB_USERNAME environment variable not set", http.StatusBadRequest)
		return
	}

	data, err := h.dashboardUseCase.GetDashboardData(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.New("dashboard").Parse(dashboardTemplate))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *HTTPHandler) DashboardJSON(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		http.Error(w, "GITHUB_USERNAME environment variable not set", http.StatusBadRequest)
		return
	}

	data, err := h.dashboardUseCase.GetDashboardData(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

const dashboardTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>GitHub Dashboard - {{.Username}}</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 20px; background: #f6f8fa; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 30px; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
        .card { background: white; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .card h2 { margin-top: 0; color: #24292f; border-bottom: 1px solid #d1d9e0; padding-bottom: 10px; }
        .commit, .repo { padding: 12px 0; border-bottom: 1px solid #f0f0f0; }
        .commit:last-child, .repo:last-child { border-bottom: none; }
        .commit-message { font-weight: 500; color: #0969da; margin-bottom: 4px; }
        .commit-meta, .repo-meta { font-size: 0.9em; color: #656d76; }
        .repo-name { font-weight: 500; color: #0969da; margin-bottom: 4px; }
        .language { display: inline-block; padding: 2px 6px; border-radius: 12px; font-size: 0.8em; background: #f3f4f6; }
        .private { color: #8b5cf6; }
        .date { color: #8b949e; }
        @media (max-width: 768px) { .grid { grid-template-columns: 1fr; } }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🐙 GitHub Dashboard</h1>
            <p>Últimas atividades de <strong>{{.Username}}</strong></p>
        </div>
        
        <div class="grid">
            <div class="card">
                <h2>📝 Últimos Commits</h2>
                {{range .Commits}}
                <div class="commit">
                    <div class="commit-message">{{.Message}}</div>
                    <div class="commit-meta">
                        📁 {{.Repository}} • 
                        <span class="date">{{.Date.Format "02/01/2006 15:04"}}</span>
                    </div>
                </div>
                {{else}}
                <p>Nenhum commit encontrado</p>
                {{end}}
            </div>
            
            <div class="card">
                <h2>📚 Repositórios Recentes</h2>
                {{range .Repositories}}
                <div class="repo">
                    <div class="repo-name">
                        {{.Name}}
                        {{if .Private}}<span class="private">🔒</span>{{end}}
                    </div>
                    <div class="repo-meta">
                        {{if .Description}}{{.Description}}<br>{{end}}
                        {{if .Language}}<span class="language">{{.Language}}</span>{{end}}
                        <span class="date">Atualizado em {{.UpdatedAt.Format "02/01/2006"}}</span>
                    </div>
                </div>
                {{else}}
                <p>Nenhum repositório encontrado</p>
                {{end}}
            </div>
        </div>
    </div>
</body>
</html>
`
