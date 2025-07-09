package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/dbtx"
	"github.com/irawankilmer/auth-service/internal/dto/response"
	"github.com/irawankilmer/auth-service/internal/model"
)

type UserRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error)
	FindUserByTokenVersion(ctx context.Context, userID string) (*model.UserModel, error)
	CheckUsername(ctx context.Context, username string) (bool, error)
	UsernameChange(ctx context.Context, user *response.UserDetailResponse, newUsername string) (bool, error)
	CheckEmail(ctx context.Context, email string) (bool, error)
	EmailChange(ctx context.Context, user *response.UserDetailResponse, newEmail string) (bool, error)
	Create(ctx context.Context, user *model.UserModel) error
	FindByID(ctx context.Context, userID string) (*response.UserDetailResponse, error)
	EmailUpdate(ctx context.Context, user *response.UserDetailResponse, newEmail string) error
	RoleUpdate(ctx context.Context, user *response.UserDetailResponse, newRoles []model.RoleModel) error
	UpdateEmailVerified(ctx context.Context, user *response.UserDetailResponse) error
	Delete(ctx context.Context, user *response.UserDetailResponse) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error) {
	const (
		queryTotal = `
			SELECT COUNT(*) FROM (
				SELECT u.id
				FROM users u
				JOIN user_roles ur ON u.id = ur.user_id
				JOIN roles r ON r.id = ur.role_id
				GROUP BY u.id
				HAVING SUM(CASE WHEN r.name IN ('admin', 'super admin') THEN 1 ELSE 0 END) = 0
			) AS filtered_users;
		`

		queryUsers = `
			SELECT 
				u.id, u.username, u.email,
				r.id AS role_id, r.name AS role_name,
				p.id AS profile_id, p.full_name
			FROM users u
			JOIN user_roles ur ON u.id = ur.user_id
			JOIN roles r ON r.id = ur.role_id
			LEFT JOIN profiles p ON p.user_id = u.id
			WHERE u.id IN (
				SELECT u.id
				FROM users u
				JOIN user_roles ur ON u.id = ur.user_id
				JOIN roles r ON r.id = ur.role_id
				GROUP BY u.id
				HAVING SUM(CASE WHEN r.name IN ('admin', 'super admin') THEN 1 ELSE 0 END) = 0
			)
			ORDER BY u.updated_at DESC
			LIMIT ? OFFSET ?;
		`
	)

	// Hitung total user kecuali admin dan super admin
	var total int
	if err := r.db.QueryRowContext(ctx, queryTotal).Scan(&total); err != nil {
		return nil, 0, apperror.New(apperror.CodeDBError, "gagal menghitung total users", err)
	}

	// Ambil data users
	rows, err := r.db.QueryContext(ctx, queryUsers, limit, offset)
	if err != nil {
		return nil, 0, apperror.New(apperror.CodeDBError, "gagal mengambil data users", err)
	}
	defer rows.Close()

	// Kelompokkan user dengan multiple role
	userMap := make(map[string]*response.UserResponse)
	var orderedIDs []string

	for rows.Next() {
		var (
			id, roleID, roleName, profileID string
			username                        sql.NullString
			email                           string
			fullName                        string
		)

		if err := rows.Scan(
			&id, &username, &email,
			&roleID, &roleName,
			&profileID, &fullName,
		); err != nil {
			return nil, 0, apperror.New(apperror.CodeDBError, "gagal scan data user", err)
		}

		user, exists := userMap[id]
		if !exists {
			// Buat instance baru
			user = &response.UserResponse{
				ID:    id,
				Email: email,
				Roles: []response.RoleResponse{},
				Profile: response.ProfileResponse{
					ID:       profileID,
					FullName: fullName,
				},
			}
			if username.Valid {
				user.Username = &username.String
			}

			userMap[id] = user
			orderedIDs = append(orderedIDs, id)
		}

		// Tambahkan role ke list
		user.Roles = append(user.Roles, response.RoleResponse{
			ID:   roleID,
			Name: roleName,
		})
	}

	// Konversi hasil map ke slice terurut
	var users []response.UserResponse
	for _, id := range orderedIDs {
		users = append(users, *userMap[id])
	}

	return users, total, nil
}

