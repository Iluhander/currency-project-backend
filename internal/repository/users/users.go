package users

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Iluhander/currency-project-backend/internal/config"
	"github.com/Iluhander/currency-project-backend/internal/model"
)

type UsersRepository struct {
	conn *sql.DB
}

func Init(cfg *config.ServiceConfig) (*UsersRepository, func(), error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	conn, err := sql.Open("postgres", psqlInfo)

	close := func() {
		conn.Close()
	}

	if err != nil {
		return nil, nil, err
	}

	return &UsersRepository{
		conn,
	}, close, nil
}

func (dbRepo *UsersRepository) ChangeCurrency(userId model.TId, amount float64) (updatedBalance int, resErr error) {
	_, err := dbRepo.conn.Query("UPDATE users SET balance = balance + amount")

	tx, err := dbRepo.conn.Begin()
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	_, err = tx.Exec("UPDATE users SET balance = balance + %d WHERE id=%s", amount, userId)
	var row *sql.Row
	row = tx.QueryRow("SELECT balance FROM users WHERE id=%s", userId)
	_, err = tx.Exec("COMMIT;")

	if err != nil {
		return 0, err
	}

	scanErr := row.Scan(&updatedBalance)
	if scanErr != nil {
		return 0, scanErr
	}

	return updatedBalance, nil
}

func (dbRepo *UsersRepository) GetOneBalance(userId model.TId) (float64, error) {
	rows, err := dbRepo.conn.Query(fmt.Sprintf("SELECT amount FROM users WHERE id=%s", userId))

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var balance float64
		readErr := rows.Scan(&balance)

		if readErr != nil {
			return 0, readErr
		}

		return balance, nil
	}

	return 0, fmt.Errorf("Missing user with id=%s", userId)
}

func (dbRepo *UsersRepository) GetBalances(offset, limit int, orderType string) ([]model.User, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	
	rows, err := dbRepo.conn.Query("SELECT id, balance, ROW_NUMBER() OVER (ORDER BY balance %s) as idx FROM users ORDER BY balance %s OFFSET %d LIMIT %d", orderType, orderType, offset, limit)

	if err != nil {
		return nil, err
	}

	users := make([]model.User, limit);
	for rows.Next() {
		curUser := model.User{}
		readErr := rows.Scan(&curUser.Id, &curUser.Balance)

		if readErr != nil {
			return nil, readErr
		}

		users = append(users, curUser)
	}

	return users, nil
}

