package migrations

import (
	"database/sql"
)

type CommonMigrations struct {
	conn *sql.DB
}

func Init(conn *sql.DB) *CommonMigrations {
	return &CommonMigrations{conn}
}

func (m *CommonMigrations) Run() error {
	_, err := m.conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	return err
}
