package users

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Iluhander/currency-project-backend/internal/model"
	"github.com/Iluhander/currency-project-backend/internal/model/users"
)

type UsersRepository struct {
	conn *sql.DB
}

type OrderedUser struct {
	Idx int `json:"index"`
	users.User
}

func Init(conn *sql.DB) *UsersRepository {
	return &UsersRepository{
		conn,
	}
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

func (dbRepo *UsersRepository) GetUsers(offset, limit int, orderField, orderType string) ([]OrderedUser, error) {
	if orderType != model.TSortAsc && orderType != model.TSortDesc {
		return make([]OrderedUser, 0), fmt.Errorf("%s is not a valid order type", orderType)
	}
	
	var wg sync.WaitGroup
	wg.Add(2)

	rows, err := dbRepo.conn.Query("SELECT id, balance, ROW_NUMBER() OVER (ORDER BY " +
		orderField + " " + orderType +
		") as idx FROM users ORDER BY " +
		orderField + " " + orderType +
		" OFFSET $1 LIMIT $2",
		offset,
		limit)

	if err != nil {
		return nil, err
	}

	users := make([]OrderedUser, 0, limit);
	for rows.Next() {
		curUser := OrderedUser{}
		readErr := rows.Scan(&curUser.Id, &curUser.Balance, &curUser.Idx)

		if readErr != nil {
			return nil, readErr
		}

		if curUser.Id == "" {
			return users, nil
		}

		users = append(users, curUser)
	}

	return users, nil
}

