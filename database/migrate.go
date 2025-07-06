package database

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{MigrationsTable: "auth_migrations"})
	if err != nil {
		return fmt.Errorf("gagal membuat driver: %w", err)
	}

	sourceDriver, err := iofs.New(MigrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("gagal load file migrasi: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "mysql", driver)
	if err != nil {
		return fmt.Errorf("gagal inisialisasi migrasi: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("gagal cek versi migrasi: %w", err)
	}

	if dirty {
		fmt.Printf("DB Dirty di versi %d. Memperbaiki...\n", version)
		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("gagal force versi %d: %w", version, err)
		}
	}

	// Jalankan DOWN
	fmt.Println("Menjalankan Down...")
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("gagal saat down...:%w", err)
	}

	// Jalankan UP
	fmt.Println("Menjalankan UP...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("gagal saat up...:%w", err)
	}

	fmt.Println("migrasi berhasil...")
	return nil
}
