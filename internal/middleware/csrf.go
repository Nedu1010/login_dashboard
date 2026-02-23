package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/login_flow/auth-service/internal/util"
)

// CSRFMiddleware validates CSRF token for state-changing requests
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF check for GET, HEAD, OPTIONS
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Skip CSRF for non-browser clients (optional - check for custom header)
		// This allows Postman/API clients to work
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" ||
			c.GetHeader("X-API-Client") != "" {
			// Browser-based request, validate CSRF
		} else if !strings.Contains(c.GetHeader("User-Agent"), "Mozilla") {
			// Likely an API client, skip CSRF
			c.Next()
			return
		}

		// Get CSRF token from cookie
		csrfCookie, err := util.GetCookie(c, util.CSRFTokenCookie)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
			c.Abort()
			return
		}

		// Get CSRF token from header
		csrfHeader := c.GetHeader("X-CSRF-Token")
		if csrfHeader == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token required in X-CSRF-Token header"})
			c.Abort()
			return
		}

		// Validate tokens match (double-submit pattern)
		if csrfCookie != csrfHeader {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
			c.Abort()
			return
		}

		c.Next()
	}
}
