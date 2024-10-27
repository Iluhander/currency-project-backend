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

func (dbRepo *UsersRepository) TryCreateUser(tx *sql.Tx, userId model.TId) error {
	_, err := tx.Exec(fmt.Sprintf("insert into users (id, balance) values (uuid('%s'), 0) on conflict do nothing", userId))

	return err
}

func (dbRepo *UsersRepository) ChangeCurrency(tx *sql.Tx, userId model.TId, amount float64) (updatedBalance float64, resErr error) {
	_, err := tx.Exec(fmt.Sprintf("UPDATE users SET balance = balance + %f WHERE id=uuid('%s')", amount, userId))
	if err != nil {
		return 0, err
	}

	var row *sql.Row
	row = tx.QueryRow(fmt.Sprintf("SELECT balance FROM users WHERE id=uuid('%s')", userId))

	scanErr := row.Scan(&updatedBalance)
	if scanErr != nil {
		return 0, scanErr
	}

	return updatedBalance, nil
}

func (dbRepo *UsersRepository) GetOneBalance(userId model.TId) (float64, error) {
	rows, err := dbRepo.conn.Query(fmt.Sprintf("SELECT amount FROM users WHERE id=uuid('%s')", userId))

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

func (dbRepo *UsersRepository) CountUsers() (int, error) {
	rows, err := dbRepo.conn.Query("SELECT COUNT(*) as total FROM users")

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var count int
		readErr := rows.Scan(&count)

		if readErr != nil {
			return 0, readErr
		}

		return count, nil
	}

	return 0, nil
}

func (dbRepo *UsersRepository) CreateUserOrder (userId, orderId model.TId, amount float64) error {
	_, err := dbRepo.conn.Exec(fmt.Sprintf("INSERT INTO users_orders (user_id, order_id, amount) values (uuid('%s'), uuid('%s'), %f)", userId, orderId, amount))

	return err
}

func (dbRepo *UsersRepository) RemoveUserOrder (userOrderRecordId model.TId) error {
	_, err := dbRepo.conn.Exec(fmt.Sprintf("DELETE FROM users_orders WHERE id=uuid('%s')", userOrderRecordId))

	return err
}

func (dbRepo *UsersRepository) GetUserOrders (userId model.TId) (orders []users.UserOrder, err error) {
	rows, err := dbRepo.conn.Query(fmt.Sprintf("SELECT id, order_id as orderId, amount, user_id as userId FROM users_orders WHERE user_id=uuid('%s')", userId))
	if err != nil {
		return make([]users.UserOrder, 0), err
	}
	
	orders = make([]users.UserOrder, 0)
	for rows.Next() {
		var order users.UserOrder
		readErr := rows.Scan(&order.Id, &order.OrderId, &order.Amount, &order.UserId)

		if readErr != nil {
			return make([]users.UserOrder, 0), readErr
		}
	}

	return orders, nil
}

