package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/login_flow/auth-service/internal/domain"
)

type TokenRepository struct {
	db *DB
}

func NewTokenRepository(db *DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	result := r.db.Client.WithContext(ctx).Create(token)
	if result.Error != nil {
		return fmt.Errorf("failed to create refresh token: %w", result.Error)
	}
	return nil
}

func (r *TokenRepository) GetByToken(ctx context.Context, tokenStr string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	result := r.db.Client.WithContext(ctx).
		Where("token = ?", tokenStr).
		First(&token)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", result.Error)
	}
	return &token, nil
}

func (r *TokenRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.RefreshToken, error) {
	var tokens []*domain.RefreshToken
	result := r.db.Client.WithContext(ctx).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Order("created_at DESC").
		Find(&tokens)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get refresh tokens: %w", result.Error)
	}
	return tokens, nil
}

func (r *TokenRepository) Revoke(ctx context.Context, tokenStr string) error {
	result := r.db.Client.WithContext(ctx).Model(&domain.RefreshToken{}).Where("token = ?", tokenStr).Update("revoked_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", result.Error)
	}
	return nil
}

func (r *TokenRepository) RevokeAllForUser(ctx context.Context, userID int64) error {
	result := r.db.Client.WithContext(ctx).Model(&domain.RefreshToken{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to revoke all tokens for user: %w", result.Error)
	}
	return nil
}

func (r *TokenRepository) CleanupExpired(ctx context.Context) error {
	result := r.db.Client.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&domain.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", result.Error)
	}
	return nil
}
