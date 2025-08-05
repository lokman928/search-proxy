package remote

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/lokman928/search-proxy/internal/common/logger"
)

type IBraveClient interface {
	Search(context.Context, SearchParams) (*SearchResponse, error)
}

var _ IBraveClient = (*BraveClient)(nil)

type BraveClient struct {
	client *resty.Client
	logger logger.ILogger
}

func NewBraveClient(baseUrl, apiKey string) *BraveClient {
	client := resty.New()
	client.SetHeaders(map[string]string{
		"Accept":               "application/json",
		"Content-Type":         "application/json",
		"X-Subscription-Token": apiKey,
	})
	client.SetBaseURL(baseUrl)

	return &BraveClient{
		client: client,
		logger: logger.NewLogger("BraveClient"),
	}
}

type SearchParams struct {
	Query string
	Count int
}

type SearchResponse struct {
	Web SearchWebResponse `json:"web"`
}

type SearchWebResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

func (c *BraveClient) Search(ctx context.Context, searchParams SearchParams) (*SearchResponse, error) {
	searchResponse := &SearchResponse{}

	c.logger.WithField("search_params", searchParams).Debug("Sending search request to Brave API")
	resp, respErr := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"q":     searchParams.Query,
			"count": fmt.Sprintf("%d", searchParams.Count),
		}).
		SetResult(searchResponse).
		Get("/res/v1/web/search")

	if respErr != nil {
		c.logger.WithField("error", respErr.Error()).Error("Failed to send request to Brave API")
		return nil, fmt.Errorf("failed to send request to Brave: %w", respErr)
	}

	if resp.IsError() {
		c.logger.WithField("response", resp.String()).Error("Brave API returned an error")
		return nil, fmt.Errorf("Brave API returned error: %s", resp.String())
	}

	c.logger.WithField("response", resp.String()).Debug("Received response from Brave API")

	return searchResponse, nil
}
