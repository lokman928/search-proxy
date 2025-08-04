package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/mock"

	"github.com/lokman928/search-proxy/internal/common/logger"
	"github.com/lokman928/search-proxy/internal/proxy/internal/middleware"
)

func createTestGinContext(method, path string, body io.Reader) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(method, path, body)
	return c
}

var _ = Describe("LogRequest", func() {
	var (
		mockLogger           *logger.MockILogger
		mockLoggerWithFields *logger.MockILogger
	)

	BeforeEach(func() {
		// Initialize the mock logger
		mockLogger = logger.NewMockILogger(GinkgoT())
		mockLoggerWithFields = logger.NewMockILogger(GinkgoT())
	})

	Context("when logging requests", func() {
		It("should log the request method and path with logger", func() {
			ginMiddleware := middleware.LogRequest(mockLogger)

			mockLogger.EXPECT().WithFields(mock.MatchedBy(func(fields map[string]any) bool {
				if value, exists := fields["method"]; !exists || value != "GET" {
					return false
				}

				if value, exists := fields["path"]; !exists || value != "/test" {
					return false
				}

				if _, exists := fields["latency"]; !exists {
					return false
				}

				if value, exists := fields["status_code"]; !exists || value != 200 {
					return false
				}

				return true
			})).Return(mockLoggerWithFields)

			mockLoggerWithFields.EXPECT().Info(mock.Anything)

			ginContext := createTestGinContext("GET", "/test", nil)
			ginMiddleware(ginContext)
		})
	})
})
