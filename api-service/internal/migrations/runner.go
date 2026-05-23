package migrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

func Run(ctx context.Context, db *sqlx.DB, dir string) error {
	if strings.TrimSpace(dir) == "" {
		dir = "db/migrations"
	}

	if err := ensureTable(ctx, db); err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		version := migrationVersion(file)

		applied, err := isApplied(ctx, db, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", filepath.Base(file), err)
		}

		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(version) VALUES ($1)`, version); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func ensureTable(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version TEXT PRIMARY KEY,
            applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
        );
    `)
	return err
}

func isApplied(ctx context.Context, db *sqlx.DB, version string) (bool, error) {
	var exists bool
	err := db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, version)
	return exists, err
}

func migrationVersion(file string) string {
	base := filepath.Base(file)
	return strings.TrimSuffix(base, ".up.sql")
}