func (r *userRepository) FindUserByTokenVersion(ctx context.Context, userID string) (*model.UserModel, error) {
	const query = `SELECT token_version FROM users WHERE id = ?`
	var user model.UserModel

	// query
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.TokenVersion)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CheckUsername(ctx context.Context, username string) (bool, error) {
	const query = `SELECT exists(SELECT 1 FROM users WHERE username = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, apperror.New("[CODE_USERNAME_CHECK_INVALID]", "gagal memeriksa username", err, 500)
	}

	return exists, nil
}

func (r *userRepository) UsernameChange(ctx context.Context, user *response.UserDetailResponse, newUsername string) (bool, error) {
	const query = `SELECT exists(SELECT 1 FROM users WHERE username = ? AND id != ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, newUsername, user.ID).Scan(&exists)
	if err != nil {
		return false, apperror.New(apperror.CodeDBError, "cek username gagal", err)
	}

	return exists, nil
}

func (r *userRepository) CheckEmail(ctx context.Context, email string) (bool, error) {
	const query = `SELECT exists(SELECT 1 FROM users WHERE email = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, apperror.New("[CODE_EMAIL_CHECK_INVALID]", "gagal memeriksa email", err, 500)
	}

	return exists, nil
}

func (r *userRepository) EmailChange(ctx context.Context, user *response.UserDetailResponse, newEmail string) (bool, error) {
	const query = `SELECT exists(SELECT 1 FROM users WHERE email = ? AND id != ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, newEmail, user.ID).Scan(&exists)
	if err != nil {
		return false, apperror.New(apperror.CodeDBError, "cek email gagal", err)
	}

	return exists, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.UserModel) error {
	return dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		const (
			queryUser = `
									INSERT INTO 
									users(id, username, email, password, token_version, email_verified, created_by_admin, google_id)
									VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
			queryUserRoles = `INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)`
			queryProfile   = `
											INSERT INTO 
											profiles(id, user_id, full_name, address, gender, image)
											VALUES(?, ?, ?, ?, ?, ?)`
		)

		// create user
		_, err := tx.ExecContext(ctx, queryUser,
			user.ID, user.Username, user.Email, user.Password, user.TokenVersion,
			user.EmailVerified, user.CreatedByAdmin, user.GoogleID,
		)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "create user gagal", err)
		}

		// create profile
		_, err = tx.ExecContext(ctx, queryProfile,
			user.Profile.ID, user.Profile.UserID, user.Profile.FullName,
			user.Profile.Address, user.Profile.Gender, user.Profile.Image,
		)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "create profile gagal", err)
		}

		// create relation user_roles
		stmt, err := tx.PrepareContext(ctx, queryUserRoles)
		if err != nil {
			return apperror.New(apperror.CodeDBPrepareError, "prepare query user_roles gagal", err)
		}
		defer stmt.Close()

		for _, r := range user.Roles {

			if _, err := stmt.ExecContext(ctx, user.ID, r.ID); err != nil {
				return apperror.New(apperror.CodeDBError, "create relation user_roles gagal", err)
			}
		}
		return nil
	})
}

func (r *userRepository) FindByID(ctx context.Context, userID string) (*response.UserDetailResponse, error) {
	const (
		queryUserWithProfile = `
			SELECT 
				u.id, u.username, u.email, u.email_verified, u.created_by_admin, u.google_id,
				p.id AS profile_id, p.full_name, p.address, p.gender, p.image
			FROM users u
			LEFT JOIN profiles p ON p.user_id = u.id
			WHERE u.id = ?
			AND NOT EXISTS (
				SELECT 1 FROM user_roles ur
				JOIN roles r ON ur.role_id = r.id
				WHERE ur.user_id = u.id AND r.name IN (?, ?)
			);
		`

		queryUserRoles = `
			SELECT r.id, r.name
			FROM user_roles ur
			JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = ?;
		`
	)

	// Variabel scan
	var (
		id, email, profileID, fullName string
		username, googleID             sql.NullString
		emailVerified, createdByAdmin  bool
		address, gender, image         sql.NullString
	)

	// Ambil user + profile
	err := r.db.QueryRowContext(ctx, queryUserWithProfile, userID, "admin", "super admin").Scan(
		&id, &username, &email, &emailVerified, &createdByAdmin, &googleID,
		&profileID, &fullName, &address, &gender, &image,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.New(apperror.CodeUserNotFound, "user tidak ditemukan", err)
		}
		return nil, apperror.New(apperror.CodeDBError, "gagal mengambil data user", err)
	}

	// Bangun objek user
	user := &response.UserDetailResponse{
		ID:             id,
		Email:          email,
		EmailVerified:  emailVerified,
		CreatedByAdmin: createdByAdmin,
		Roles:          []response.RoleResponse{},
		Profile: response.ProfileDetailResponse{
			ID:       profileID,
			FullName: fullName,
		},
	}

	// Handle nullable fields
	if username.Valid {
		user.Username = &username.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
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

	// Ambil roles
	rows, err := r.db.QueryContext(ctx, queryUserRoles, userID)
	if err != nil {
		return nil, apperror.New(apperror.CodeDBError, "gagal mengambil data roles", err)
	}
	defer rows.Close()

	for rows.Next() {
		var role response.RoleResponse
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, apperror.New(apperror.CodeDBError, "gagal membaca data role", err)
		}
		user.Roles = append(user.Roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.New(apperror.CodeDBError, "error saat membaca rows role", err)
	}

	return user, nil
}

func (r *userRepository) EmailUpdate(ctx context.Context, user *response.UserDetailResponse, newEmail string) error {
	return dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		const (
			queryUser         = `UPDATE users SET email = ? WHERE id = ?`
			queryEmailHistory = `INSERT INTO email_history(email) VALUES(?)`
		)

		// update email
		_, err := tx.ExecContext(ctx, queryUser, newEmail, user.ID)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "update email gagal", err)
		}

		// insert email
		_, err = tx.ExecContext(ctx, queryEmailHistory, user.Email)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "insert email_history gagal", err)
		}
		return nil
	})
}

func (r *userRepository) RoleUpdate(ctx context.Context, user *response.UserDetailResponse, newRoles []model.RoleModel) error {
	return dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		const (
			queryDelete = `DELETE FROM user_roles WHERE user_id = ?`
			queryInsert = `INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)`
		)

		// hapus semua roles lama
		if _, err := tx.ExecContext(ctx, queryDelete, user.ID); err != nil {
			return apperror.New(apperror.CodeDBError, "delete roles gagal", err)
		}

		// memastikan roles baru memang ada
		if len(newRoles) == 0 {
			return apperror.New(apperror.CodeInvalidInput, "roles tidak boleh kosong", errors.New("roles tidak boleh kosong"))
		}

		// prepare roles baru
		stmt, err := tx.PrepareContext(ctx, queryInsert)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "prepare create roles gagal", err)
		}
		defer stmt.Close()

		// insert roles baru
		for _, r := range newRoles {
			if _, err := stmt.ExecContext(ctx, user.ID, r.ID); err != nil {
				return apperror.New(apperror.CodeDBError, "query create roles gagal", err)
			}
		}

		return nil
	})
}

func (r *userRepository) UpdateEmailVerified(ctx context.Context, user *response.UserDetailResponse) error {
	const query = `UPDATE users SET email_verified = true WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, query, user.ID); err != nil {
		return apperror.New(apperror.CodeDBError, "query update email verified gagal", err)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, user *response.UserDetailResponse) error {
	const (
		queryUsers           = `DELETE FROM users WHERE id = ?`
		queryUsernameHistory = `INSERT INTO username_history(username) VALUES(?)`
		queryEmailHistory    = `INSERT INTO email_history(email) VALUES(?)`
	)

	return dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, queryUsers, user.ID)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "delete user gagal", err)
		}

		_, err = tx.ExecContext(ctx, queryUsernameHistory, user.Username)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "insert username_history gagal", err)
		}

		_, err = tx.ExecContext(ctx, queryEmailHistory, user.Email)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "insert email_history gagal", err)
		}

		return nil
	})
}
