package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// DuckDuckGoProvider searches DuckDuckGo's instant answer API.
type DuckDuckGoProvider struct {
	client *http.Client
	limiter *tokenBucket
}

// NewDuckDuckGoProvider creates a rate-limited DuckDuckGo provider.
func NewDuckDuckGoProvider() *DuckDuckGoProvider {
	return &DuckDuckGoProvider{
		client:  &http.Client{Timeout: 10 * time.Second},
		limiter: newTokenBucket(time.Second, 1),
	}
}

type ddgResponse struct {
	AbstractText string `json:"AbstractText"`
	AbstractURL  string `json:"AbstractURL"`
	Heading      string `json:"Heading"`
	RelatedTopics []struct {
		Text string `json:"Text"`
		FirstURL string `json:"FirstURL"`
		Result string `json:"Result"`
	} `json:"RelatedTopics"`
}

// Search performs a DuckDuckGo instant answer search.
func (p *DuckDuckGoProvider) Search(ctx context.Context, query string) ([]Result, error) {
	p.limiter.Wait()

	endpoint := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=json&no_html=1&skip_disambig=1", url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MyQQBot/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("duckduckgo: unexpected status %d", resp.StatusCode)
	}

	var data ddgResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("duckduckgo: decode failed: %w", err)
	}

	var results []Result
	if data.AbstractText != "" && data.AbstractURL != "" {
		results = append(results, Result{
			Title:   data.Heading,
			URL:     data.AbstractURL,
			Snippet: data.AbstractText,
		})
	}

	for _, rt := range data.RelatedTopics {
		if rt.Text == "" {
			continue
		}
		results = append(results, Result{
			Title:   rt.FirstURL,
			URL:     rt.FirstURL,
			Snippet: rt.Text,
		})
	}

	return results, nil
}

// tokenBucket is a trivial 1 req/s rate limiter.
type tokenBucket struct {
	mu       sync.Mutex
	last     time.Time
	interval time.Duration
}

func newTokenBucket(interval time.Duration, _ int) *tokenBucket {
	return &tokenBucket{interval: interval}
}

func (tb *tokenBucket) Wait() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	if elapsed := time.Since(tb.last); elapsed < tb.interval {
		time.Sleep(tb.interval - elapsed)
	}
	tb.last = time.Now()
}
