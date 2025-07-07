package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
)

type UsernameHistoryRepository interface {
	IsUsernameExists(ctx context.Context, username string) (bool, error)
	IsUsernameChangeAllowed(ctx context.Context, oldUsername, newUsername string) (bool, error)
}

type usernameHistoryRepository struct {
	db *sql.DB
}

func NewUsernameHistoryRepository(db *sql.DB) UsernameHistoryRepository {
	return &usernameHistoryRepository{db: db}
}

func (r *usernameHistoryRepository) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	const query = `SELECT exists(SELECT 1 FROM username_history WHERE username = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, apperror.New("[CODE_USERNAME_CHECK_INVALID]", "gagal memeriksa username history", err, 500)
	}

	return exists, nil
}

func (r *usernameHistoryRepository) IsUsernameChangeAllowed(ctx context.Context, oldUsername, newUsername string) (bool, error) {
	// Jika username tidak diubah
	if oldUsername == newUsername {
		return true, nil
	}

	// Cek username
	const query = `SELECT EXISTS(SELECT 1 FROM username_history WHERE username = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, newUsername).Scan(&exists)
	if err != nil {
		return false, apperror.New("[USERNAME_HISTORY_CHECK_FAILED]", "gagal memeriksa riwayat username", err, 500)
	}

	return !exists, nil
}
