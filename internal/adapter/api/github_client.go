package api

import (
	"encoding/json"
	"fmt"
	"gowidget/internal/domain"
	"net/http"
	"sort"
	"time"
)

type GitHubClient struct {
	token      string
	httpClient *http.Client
}

type GitHubCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Date string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
	HTMLURL string `json:"html_url"`
}

type GitHubRepository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Private     bool   `json:"private"`
	UpdatedAt   string `json:"updated_at"`
	HTMLURL     string `json:"html_url"`
}

type GitHubEvent struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Commits []struct {
			SHA     string `json:"sha"`
			Message string `json:"message"`
		} `json:"commits,omitempty"`
	} `json:"payload"`
	CreatedAt string `json:"created_at"`
}

func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *GitHubClient) makeRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if g.token != "" {
		req.Header.Set("Authorization", "token "+g.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	return g.httpClient.Do(req)
}

func (g *GitHubClient) GetUserCommits(username string, limit int) ([]domain.Commit, error) {
	// Get user events to find push events with commits
	url := fmt.Sprintf("https://api.github.com/users/%s/events?per_page=100", username)
	resp, err := g.makeRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var events []GitHubEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	var commits []domain.Commit
	for _, event := range events {
		if event.Type == "PushEvent" && len(event.Payload.Commits) > 0 {
			for _, commit := range event.Payload.Commits {
				date, _ := time.Parse(time.RFC3339, event.CreatedAt)

				commits = append(commits, domain.Commit{
					Message:    commit.Message,
					Repository: event.Repo.Name,
					Date:       date,
					SHA:        commit.SHA,
					URL:        fmt.Sprintf("https://github.com/%s/commit/%s", event.Repo.Name, commit.SHA),
				})
			}
		}
	}

	// Sort by date (newest first)
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})

	// Limit results
	if len(commits) > limit {
		commits = commits[:limit]
	}

	return commits, nil
}

func (g *GitHubClient) GetUserRepositories(username string, limit int) ([]domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=%d", username, limit)
	resp, err := g.makeRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var githubRepos []GitHubRepository
	if err := json.NewDecoder(resp.Body).Decode(&githubRepos); err != nil {
		return nil, err
	}

	var repositories []domain.Repository
	for _, repo := range githubRepos {
		updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)

		repositories = append(repositories, domain.Repository{
			Name:        repo.Name,
			FullName:    repo.FullName,
			UpdatedAt:   updatedAt,
			Description: repo.Description,
			Language:    repo.Language,
			Private:     repo.Private,
			URL:         repo.HTMLURL,
		})
	}

	return repositories, nil
}
