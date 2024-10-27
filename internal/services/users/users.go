package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"

	"github.com/Iluhander/currency-project-backend/internal/model"
	"github.com/Iluhander/currency-project-backend/internal/model/plugins"
	"github.com/Iluhander/currency-project-backend/internal/model/users"
	repository "github.com/Iluhander/currency-project-backend/internal/repository/users"
	pluginsServices "github.com/Iluhander/currency-project-backend/internal/services/plugins"
)

type UsersService struct {
	dbRepo *repository.UsersRepository
	conn *sql.DB

	pluginsService *pluginsServices.PluginsService
}

type GetUsers struct {
	Users []repository.OrderedUser
	Pagination model.PaginationDto
}

type PaymentLinkContainer struct {
	Status string `json:"status"`
	Link string `json:"paymentLink"`
}

type PaymentOpContainer struct {
	Operation []PaymentLinkContainer `json:"Operation"`
}

type OrderStatus struct {
	Data PaymentOpContainer `json:"Data"`
}

func Init(dbRepo *repository.UsersRepository, conn *sql.DB, pluginsService *pluginsServices.PluginsService) *UsersService  {
	return &UsersService{
		dbRepo,
		conn,
		pluginsService,
	}
}

func (s *UsersService) AddCurrency(userId model.TId, amount float64) (updatedBalance float64, resErr error) {
	if amount < 0 {
		return 0, fmt.Errorf("adding negative currency amount %f is prohibitedd: %w", amount, model.InvalidDataErr)
	}

	if amount == 0 {
		return 0, fmt.Errorf("adding zero currency is prohibitedd: %w", model.InvalidDataErr)
	}

	t, tErr := s.conn.Begin()

	if tErr != nil {
		return 0, tErr
	}

	createErr := s.dbRepo.TryCreateUser(t, userId)
	if createErr != nil {
		t.Rollback()

		return 0, createErr
	}

	changedVal, err := s.dbRepo.ChangeCurrency(t, userId, amount)

	if err != nil {
		t.Rollback()

		return 0, err
	}

	t.Commit()

	return changedVal, nil
}

func (s *UsersService) SubtractCurrency(userId model.TId, amount float64) (updatedBalance float64, resErr error) {
	if amount < 0 {
		return 0, fmt.Errorf("subtracting negative currency %f amount is prohibited: %w", amount, model.InvalidDataErr)
	}

	if amount == 0 {
		return 0, fmt.Errorf("subtracting zero currency is prohibited: %w", model.InvalidDataErr)
	}

	t, tErr := s.conn.Begin()

	if tErr != nil {
		return 0, tErr
	}

	createErr := s.dbRepo.TryCreateUser(t, userId)
	if createErr != nil {
		t.Rollback()

		return 0, createErr
	}

	changedVal, err := s.dbRepo.ChangeCurrency(t, userId, -1 * amount)

	if err != nil {
		t.Rollback()

		return 0, err
	}

	return changedVal, nil
}

func (s *UsersService) GetUserBalance(userId model.TId) (float64, error) {
	curBalance, curBalanceErr := s.dbRepo.GetOneBalance(userId)
	if curBalanceErr != nil {
		return 0, curBalanceErr
	}

	curOrders, curOrdersErr := s.dbRepo.GetUserOrders(userId)
	if curOrdersErr != nil {
		return 0, curOrdersErr
	}

	curPlugins := s.pluginsService.GetPipeline("")

	var paymentPlugin *plugins.Plugin
	for _, v := range curPlugins {
		if v.Type == plugins.TPaymentPlugin {
			paymentPlugin = v

			break
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(len(curOrders))

	t, _ := s.conn.Begin()

	sumBalance := curBalance
	for _, v := range curOrders {
		go func(v users.UserOrder) {
			resp, err := http.Get(fmt.Sprintf("http://%s/api/status/%s", paymentPlugin, v.OrderId))

			if err == nil {
				var cResp OrderStatus
				if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
				}

				if cResp.Data.Operation[0].Status == "APPROVED" {
					s.dbRepo.ChangeCurrency(t, v.UserId, v.Amount)
					s.dbRepo.RemoveUserOrder(v.Id)

					sumBalance += v.Amount
				}
			}

			defer resp.Body.Close()

			wg.Done()
		}(v)
	}

	wg.Wait()
	t.Commit()

	return sumBalance, nil
}

func (s *UsersService) GetUsers(offset, limit int, orderField, orderType string) (*GetUsers, error) {
	if orderField == "" {
		orderField = "balance"
	}
	
	if orderField != "balance" {
		return nil, fmt.Errorf("Field %s cannot be used for sorting: %w", orderField, model.InvalidDataErr)
	}

	if orderField != "balance" {
		return nil, fmt.Errorf("Field %s cannot be used for sorting: %w", orderField, model.InvalidDataErr)
	}
	
	users, usersErr := s.dbRepo.GetUsers(offset, limit, orderField, orderType)

	if usersErr != nil {
		return nil, usersErr
	}

	total, totalErr := s.dbRepo.CountUsers()

	if totalErr != nil {
		return nil, totalErr
	}

	return &GetUsers{
		users,
		model.PaginationDto{
			total,
			limit,
			int(math.Ceil(float64(total) / float64(limit))),
			offset / limit,
		},
	}, nil
}
