package usecase

import (
	"gowidget/internal/domain"
)

type Dashboard struct {
	githubService domain.GitHubService
}

type DashboardData struct {
	Commits      []domain.Commit     `json:"commits"`
	Repositories []domain.Repository `json:"repositories"`
	Username     string              `json:"username"`
}

func NewDashboard(githubService domain.GitHubService) *Dashboard {
	return &Dashboard{
		githubService: githubService,
	}
}

func (d *Dashboard) GetDashboardData(username string) (*DashboardData, error) {
	commits, err := d.githubService.GetUserCommits(username, 10)
	if err != nil {
		return nil, err
	}

	repositories, err := d.githubService.GetUserRepositories(username, 10)
	if err != nil {
		return nil, err
	}

	return &DashboardData{
		Commits:      commits,
		Repositories: repositories,
		Username:     username,
	}, nil
}
