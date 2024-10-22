package model

const (
	TAuthPlugin = 1
	TPaymentPlugin = 2
	TStatisticsPlugin = 3
)

type Plugin struct {
	Id TId
	Host string
	Type int
}