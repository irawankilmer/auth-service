package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
)

type EmailHistoryRepository interface {
	IsEmailExists(ctx context.Context, email string) (bool, error)
}

type emailHistoryRepository struct {
	db *sql.DB
}

func NewEmailHistoryRepository(db *sql.DB) EmailHistoryRepository {
	return &emailHistoryRepository{db: db}
}

func (r *emailHistoryRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	const query = `SELECT exists(SELECT 1 FROM email_history WHERE email = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, apperror.New(apperror.CodeDBError, "gagal memeriksa email", err)
	}

	return exists, nil
}
