package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/dbtx"
	"github.com/irawankilmer/auth-service/internal/model"
	"net/http"
)

type AuthRepository interface {
	IdentifierCheck(ctx context.Context, identifier string) (*model.UserModel, error)
	UpdateTokenVersion(ctx context.Context, userID, newTokenVersion string) error
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) IdentifierCheck(ctx context.Context, identifier string) (*model.UserModel, error) {
	var user model.UserModel
	var roles []model.RoleModel
	err := dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		const (
			queryUsers = `SELECT id, password, email_verified FROM users WHERE username = ? OR email = ? LIMIT 1`
			queryROles = `SELECT r.id, r.name FROM roles r JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = ?`
		)

		// query user
		err := tx.QueryRowContext(ctx, queryUsers, identifier, identifier).
			Scan(&user.ID, &user.Password, &user.EmailVerified)
		if err != nil {
			if err == sql.ErrNoRows {
				return apperror.New("[IDENTIFIER_NOT_FOUND]", "username atau email salah", err, http.StatusUnauthorized)
			}

			return apperror.New(apperror.CodeDBError, "query check identifier gagal", err)
		}

		// query roles
		rows, err := tx.QueryContext(ctx, queryROles, user.ID)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "query roles gagal", err)
		}
		defer rows.Close()

		for rows.Next() {
			var role model.RoleModel
			if err := rows.Scan(&role.ID, &role.Name); err != nil {
				return apperror.New(apperror.CodeDBError, "roles gagal scan", err)
			}

			roles = append(roles, role)
		}

		if err := rows.Err(); err != nil {
			return apperror.New(apperror.CodeDBError, "gagal setelah iterasi roles", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	user.Roles = roles

	return &user, nil
}

func (r *authRepository) UpdateTokenVersion(ctx context.Context, userID, newTokenVersion string) error {
	const query = `UPDATE users SET token_version = ? WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, query, newTokenVersion, userID); err != nil {
		return apperror.New(apperror.CodeDBError, "update token_version gagal", err)
	}

	return nil
}
