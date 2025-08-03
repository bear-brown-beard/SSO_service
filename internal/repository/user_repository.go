package repository

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"sso/internal/models"

	"github.com/antibomberman/qb"
)

type UserRepository interface {
	FindByPhone(phone string) (*models.User, error)
	Create(phone string) (*models.User, error)
}

type userRepository struct {
	qb qb.QueryBuilderInterface
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		qb: qb.New("mysql", db),
	}
}

func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User

	found, err := r.qb.From("user").
		Where("phone = ?", phone).
		WhereNull("deleted_at").
		Limit(1).
		First(&user)

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	if !found {
		return nil, nil
	}

	return &user, nil
}

func (r *userRepository) Create(phone string) (*models.User, error) {
	now := time.Now().Unix()
	username := fmt.Sprintf("user_%d", now)
	authKey := GenerateAuthKey()

	id, err := r.qb.From("user").CreateMap(map[string]any{
		"username":      username,
		"auth_key":      authKey,
		"lang":          "ru",
		"password_hash": "",
		"status":        10,
		"created_at":    now,
		"updated_at":    now,
		"phone":         phone,
		"is_guest":      false,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.User{
		ID:        id.(int64),
		Username:  username,
		AuthKey:   authKey,
		Lang:      "ru",
		Status:    10,
		CreatedAt: now,
		UpdatedAt: now,
		Phone:     sql.NullString{String: phone, Valid: true},
		IsGuest:   sql.NullBool{Bool: false, Valid: true},
	}, nil
}

func GenerateAuthKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
