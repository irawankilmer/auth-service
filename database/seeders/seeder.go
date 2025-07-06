package seeders

import (
	"database/sql"
	"github.com/irawankilmer/auth-service/pkg/utils"
)

func SeedsRun(db *sql.DB, u utils.Utility) error {
	if err := Role(db, u); err != nil {
		return err
	}

	if err := User(db, u); err != nil {
		return err
	}

	return nil
}
