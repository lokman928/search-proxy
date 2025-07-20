package ratelimiter_test

import (
	"context"
	"time"

	"github.com/lokman928/search-proxy/internal/common/ratelimiter"
	"github.com/lokman928/search-proxy/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokenRateLimiter", func() {
	var (
		limiter *ratelimiter.TokenRateLimiter
		cfg     *config.RateLimiterConfig
	)

	BeforeEach(func() {
		cfg = &config.RateLimiterConfig{
			Enable:         true,
			MaxConcurrency: 3,
			CooldownTime:   100, // 100ms cooldown
		}
		limiter = ratelimiter.NewTokenRateLimiter(cfg)
	})

	Context("when creating a new rate limiter", func() {
		It("should initialize with the correct number of tokens", func() {
			// We can verify this by getting all tokens and ensuring we get exactly MaxConcurrency tokens
			tokens := make([]*ratelimiter.Token, 0, cfg.MaxConcurrency)

			// Get all available tokens
			for i := 0; i < cfg.MaxConcurrency; i++ {
				token, err := limiter.GetToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token).NotTo(BeNil())
				tokens = append(tokens, token)
			}

			Expect(len(tokens)).To(Equal(cfg.MaxConcurrency))
		})
	})

	Context("when getting tokens", func() {
		It("should return a token when tokens are available", func() {
			token, err := limiter.GetToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeNil())
		})

		It("should return a token with context when tokens are available", func() {
			ctx := context.Background()
			token, err := limiter.GetTokenWithContext(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeNil())
		})

		It("should block when no tokens are available", func() {
			tokens := make([]*ratelimiter.Token, 0, cfg.MaxConcurrency)

			// Get all available tokens
			for i := 0; i < cfg.MaxConcurrency; i++ {
				token, err := limiter.GetToken()
				Expect(err).NotTo(HaveOccurred())
				tokens = append(tokens, token)
			}

			// Try to get one more token in a goroutine
			gotToken := false
			go func() {
				defer GinkgoRecover()
				token, err := limiter.GetToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token).NotTo(BeNil())
				gotToken = true
			}()

			// Wait a bit to ensure the goroutine is blocked
			Consistently(func() bool { return gotToken }, 50*time.Millisecond).Should(BeFalse())

			// Release a token
			Expect(limiter.Release(tokens[0])).To(BeTrue())

			// Wait for the cooldown and a bit more to ensure the token is released
			time.Sleep(time.Duration(cfg.CooldownTime)*time.Millisecond + 50*time.Millisecond)

			// The goroutine should have received the token by now
			Eventually(func() bool { return gotToken }, 200*time.Millisecond).Should(BeTrue())
		})

		It("should return an error when context is canceled", func() {
			tokens := make([]*ratelimiter.Token, 0, cfg.MaxConcurrency)

			// Get all available tokens
			for i := 0; i < cfg.MaxConcurrency; i++ {
				token, err := limiter.GetToken()
				Expect(err).NotTo(HaveOccurred())
				tokens = append(tokens, token)
			}

			// Try to get one more token with a cancelable context
			ctx, cancel := context.WithCancel(context.Background())

			// Cancel the context after a short delay
			go func() {
				time.Sleep(50 * time.Millisecond)
				cancel()
			}()

			// This should return with an error
			token, err := limiter.GetTokenWithContext(ctx)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(context.Canceled))
			Expect(token).To(BeNil())
		})
	})

	Context("when releasing tokens", func() {
		It("should return the token to the pool after cooldown", func() {
			// Get all available tokens
			tokens := make([]*ratelimiter.Token, 0, cfg.MaxConcurrency)
			for i := 0; i < cfg.MaxConcurrency; i++ {
				token, err := limiter.GetToken()
				Expect(err).NotTo(HaveOccurred())
				tokens = append(tokens, token)
			}

			// No more tokens should be available
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			token, err := limiter.GetTokenWithContext(ctx)
			Expect(err).To(HaveOccurred())
			Expect(token).To(BeNil())

			// Release a token
			Expect(limiter.Release(tokens[0])).To(BeTrue())

			// Wait for the cooldown
			time.Sleep(time.Duration(cfg.CooldownTime)*time.Millisecond + 10*time.Millisecond)

			// Now we should be able to get a token
			token, err = limiter.GetToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeNil())
		})

		It("should allow multiple tokens to be released and reacquired", func() {
			// Get all available tokens
			tokens := make([]*ratelimiter.Token, 0, cfg.MaxConcurrency)
			for i := 0; i < cfg.MaxConcurrency; i++ {
				token, err := limiter.GetToken()
				Expect(err).NotTo(HaveOccurred())
				tokens = append(tokens, token)
			}

			// Release all tokens
			for _, token := range tokens {
				Expect(limiter.Release(token)).To(BeTrue())
			}

			// Wait for the cooldown
			time.Sleep(time.Duration(cfg.CooldownTime)*time.Millisecond + 10*time.Millisecond)

			// We should be able to get all tokens again
			newTokens := make([]*ratelimiter.Token, 0, cfg.MaxConcurrency)
			for i := 0; i < cfg.MaxConcurrency; i++ {
				token, err := limiter.GetToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token).NotTo(BeNil())
				newTokens = append(newTokens, token)
			}

			Expect(len(newTokens)).To(Equal(cfg.MaxConcurrency))
		})
	})
})
