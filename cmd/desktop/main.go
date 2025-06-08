package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"gowidget/internal/adapter/api"
	"gowidget/internal/config"
	"gowidget/internal/domain"
	"gowidget/internal/usecase"
)

type GitHubWidget struct {
	app           fyne.App
	window        fyne.Window
	githubUseCase *usecase.Dashboard
	username      string

	commitsContainer *fyne.Container
	reposContainer   *fyne.Container
	statusLabel      *widget.Label
	refreshButton    *widget.Button
}

func main() {
	cfg := config.Load()

	if cfg.GitHubUsername == "" {
		log.Fatal("❌ GITHUB_USERNAME environment variable is required")
	}

	githubClient := api.NewGitHubClient(cfg.GitHubToken)
	dashboardUseCase := usecase.NewDashboard(githubClient)

	widget := NewGitHubWidget(dashboardUseCase, cfg.GitHubUsername)
	widget.Show()
}

func NewGitHubWidget(useCase *usecase.Dashboard, username string) *GitHubWidget {
	myApp := app.New()
	myApp.SetIcon(theme.ComputerIcon())

	w := &GitHubWidget{
		app:           myApp,
		githubUseCase: useCase,
		username:      username,
	}

	w.setupUI()
	w.loadData()

	return w
}

func (w *GitHubWidget) setupUI() {
	w.window = w.app.NewWindow(fmt.Sprintf("GitHub Dashboard - %s", w.username))
	w.window.Resize(fyne.NewSize(800, 600))
	w.window.SetFixedSize(false)

	title := widget.NewRichTextFromMarkdown("# 🐙 GitHub Dashboard")
	title.Wrapping = fyne.TextWrapWord

	userLabel := widget.NewLabel(fmt.Sprintf("Usuário: %s", w.username))

	w.statusLabel = widget.NewLabel("Carregando...")
	w.refreshButton = widget.NewButton("🔄 Atualizar", w.loadData)

	headerContainer := container.NewVBox(
		title,
		userLabel,
		container.NewHBox(layout.NewSpacer(), w.statusLabel, w.refreshButton, layout.NewSpacer()),
		widget.NewSeparator(),
	)

	w.commitsContainer = container.NewVBox()
	w.reposContainer = container.NewVBox()

	commitsScroll := container.NewScroll(w.commitsContainer)
	commitsScroll.SetMinSize(fyne.NewSize(380, 400))

	reposScroll := container.NewScroll(w.reposContainer)
	reposScroll.SetMinSize(fyne.NewSize(380, 400))

	commitsCard := container.NewBorder(
		widget.NewRichTextFromMarkdown("## 📝 Últimos Commits"),
		nil, nil, nil,
		commitsScroll,
	)

	reposCard := container.NewBorder(
		widget.NewRichTextFromMarkdown("## 📚 Repositórios Recentes"),
		nil, nil, nil,
		reposScroll,
	)

	content := container.NewHSplit(commitsCard, reposCard)
	content.SetOffset(0.5)

	mainContainer := container.NewBorder(
		headerContainer,
		nil, nil, nil,
		content,
	)

	w.window.SetContent(mainContainer)

	go w.autoRefresh()
}

func (w *GitHubWidget) loadData() {
	w.statusLabel.SetText("Carregando...")
	w.refreshButton.Disable()

	go func() {
		data, err := w.githubUseCase.GetDashboardData(w.username)
		if err != nil {
			w.statusLabel.SetText(fmt.Sprintf("❌ Erro: %v", err))
			w.refreshButton.Enable()
			return
		}

		w.app.SendNotification(fyne.NewNotification("GitHub Widget", "Dados atualizados com sucesso!"))

		w.updateCommits(data.Commits)
		w.updateRepositories(data.Repositories)
		w.statusLabel.SetText(fmt.Sprintf("✅ Atualizado às %s", time.Now().Format("15:04:05")))
		w.refreshButton.Enable()
	}()
}

func (w *GitHubWidget) updateCommits(commits []domain.Commit) {
	w.commitsContainer.Objects = nil

	if len(commits) == 0 {
		w.commitsContainer.Add(widget.NewLabel("Nenhum commit encontrado"))
		w.commitsContainer.Refresh()
		return
	}

	for _, commit := range commits {
		commitCard := w.createCommitCard(commit)
		w.commitsContainer.Add(commitCard)
		w.commitsContainer.Add(widget.NewSeparator())
	}

	w.commitsContainer.Refresh()
}

func (w *GitHubWidget) updateRepositories(repos []domain.Repository) {
	w.reposContainer.Objects = nil

	if len(repos) == 0 {
		w.reposContainer.Add(widget.NewLabel("Nenhum repositório encontrado"))
		w.reposContainer.Refresh()
		return
	}

	for _, repo := range repos {
		repoCard := w.createRepoCard(repo)
		w.reposContainer.Add(repoCard)
		w.reposContainer.Add(widget.NewSeparator())
	}

	w.reposContainer.Refresh()
}

func (w *GitHubWidget) createCommitCard(commit domain.Commit) *fyne.Container {
	message := widget.NewRichTextFromMarkdown(fmt.Sprintf("**%s**", commit.Message))
	message.Wrapping = fyne.TextWrapWord

	meta := widget.NewLabel(fmt.Sprintf("📁 %s • %s",
		commit.Repository,
		commit.Date.Format("02/01/2006 15:04")))

	openButton := widget.NewButton("🔗 Ver Commit", func() {
		w.app.OpenURL(parseURL(commit.URL))
	})
	openButton.Importance = widget.LowImportance

	return container.NewVBox(
		message,
		meta,
		openButton,
	)
}

func (w *GitHubWidget) createRepoCard(repo domain.Repository) *fyne.Container {
	name := widget.NewRichTextFromMarkdown(fmt.Sprintf("**%s** %s",
		repo.Name,
		func() string {
			if repo.Private {
				return "🔒"
			}
			return ""
		}()))

	description := ""
	if repo.Description != "" {
		description = repo.Description
	}

	meta := fmt.Sprintf("%s\n", description)
	if repo.Language != "" {
		meta += fmt.Sprintf("🏷️ %s • ", repo.Language)
	}
	meta += fmt.Sprintf("Atualizado em %s", repo.UpdatedAt.Format("02/01/2006"))

	metaLabel := widget.NewLabel(meta)
	metaLabel.Wrapping = fyne.TextWrapWord

	openButton := widget.NewButton("🔗 Ver Repositório", func() {
		w.app.OpenURL(parseURL(repo.URL))
	})
	openButton.Importance = widget.LowImportance

	return container.NewVBox(
		name,
		metaLabel,
		openButton,
	)
}

func (w *GitHubWidget) autoRefresh() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		w.loadData()
	}
}

func (w *GitHubWidget) Show() {
	w.window.ShowAndRun()
}

func parseURL(urlStr string) *url.URL {
	u, _ := url.Parse(urlStr)
	return u
}
