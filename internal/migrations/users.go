package migrations

import (
	"database/sql"
)

type UsersMigrations struct {
	conn *sql.DB
}

func Init(conn *sql.DB) *UsersMigrations {
	return &UsersMigrations{conn}
}

func (m *UsersMigrations) Run() error {
	_, err := m.conn.Exec("CREATE TABLE IF NOT EXISTS public.users (\n" +
	"id uuid NOT NULL DEFAULT uuid_generate_v4(),\n" +
	"balance numeric NULL DEFAULT 0,\n" +
	"CONSTRAINT users_pkey PRIMARY KEY (id)\n" +
	")",
	)

	return err
}
