package repository

import (
	"database/sql"
	"fmt"

	"sso/internal/models"

	"github.com/antibomberman/qb"
)

type TestAccountRepository interface {
	FindByPhone(phone string) (*models.TestAccount, error)
	UpdateCode(phone, code string) error
}

type testAccountRepository struct {
	qb qb.QueryBuilderInterface
}

func NewTestAccountRepository(db *sql.DB) TestAccountRepository {
	return &testAccountRepository{
		qb: qb.New("mysql", db),
	}
}

func (r *testAccountRepository) FindByPhone(phone string) (*models.TestAccount, error) {
	var account models.TestAccount

	found, err := r.qb.From("test_accounts").
		Where("phone = ?", phone).
		Limit(1).
		First(&account)

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	if !found {
		return nil, nil
	}

	return &account, nil
}

func (r *testAccountRepository) UpdateCode(phone, code string) error {
	err := r.qb.From("test_accounts").
		Where("phone = ?", phone).
		UpdateMap(map[string]any{
			"code": code,
		})

	if err != nil {
		return fmt.Errorf("failed to update code: %w", err)
	}

	return nil
}
