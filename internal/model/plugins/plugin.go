package plugins

import "github.com/Iluhander/currency-project-backend/internal/model"

const (
	TAuthPlugin = 1
	TPaymentPlugin = 2
	TStatisticsPlugin = 3
)

type Plugin struct {
	Id model.TId `json:"id"`
	Host string `json:"host"`
	Type int `json:"type"`
}