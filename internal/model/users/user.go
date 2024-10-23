package users

import "github.com/Iluhander/currency-project-backend/internal/model"

type User struct {
	Id model.TId `json:"id"`
	Balance float64 `json:"balance"`
}
