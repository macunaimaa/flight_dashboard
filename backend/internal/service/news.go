package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Source      string `json:"source"`
	ImageURL    string `json:"imageUrl"`
	PublishedAt string `json:"publishedAt"`
}

type NewsSearchResult struct {
	Articles []NewsArticle `json:"articles"`
	Total    int           `json:"total"`
}

type NewsService struct {
	apiKey     string
	httpClient *http.Client
}

func NewNewsService(apiKey string) *NewsService {
	return &NewsService{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *NewsService) Search(ctx context.Context, query string, lang string) (*NewsSearchResult, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("NEWS_API_KEY not configured")
	}

	params := url.Values{}
	params.Set("q", query)
	params.Set("lang", lang)
	params.Set("token", s.apiKey)
	params.Set("max", "10")

	reqURL := fmt.Sprintf("https://gnews.io/api/v4/search?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("news API returned status %d", resp.StatusCode)
	}

	var raw struct {
		Articles []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			Source      struct {
				Name string `json:"name"`
			} `json:"source"`
			Image       string `json:"image"`
			PublishedAt string `json:"publishedAt"`
		} `json:"articles"`
		TotalArticles int `json:"totalArticles"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	result := &NewsSearchResult{
		Total: raw.TotalArticles,
	}

	for _, a := range raw.Articles {
		result.Articles = append(result.Articles, NewsArticle{
			Title:       a.Title,
			Description: a.Description,
			URL:         a.URL,
			Source:      a.Source.Name,
			ImageURL:    a.Image,
			PublishedAt: a.PublishedAt,
		})
	}

	return result, nil
}
