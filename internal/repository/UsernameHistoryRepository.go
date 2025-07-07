package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
)

type UsernameHistoryRepository interface {
	IsUsernameExists(ctx context.Context, newUsername string) (bool, error)
}

type usernameHistoryRepository struct {
	db *sql.DB
}

func NewUsernameHistoryRepository(db *sql.DB) UsernameHistoryRepository {
	return &usernameHistoryRepository{db: db}
}

func (r *usernameHistoryRepository) IsUsernameExists(ctx context.Context, newUsername string) (bool, error) {
	const query = `SELECT exists(SELECT 1 FROM username_history WHERE username = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, newUsername).Scan(&exists)
	if err != nil {
		return false, apperror.New("[CODE_USERNAME_CHECK_INVALID]", "gagal memeriksa username history", err, 500)
	}

	return exists, nil
}
