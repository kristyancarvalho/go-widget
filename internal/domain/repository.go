package domain

import "time"

type Repository struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Private     bool      `json:"private"`
	URL         string    `json:"url"`
}
