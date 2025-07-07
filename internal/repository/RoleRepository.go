package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/model"
	"strings"
)

type RoleRepository interface {
	CheckRoles(ctx context.Context, roles []string) ([]model.RoleModel, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

const maxRolesAllowed = 20

var (
	ErrEmptyRoles    = errors.New("roles tidak boleh kosong")
	ErrTooManyRoles  = errors.New("jumlah roles melebihi batas maksimal")
	ErrRolesNotFound = errors.New("satu atau lebih roles tidak ditemukan di database")
)

func (r *roleRepository) CheckRoles(ctx context.Context, roles []string) ([]model.RoleModel, error) {
	if len(roles) == 0 {
		return nil, apperror.New(apperror.CodeBadRequest, ErrEmptyRoles.Error(), ErrEmptyRoles)
	}

	if len(roles) > maxRolesAllowed {
		return nil, apperror.New(apperror.CodeBadRequest, ErrTooManyRoles.Error(), ErrTooManyRoles)
	}

	// Hilangkan duplikasi
	uniqueMap := make(map[string]struct{}, len(roles))
	uniqueRoles := make([]string, 0, len(roles))

	for _, role := range roles {
		if _, exists := uniqueMap[role]; !exists {
			uniqueMap[role] = struct{}{}
			uniqueRoles = append(uniqueRoles, role)
		}
	}

	// Siapkan placeholders dan args untuk query
	placeholders := strings.Repeat("?,", len(uniqueRoles))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]interface{}, len(uniqueRoles))
	for i, role := range uniqueRoles {
		args[i] = role
	}

	query := fmt.Sprintf(`SELECT id, name FROM roles WHERE name IN (%s)`, placeholders)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.New(apperror.CodeDBError, "gagal menjalankan query roles", err)
	}
	defer rows.Close()

	var foundRoles []model.RoleModel
	for rows.Next() {
		var role model.RoleModel
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, apperror.New(apperror.CodeDBError, "gagal membaca data role dari database", err)
		}
		foundRoles = append(foundRoles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.New(apperror.CodeDBError, "terjadi error saat iterasi hasil query roles", err)
	}

	if len(foundRoles) != len(uniqueRoles) {
		return nil, apperror.New(apperror.CodeRoleNotFound, ErrRolesNotFound.Error(), ErrRolesNotFound)
	}

	return foundRoles, nil
}
