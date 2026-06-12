package search

import "context"

// Result represents a single web search result.
type Result struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// Provider is the abstraction for web search backends.
type Provider interface {
	Search(ctx context.Context, query string) ([]Result, error)
}
