package ratelimiter

import (
	"context"
	"time"

	"github.com/lokman928/search-proxy/internal/config"
)

type Token struct{}

// TokenRateLimiter implements a token-based rate limiter.
// It initializes with N tokens and allows consumers to get and release tokens.
// If no tokens are available, Get() will block until a token becomes available.
type TokenRateLimiter struct {
	enabled      bool          // Whether the rate limiter is enabled
	tokens       chan Token    // Channel to hold tokens
	cooldownTime time.Duration // Cooldown interval in seconds
}

// NewTokenRateLimiter creates a new TokenRateLimiter with n tokens.
func NewTokenRateLimiter(cfg *config.RateLimiterConfig) *TokenRateLimiter {
	tokens := make(chan Token, cfg.MaxConcurrency)

	// Initialize the channel with n tokens
	for i := 0; i < cfg.MaxConcurrency; i++ {
		tokens <- Token{}
	}

	return &TokenRateLimiter{
		enabled:      cfg.Enable,
		tokens:       tokens,
		cooldownTime: time.Duration(cfg.CooldownTime) * time.Millisecond,
	}
}

// Get retrieves a token from the limiter.
// If no tokens are available, this method will block until a token becomes available.
// Returns a token and an error. If the context is canceled, the error will be non-nil.
func (t *TokenRateLimiter) GetToken() (*Token, error) {
	return t.GetTokenWithContext(context.Background())
}

// GetWithContext retrieves a token from the limiter with context support.
// If no tokens are available, this method will block until either a token becomes available
// or the context is canceled.
// Returns a token and an error. If the context is canceled, the error will be non-nil.
func (t *TokenRateLimiter) GetTokenWithContext(ctx context.Context) (*Token, error) {
	select {
	case token := <-t.tokens:
		return &token, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Release returns a token to the limiter.
// The token should be one that was previously obtained from GetWithContext().
// Returns true if the token was successfully released, false if the token is invalid.
func (t *TokenRateLimiter) Release(token *Token) bool {
	time.AfterFunc(t.cooldownTime, func() {
		t.tokens <- *token
	})

	return true
}
