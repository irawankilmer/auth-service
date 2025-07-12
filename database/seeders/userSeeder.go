package seeders

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/dbtx"
	"github.com/irawankilmer/auth-service/pkg/utils"
	"net/http"
	"time"
)

func User(db *sql.DB, u utils.Utility) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return dbtx.WithTxContext(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		// Get data role super admin
		var roleID string
		err := tx.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = ?`, "super admin").Scan(&roleID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("role super admin tidak ditemukan: %w", err)
			}
			return fmt.Errorf("query role gagal: %w", err)
		}

		// generate token version
		tokenVersion, err := u.UUIDGenerate()
		if err != nil {
			return apperror.New("[UUID_GENERATE_VALIED]", "generate UUID gagal", err, http.StatusInternalServerError)
		}

		// Create user
		hashPass, err := u.HashGenerate("superadmin")
		if err != nil {
			return fmt.Errorf("generate hash password gagal: %w", err)
		}
		userID := u.ULIDGenerate()
		_, err = tx.ExecContext(ctx,
			`INSERT INTO users(id, username, email, password, email_verified, created_by_admin, token_version) VALUES(?, ?, ?, ?, ?, ?, ?)`,
			userID, "superadmin", "superadmin@gmail.com", hashPass, true, false, tokenVersion,
		)
		if err != nil {
			return fmt.Errorf("query users gagal: %w", err)
		}

		// create profiles
		_, err = tx.ExecContext(ctx, `INSERT INTO profiles(id, user_id, full_name, address, gender, image) VALUES(?, ?, ?, ?, ?, ?)`,
			u.ULIDGenerate(), userID, "Super Admin Pertama", "Samarang awi", 1, "default.jpg")
		if err != nil {
			return fmt.Errorf("create profiles gagal: %w", err)
		}

		// relasi user_roles
		_, err = tx.ExecContext(ctx, `INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)`, userID, roleID)
		if err != nil {
			return fmt.Errorf("create user_roles gagal: %w", err)
		}
		return nil
	})
}
