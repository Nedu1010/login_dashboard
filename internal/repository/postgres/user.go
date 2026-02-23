package postgres

import (
	"context"
	"fmt"

	"github.com/login_flow/auth-service/internal/domain"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	result := r.db.Client.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	result := r.db.Client.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", result.Error)
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Client.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", result.Error)
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	result := r.db.Client.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.Client.WithContext(ctx).Delete(&domain.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	return nil
}

func (r *UserRepository) MarkAsVerified(ctx context.Context, id int64) error {
	result := r.db.Client.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("verified", true)
	if result.Error != nil {
		return fmt.Errorf("failed to mark user as verified: %w", result.Error)
	}
	return nil
}
