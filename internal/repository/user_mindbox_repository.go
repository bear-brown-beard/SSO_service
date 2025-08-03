package repository

import (
	"database/sql"
	"fmt"

	"sso/internal/models"

	"github.com/antibomberman/qb"
)

type UserMindBoxRepository interface {
	FindByUserID(userID int64) (*models.UserMindBox, error)
	Create(userID int64) (*models.UserMindBox, error)
}

type userMindBoxRepository struct {
	qb qb.QueryBuilderInterface
}

func NewUserMindBoxRepository(db *sql.DB) UserMindBoxRepository {
	return &userMindBoxRepository{
		qb: qb.New("mysql", db),
	}
}

func (r *userMindBoxRepository) FindByUserID(userID int64) (*models.UserMindBox, error) {
	var userMindBox models.UserMindBox

	found, err := r.qb.From("user_mind_box").
		Where("user_id = ?", userID).
		Limit(1).
		First(&userMindBox)

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	if !found {
		return nil, nil
	}

	return &userMindBox, nil
}

func (r *userMindBoxRepository) Create(userID int64) (*models.UserMindBox, error) {
	id, err := r.qb.From("user_mind_box").CreateMap(map[string]any{
		"user_id":                  userID,
		"consent_to_mailings":      false,
		"loyalty_program_enrolled": false,
		"is_phone_confirm":         false,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create record in user_mind_box: %w", err)
	}

	return &models.UserMindBox{
		ID:                     id.(int64),
		UserID:                 sql.NullInt64{Int64: userID, Valid: true},
		ConsentToMailings:      sql.NullBool{Bool: false, Valid: true},
		LoyaltyProgramEnrolled: sql.NullBool{Bool: false, Valid: true},
		IsPhoneConfirm:         sql.NullBool{Bool: false, Valid: true},
	}, nil
}
