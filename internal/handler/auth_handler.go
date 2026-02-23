package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/login_flow/auth-service/internal/config"
	"github.com/login_flow/auth-service/internal/service"
	"github.com/login_flow/auth-service/internal/util"
	"github.com/login_flow/auth-service/pkg/validator"
)

type AuthHandler struct {
	authService *service.AuthService
	csrfService *service.CSRFService
	cfg         *config.Config
}

func NewAuthHandler(authService *service.AuthService, csrfService *service.CSRFService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		csrfService: csrfService,
		cfg:         cfg,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req validator.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validator.ValidatePassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password must be at least 8 characters and contain uppercase, lowercase, and number",
		})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful",
		"user":    user.ToResponse(),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req validator.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accessToken, refreshToken, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if errors.Is(err, service.ErrUserNotVerified) {
			c.JSON(http.StatusForbidden, gin.H{"error": "email not verified"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	util.SetAccessTokenCookie(c, accessToken, &h.cfg.Cookie, int(h.cfg.JWT.AccessExpiry.Seconds()))
	util.SetRefreshTokenCookie(c, refreshToken, &h.cfg.Cookie, int(h.cfg.JWT.RefreshExpiry.Seconds()))

	csrfToken := h.csrfService.GenerateToken()
	util.SetCSRFTokenCookie(c, csrfToken, &h.cfg.Cookie, int(h.cfg.JWT.RefreshExpiry.Seconds()))

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    user.ToResponse(),
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := util.GetCookie(c, util.RefreshTokenCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	newAccessToken, newRefreshToken, err := h.authService.RefreshAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	util.SetAccessTokenCookie(c, newAccessToken, &h.cfg.Cookie, int(h.cfg.JWT.AccessExpiry.Seconds()))
	util.SetRefreshTokenCookie(c, newRefreshToken, &h.cfg.Cookie, int(h.cfg.JWT.RefreshExpiry.Seconds()))

	csrfToken := h.csrfService.GenerateToken()
	util.SetCSRFTokenCookie(c, csrfToken, &h.cfg.Cookie, int(h.cfg.JWT.RefreshExpiry.Seconds()))

	c.JSON(http.StatusOK, gin.H{
		"message": "token refreshed successfully",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := util.GetCookie(c, util.RefreshTokenCookie)
	if err != nil {
		util.ClearAuthCookies(c, &h.cfg.Cookie)
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	util.ClearAuthCookies(c, &h.cfg.Cookie)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
