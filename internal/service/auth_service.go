package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/login_flow/auth-service/internal/config"
	"github.com/login_flow/auth-service/internal/domain"
	"github.com/login_flow/auth-service/pkg/crypto"
	"github.com/login_flow/auth-service/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUserNotVerified    = errors.New("email not verified")
)

type AuthService struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
	cfg       *config.Config
}

func NewAuthService(userRepo domain.UserRepository, tokenRepo domain.TokenRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		cfg:       cfg,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
		Verified: false, // Email verification required
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// TODO: Send verification email
	// For now, auto-verify for testing
	if err := s.userRepo.MarkAsVerified(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, user *domain.User, err error) {
	// Get user by email
	user, err = s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	// Verify password
	if err := crypto.ComparePassword(user.Password, password); err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	// Check if email is verified
	if !user.Verified {
		return "", "", nil, ErrUserNotVerified
	}

	// Generate access token
	accessToken, err = jwt.GenerateAccessToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.AccessExpiry)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err = crypto.GenerateRandomToken(32)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in database
	refreshTokenModel := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshExpiry),
	}

	if err := s.tokenRepo.Create(ctx, refreshTokenModel); err != nil {
		return "", "", nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, user, nil
}

// RefreshAccessToken generates a new access token using a refresh token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenStr string) (string, string, error) {
	// Get refresh token from database
	refreshToken, err := s.tokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	// Check if token is valid
	if !refreshToken.IsValid() {
		return "", "", ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, refreshToken.UserID)
	if err != nil {
		return "", "", ErrUserNotFound
	}

	// Generate new access token
	newAccessToken, err := jwt.GenerateAccessToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.AccessExpiry)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Rotate refresh token (optional but recommended)
	newRefreshToken, err := crypto.GenerateRandomToken(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Revoke old refresh token
	if err := s.tokenRepo.Revoke(ctx, refreshTokenStr); err != nil {
		return "", "", fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Store new refresh token
	newRefreshTokenModel := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshExpiry),
	}

	if err := s.tokenRepo.Create(ctx, newRefreshTokenModel); err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout revokes a refresh token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenRepo.Revoke(ctx, refreshToken)
}

// ValidateAccessToken validates an access token and returns claims
func (s *AuthService) ValidateAccessToken(tokenStr string) (*jwt.Claims, error) {
	return jwt.ValidateToken(tokenStr, s.cfg.JWT.Secret)
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
