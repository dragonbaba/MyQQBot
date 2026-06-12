package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TavilyProvider searches using the Tavily API.
type TavilyProvider struct {
	apiKey string
	client *http.Client
}

// NewTavilyProvider creates a new Tavily provider.
func NewTavilyProvider(apiKey string) *TavilyProvider {
	return &TavilyProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type tavilyRequest struct {
	APIKey      string `json:"api_key"`
	Query       string `json:"query"`
	SearchDepth string `json:"search_depth"`
	MaxResults  int    `json:"max_results"`
}

type tavilyResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

type tavilyResponse struct {
	Results []tavilyResult `json:"results"`
}

// Search performs a Tavily search.
func (p *TavilyProvider) Search(ctx context.Context, query string) ([]Result, error) {
	if p.apiKey == "" {
		return nil, errors.New("Tavily API key not configured")
	}

	body, err := json.Marshal(tavilyRequest{
		APIKey:      p.apiKey,
		Query:       query,
		SearchDepth: "basic",
		MaxResults:  5,
	})
	if err != nil {
		return nil, fmt.Errorf("tavily: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.tavily.com/search", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("tavily: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tavily: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("tavily: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var data tavilyResponse
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, fmt.Errorf("tavily: decode failed: %w", err)
	}

	results := make([]Result, 0, len(data.Results))
	for _, r := range data.Results {
		results = append(results, Result{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Content,
		})
	}

	return results, nil
}
