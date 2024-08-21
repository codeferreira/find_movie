package omdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Result struct {
	Search       []SearchResult `json:"search"`
	TotalResults string         `json:"totalResults"`
	Response     string         `json:"response"`
}
type SearchResult struct {
	Title  string `json:"title"`
	Year   string `json:"year"`
	ImdbID string `json:"imdbID"`
	Type   string `json:"yype"`
	Poster string `json:"poster"`
}

func Search(apiKey string, query string) (Result, error) {
	q := url.Values{}
	q.Add("apikey", apiKey)
	q.Add("s", query)
	resp, err := http.Get("https://www.omdbapi.com/?" + q.Encode())
	if err != nil {
		return Result{}, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Result{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
