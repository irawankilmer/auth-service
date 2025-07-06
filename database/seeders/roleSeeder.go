package seeders

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/dbtx"
	"github.com/irawankilmer/auth-service/pkg/utils"
	"time"
)

func Role(db *sql.DB, u utils.Utility) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return dbtx.WithTxContext(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		query := `INSERT INTO roles(id, name) VALUES(?, ?)`
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return fmt.Errorf("prepare role gagal: %w", err)
		}
		defer stmt.Close()

		roles := []string{"super admin", "admin", "editor", "penulis", "tamu"}
		for _, r := range roles {
			_, err := stmt.ExecContext(ctx, u.ULIDGenerate(), r)
			if err != nil {
				return fmt.Errorf("query inser roles gagal: %w", err)
			}
		}
		return nil
	})
}
