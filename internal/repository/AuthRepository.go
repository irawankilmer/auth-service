package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/dbtx"
	"github.com/irawankilmer/auth-service/internal/dto/response"
	"github.com/irawankilmer/auth-service/internal/model"
	"net/http"
)

type AuthRepository interface {
	IdentifierCheck(ctx context.Context, identifier string) (*model.UserModel, error)
	UpdateTokenVersion(ctx context.Context, userID, newTokenVersion string) error
	Me(ctx context.Context, userID string) (*response.UserDetailResponse, error)
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

func (r *authRepository) Me(ctx context.Context, userID string) (*response.UserDetailResponse, error) {
	const (
		query = `
							SELECT 
								u.id, u.username, u.email, u.email_verified,
								p.id, p.full_name, p.address, p.gender, p.image
							FROM users u 
							INNER JOIN profiles p ON u.id = p.user_id
							WHERE u.id = ?
						`
		queryRoles = `
									SELECT
										r.id, r.name
									FROM roles r
									INNER JOIN user_roles ur ON ur.role_id = r.id
									WHERE ur.user_id = ?
									`
	)
	user := &response.UserDetailResponse{
		Profile: response.ProfileDetailResponse{},
	}
	var roles []response.RoleResponse

	err := dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		// query user dan profiles
		var username, fullName, address, gender, image sql.NullString
		if err := tx.QueryRowContext(ctx, query, userID).Scan(
			&user.ID, &username, &user.Email, &user.EmailVerified, &user.Profile.ID, &fullName, &address, &gender, &image,
		); err != nil {
			if err == sql.ErrNoRows {
				return apperror.New("[USER_NOT_FOUND]", "user tidak ditemukan", err, http.StatusUnauthorized)
			}
			return apperror.New(apperror.CodeDBError, "query users gagal", err)
		}
		if username.Valid {
			user.Username = &username.String
		}
		if fullName.Valid {
			user.Profile.FullName = &fullName.String
		}
		if address.Valid {
			user.Profile.Address = &address.String
		}
		if gender.Valid {
			user.Profile.Gender = &gender.String
		}
		if image.Valid {
			user.Profile.Image = &image.String
		}

		// query roles
		rows, err := tx.QueryContext(ctx, queryRoles, userID)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "query roles gagal", err)
		}
		defer rows.Close()

		// scan roles
		for rows.Next() {
			var role response.RoleResponse
			if err := rows.Scan(&role.ID, &role.Name); err != nil {
				return apperror.New(apperror.CodeDBError, "scan query roles gagal", err)
			}
			roles = append(roles, role)
		}

		// cek error rows roles
		if err := rows.Err(); err != nil {
			return apperror.New(apperror.CodeDBError, "gagal setelah iterasi roles", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	user.Roles = roles

	return user, nil
}
