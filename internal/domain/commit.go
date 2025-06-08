package domain

import "time"

type Commit struct {
	Message    string    `json:"message"`
	Repository string    `json:"repository"`
	Date       time.Time `json:"date"`
	SHA        string    `json:"sha"`
	URL        string    `json:"url"`
}
