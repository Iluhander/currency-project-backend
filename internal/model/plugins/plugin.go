package plugins

import "github.com/Iluhander/currency-project-backend/internal/model"

const (
	TAuthPlugin = "AUTHENTICATION"
	TPaymentPlugin = "PAYMENT"
	TStatisticsPlugin = "STATISTIC"
)

type Plugin struct {
	Id model.TId `json:"id"`
	Host string `json:"host"`
	Type string `json:"type"`
}