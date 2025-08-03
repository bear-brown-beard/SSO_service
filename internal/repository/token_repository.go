package repository

import (
	"database/sql"
	"fmt"
	"time"

	"sso/internal/models"

	"github.com/antibomberman/qb"
)

type TokenRepository interface {
	Create(userID int64, token, agent, ip string, expireAt time.Time) error
	FindByToken(token string) (*models.UserAccessToken, error)
	Deactivate(token string) error
	DeactivateAllUserTokens(userID int64) error
}

type tokenRepository struct {
	qb qb.QueryBuilderInterface
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{
		qb: qb.New("mysql", db),
	}
}

func (r *tokenRepository) Create(userID int64, token, agent, ip string, expireAt time.Time) error {
	if err := r.DeactivateAllUserTokens(userID); err != nil {
		return fmt.Errorf("failed to deactivate existing tokens: %w", err)
	}

	_, err := r.qb.From("user_access_tokens").CreateMap(map[string]any{
		"user_id":   userID,
		"token":     token,
		"agent":     agent,
		"ip":        ip,
		"expire_at": expireAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	return nil
}

func (r *tokenRepository) FindByToken(token string) (*models.UserAccessToken, error) {
	var accessToken models.UserAccessToken

	found, err := r.qb.From("user_access_tokens").
		Where("token = ?", token).
		Where("expire_at > NOW()").
		Limit(1).
		First(&accessToken)

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	if !found {
		return nil, nil
	}

	return &accessToken, nil
}

func (r *tokenRepository) Deactivate(token string) error {
	err := r.qb.From("user_access_tokens").
		Where("token = ?", token).
		UpdateMap(map[string]any{
			"expire_at": "NOW()",
		})

	if err != nil {
		return fmt.Errorf("failed to deactivate token: %w", err)
	}
	return nil
}

func (r *tokenRepository) DeactivateAllUserTokens(userID int64) error {
	err := r.qb.From("user_access_tokens").
		Where("user_id = ?", userID).
		UpdateMap(map[string]any{
			"expire_at": "NOW()",
		})

	if err != nil {
		return fmt.Errorf("failed to deactivate user tokens: %w", err)
	}
	return nil
}
