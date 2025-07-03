package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bytedance/sonic"
	"github.com/lokman928/search-proxy/internal/common"
	"github.com/lokman928/search-proxy/internal/common/ratelimiter"
)

type ServiceConfig struct {
	BaseUrl     string
	ApiKey      string
	RateLimiter *ratelimiter.TokenRateLimiter
}

type Service struct {
	baseUrl     string
	apiKey      string
	ratelimiter *ratelimiter.TokenRateLimiter
}

func NewService(cfg *ServiceConfig) *Service {
	return &Service{
		baseUrl:     cfg.BaseUrl,
		apiKey:      cfg.ApiKey,
		ratelimiter: cfg.RateLimiter,
	}
}

func (s *Service) Search(ctx context.Context, query string, count int) ([]common.SearchResult, error) {
	token, getTokenErr := s.ratelimiter.GetTokenWithContext(ctx)
	if getTokenErr != nil {
		return nil, fmt.Errorf("failed to get rate limiter token: %w", getTokenErr)
	}
	defer s.ratelimiter.Release(token)

	resp, err := s.sendRequestToBrave("/res/v1/web/search", map[string]string{
		"q":     query,
		"count": fmt.Sprintf("%d", count),
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	root, rootErr := sonic.Get(body)
	if rootErr != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", rootErr)
	}

	rawResults, rawResultsErr := root.Get("web").Get("results").ArrayUseNode()
	if rawResultsErr != nil {
		return nil, fmt.Errorf("failed to get web results: %w", rawResultsErr)
	}

	results := make([]common.SearchResult, 0, len(rawResults))
	for _, rawResult := range rawResults {
		url, _ := rawResult.Get("url").String()

		title, _ := rawResult.Get("title").String()

		snippet, _ := rawResult.Get("description").String()

		results = append(results, common.SearchResult{
			Link:    url,
			Title:   title,
			Snippet: snippet,
		})
	}

	return results, nil
}

func (s *Service) getUrl(path string) (string, error) {
	u, err := url.Parse(s.baseUrl)
	if err != nil {
		return "", err
	}

	return u.ResolveReference(&url.URL{Path: path}).String(), nil
}

func (s *Service) sendRequestToBrave(path string, queryParams map[string]string) (*http.Response, error) {
	url, urlErr := s.getUrl(path)
	if urlErr != nil {
		return nil, fmt.Errorf("failed to get URL: %w", urlErr)
	}

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("failed to create request: %w", reqErr)
	}

	reqUrl := req.URL.Query()
	for key, value := range queryParams {
		reqUrl.Set(key, value)
	}
	req.URL.RawQuery = reqUrl.Encode()

	req.Header = http.Header{
		"Accept":               {"application/json"},
		"X-Subscription-Token": {s.apiKey},
	}

	client := &http.Client{}
	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, fmt.Errorf("failed to send request: %w", respErr)
	}

	return resp, nil
}
