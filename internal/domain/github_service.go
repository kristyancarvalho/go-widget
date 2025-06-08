package domain

type GitHubService interface {
	GetUserCommits(username string, limit int) ([]Commit, error)
	GetUserRepositories(username string, limit int) ([]Repository, error)
}
