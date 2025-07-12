package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/dbtx"
	"github.com/irawankilmer/auth-service/internal/model"
	"net/http"
)

type UserSessionRepository interface {
	Create(ctx context.Context, data *model.UserSession) error
	FindRefreshToken(ctx context.Context, hashed string) (*model.UserSession, error)
	GetTokenVersionByUserID(ctx context.Context, userID string) (*model.UserModel, error)
	Revoked(ctx context.Context, usID string) error
	RevokeAllSessionByUserID(ctx context.Context, userID string) error
}

type userSessionRepositoryImpl struct {
	db *sql.DB
}

func NewUserSessionRepository(db *sql.DB) UserSessionRepository {
	return &userSessionRepositoryImpl{db: db}
}

func (r *userSessionRepositoryImpl) Create(ctx context.Context, data *model.UserSession) error {
	const query = `
									INSERT
									INTO user_sessions
										(id, user_id, refresh_token_hash, device_id, ip_address, user_agent, expires_at)
									VALUES(?, ?, ?, ?, ?, ?, ?)
								`
	if _, err := r.db.ExecContext(ctx, query,
		data.ID, data.UserID, data.RefreshTokenHash, data.DeviceID, data.IPAddress, data.UserAgent, data.ExpiresAt,
	); err != nil {
		return apperror.New(apperror.CodeDBError, "query user sessions gagal", err)
	}

	return nil
}

func (r *userSessionRepositoryImpl) FindRefreshToken(ctx context.Context, hashed string) (*model.UserSession, error) {
	const query = `SELECT id, user_id, refresh_token_hash, revoked, expires_at 
									FROM user_sessions WHERE refresh_token_hash = ?`
	var us model.UserSession
	if err := r.db.QueryRowContext(ctx, query, hashed).Scan(
		&us.ID, &us.UserID, &us.RefreshTokenHash, &us.Revoked, &us.ExpiresAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.New("[REFRESH_TOKEN_NOT_FOUND]", "refresh token tidak ditemukan", err, http.StatusUnauthorized)
		}

		return nil, apperror.New(apperror.CodeDBError, "query refresh token gagal", err)
	}

	return &us, nil
}

func (r *userSessionRepositoryImpl) GetTokenVersionByUserID(ctx context.Context, userID string) (*model.UserModel, error) {
	const (
		queryUser  = `SELECT id, token_version, email_verified FROM users WHERE id = ?`
		queryRoles = `SELECT r.id, r.name FROM roles r INNER JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = ?`
	)

	user := model.UserModel{
		Roles: []model.RoleModel{},
	}

	// stored procedure
	err := dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		// query user
		if err := tx.QueryRowContext(ctx, queryUser, userID).Scan(&user.ID, &user.TokenVersion, &user.EmailVerified); err != nil {
			if err == sql.ErrNoRows {
				return apperror.New(apperror.CodeUserNotFound, "user tidak ditemukan", err)
			}

			return apperror.New(apperror.CodeDBError, "query users gagal", err)
		}

		// query roles
		rows, err := tx.QueryContext(ctx, queryRoles, user.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				return apperror.New(apperror.CodeRoleNotFound, "role dari user tersebut tidak ditemukan", err)
			}

			return apperror.New(apperror.CodeDBError, "query roles gagal", err)
		}
		defer rows.Close()

		// scan roles
		for rows.Next() {
			var role model.RoleModel
			if err := rows.Scan(&role.ID, &role.Name); err != nil {
				return apperror.New(apperror.CodeDBError, "scan roles gagal", err)
			}

			user.Roles = append(user.Roles, role)
		}

		// cek error rows
		if err := rows.Err(); err != nil {
			return apperror.New(apperror.CodeDBError, "gagal setelah iterasi", err)
		}

		return nil
	})

	// cek error stored procedure
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userSessionRepositoryImpl) Revoked(ctx context.Context, usID string) error {
	const query = `UPDATE user_sessions SET revoked = true WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, query, usID); err != nil {
		return apperror.New(apperror.CodeDBError, "query revoked gagal", err)
	}

	return nil
}

func (r *userSessionRepositoryImpl) RevokeAllSessionByUserID(ctx context.Context, userID string) error {
	const query = `UPDATE user_sessions SET revoked = true WHERE user_id = ?`
	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return apperror.New(apperror.CodeDBError, "revoke semua sesi gagal", err)
	}

	return nil
}
