package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/model"
	"net/http"
)

type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *model.EmailVerificationModel) error
	FindByToken(ctx context.Context, token string) (*model.EmailVerificationModel, error)
	MarkAsUsed(ctx context.Context, evID string) error
}

type emailVerificationRepository struct {
	db *sql.DB
}

func NewEmailVerificationRepository(db *sql.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(ctx context.Context, ev *model.EmailVerificationModel) error {
	const query = `INSERT INTO email_verifications(id, user_id, token, expires_at) VALUES(?, ?, ?, ?)`
	if _, err := r.db.ExecContext(ctx, query, ev.ID, ev.UserID, ev.Token, ev.ExpiresAt); err != nil {
		return apperror.New(apperror.CodeDBError, "query insert email_verifications gagal", err)
	}

	return nil
}

func (r *emailVerificationRepository) FindByToken(ctx context.Context, token string) (*model.EmailVerificationModel, error) {
	const query = `SELECT id, user_id, expires_at, is_used FROM email_verifications WHERE token = ? LIMIT 1`
	var ev model.EmailVerificationModel
	if err := r.db.QueryRowContext(ctx, query, token).Scan(&ev.ID, &ev.UserID, &ev.ExpiresAt, &ev.IsUsed); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.New("[TOKEN_NOT_FOUND]", "token tidak ditemukan", err, http.StatusUnauthorized)
		}

		return nil, apperror.New(apperror.CodeDBError, "query select email verification by token gagal", err)
	}

	return &ev, nil
}

func (r *emailVerificationRepository) MarkAsUsed(ctx context.Context, evID string) error {
	const query = `UPDATE email_verifications SET is_used = true WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, query, evID); err != nil {
		return apperror.New(apperror.CodeDBError, "query update is_used gagal", err)
	}

	return nil
}
