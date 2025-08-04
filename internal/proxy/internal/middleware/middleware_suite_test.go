package middleware_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMiddleware(t *testing.T) {
	BeforeSuite(func() {
		// Set the Gin mode to TestMode to avoid any side effects during testing
		gin.SetMode(gin.TestMode)
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Middleware Suite")
}
