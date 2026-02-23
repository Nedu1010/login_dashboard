package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/login_flow/auth-service/internal/config"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
	CSRFTokenCookie    = "csrf_token"
)

// SetAccessTokenCookie sets the access token as HTTP-only cookie
func SetAccessTokenCookie(c *gin.Context, token string, cfg *config.CookieConfig, maxAge int) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		AccessTokenCookie,
		token,
		maxAge,
		"/",
		cfg.Domain,
		cfg.Secure,
		true, // HttpOnly
	)
}

// SetRefreshTokenCookie sets the refresh token as HTTP-only cookie
func SetRefreshTokenCookie(c *gin.Context, token string, cfg *config.CookieConfig, maxAge int) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		RefreshTokenCookie,
		token,
		maxAge,
		"/",
		cfg.Domain,
		cfg.Secure,
		true, // HttpOnly
	)
}

// SetCSRFTokenCookie sets the CSRF token as readable cookie (NOT HttpOnly)
func SetCSRFTokenCookie(c *gin.Context, token string, cfg *config.CookieConfig, maxAge int) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		CSRFTokenCookie,
		token,
		maxAge,
		"/",
		cfg.Domain,
		cfg.Secure,
		false, // NOT HttpOnly - JS needs to read this
	)
}

// ClearAuthCookies removes all authentication cookies
func ClearAuthCookies(c *gin.Context, cfg *config.CookieConfig) {
	c.SetCookie(AccessTokenCookie, "", -1, "/", cfg.Domain, cfg.Secure, true)
	c.SetCookie(RefreshTokenCookie, "", -1, "/", cfg.Domain, cfg.Secure, true)
	c.SetCookie(CSRFTokenCookie, "", -1, "/", cfg.Domain, cfg.Secure, false)
}

// GetCookie retrieves a cookie value
func GetCookie(c *gin.Context, name string) (string, error) {
	return c.Cookie(name)
}
