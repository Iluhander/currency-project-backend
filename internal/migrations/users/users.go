package users

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
	_, createUsersErr := m.conn.Exec("CREATE TABLE IF NOT EXISTS public.users (\n" +
	"id uuid NOT NULL DEFAULT uuid_generate_v4(),\n" +
	"balance numeric NULL DEFAULT 0,\n" +
	"CONSTRAINT users_pkey PRIMARY KEY (id)\n" +
	")",
	)

	if createUsersErr != nil {
		return createUsersErr
	}

	_, createUsersOrderErr := m.conn.Exec("CREATE TABLE IF NOT EXISTS public.users_orders (\n" +
		"id uuid NOT NULL DEFAULT uuid_generate_v4(),\n" +
		"user_id uuid NOT NULL,\n" +
		"order_id uuid NOT NULL,\n" +
		"amount numeric NULL DEFAULT 0,\n" +
		"CONSTRAINT users_orders_pkey PRIMARY KEY (id),\n" +
		"CONSTRAINT users_orders_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id)\n" +
		")",
	)

	return createUsersOrderErr
}
