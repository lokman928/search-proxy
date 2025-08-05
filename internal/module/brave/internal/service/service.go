package service

import (
	"context"
	"fmt"

	"github.com/lokman928/search-proxy/internal/common"
	"github.com/lokman928/search-proxy/internal/common/logger"
	"github.com/lokman928/search-proxy/internal/common/ratelimiter"
	"github.com/lokman928/search-proxy/internal/module/brave/internal/remote"
)

type ServiceConfig struct {
	BaseUrl     string
	ApiKey      string
	RateLimiter *ratelimiter.TokenRateLimiter
}

type Service struct {
	braveClient remote.IBraveClient
	ratelimiter *ratelimiter.TokenRateLimiter
	logger      logger.ILogger
}

func NewService(cfg *ServiceConfig) *Service {
	return &Service{
		braveClient: remote.NewBraveClient(cfg.BaseUrl, cfg.ApiKey),
		ratelimiter: cfg.RateLimiter,
		logger:      logger.NewLogger("BraveService"),
	}
}

func (svc *Service) Search(ctx context.Context, query string, count int) ([]common.SearchResult, error) {
	svc.logger.WithFields(map[string]any{
		"query": query,
		"count": count,
	}).Info("Start searching with Brave")

	token, getTokenErr := svc.ratelimiter.GetTokenWithContext(ctx)
	if getTokenErr != nil {
		return nil, fmt.Errorf("failed to get rate limiter token: %w", getTokenErr)
	}
	defer svc.ratelimiter.Release(token)

	searchResp, searchErr := svc.braveClient.Search(ctx, remote.SearchParams{
		Query: query,
		Count: count,
	})

	if searchErr != nil {
		svc.logger.WithField("error", searchErr.Error()).Error("Failed to search with Brave")
		return nil, fmt.Errorf("failed to search by Brave API: %w", searchErr)
	}

	// Convert the search response to SearchResult
	searchResult := searchResp.Web.Results
	results := make([]common.SearchResult, 0, len(searchResult))
	for _, res := range searchResult {
		results = append(results, common.SearchResult{
			Link:    res.URL,
			Title:   res.Title,
			Snippet: res.Description,
		})
	}

	svc.logger.Info("Search completed with Brave")

	return results, nil
}
