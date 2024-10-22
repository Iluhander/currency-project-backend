package users

import (
	"fmt"

	"github.com/Iluhander/currency-project-backend/internal/model"
	repository "github.com/Iluhander/currency-project-backend/internal/repository/users"
)

type UsersService struct {
	dbRepo *repository.UsersRepository
}

func Init(dbRepo *repository.UsersRepository) *UsersService  {
	return &UsersService{
		dbRepo,
	}
}

func (s *UsersService) AddCurrency(userId model.TId, amount float64) (updatedBalance int, resErr error) {
	if amount < 0 {
		return 0, fmt.Errorf("adding negative currency amount %f is prohibited", amount)
	}

	if amount == 0 {
		return 0, fmt.Errorf("adding zero currency is prohibited")
	}

	return s.dbRepo.ChangeCurrency(userId, amount)
}

func (s *UsersService) SubtractCurrency(userId model.TId, amount float64) (updatedBalance int, resErr error) {
	if amount < 0 {
		return 0, fmt.Errorf("subtracting negative currency %f amount is prohibited", amount)
	}

	if amount == 0 {
		return 0, fmt.Errorf("subtracting zero currency is prohibited")
	}

	return s.dbRepo.ChangeCurrency(userId, -1 * amount)
}

func (s *UsersService) GetUserBalance(userId model.TId) (float64, error) {
	return s.dbRepo.GetOneBalance(userId)
}

func (s *UsersService) GetUsersBalances(offset, limit int, orderType string) ([]model.User, error) {
	return s.dbRepo.GetBalances(offset, limit, orderType)
}
