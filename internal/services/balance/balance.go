package balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Iluhander/currency-project-backend/internal/model"
	"github.com/Iluhander/currency-project-backend/internal/model/plugins"
	"github.com/Iluhander/currency-project-backend/internal/repository/users"
	pluginsServices "github.com/Iluhander/currency-project-backend/internal/services/plugins"
	usersServices "github.com/Iluhander/currency-project-backend/internal/services/users"
)

type BalanceService struct {
	usersService *usersServices.UsersService
	usersRepository *users.UsersRepository
	pluginsService *pluginsServices.PluginsService
}

type PayRequest struct {
	Amount float64 `json:"amount"`
	String string `json:"string"`
}

type PayResponse struct {
	Id string `json:"id"`
}

func Init(usersService *usersServices.UsersService, usersRepository *users.UsersRepository, pluginsService *pluginsServices.PluginsService) *BalanceService {
	return &BalanceService{
		usersService,
		usersRepository,
		pluginsService,
	}
}

func (b *BalanceService) SubtractCurrency(userId model.TId, amount float64) (float64, error) {
	return b.usersService.SubtractCurrency(userId, amount)	
}

func (b *BalanceService) AddCurrency(userId model.TId, amount, cost float64) (string, error) {
	curPlugins := b.pluginsService.GetPipeline("")

	var paymentPlugin *plugins.Plugin
	for _, v := range curPlugins {
		if v.Type == plugins.TPaymentPlugin {
			paymentPlugin = v

			break
		}
	}

	reqBody := PayRequest{
		Amount: amount,
		String: fmt.Sprintf("Purchase by user with id=%s of %f digital items", userId, amount),
	}

	bodyBytes, _ := json.Marshal(reqBody)
	body := bytes.NewBuffer(bodyBytes)

	req, formingError := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/api/payment", paymentPlugin.Host), body)

	if formingError != nil {
		return "", model.InternalErr
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to the payment plugin failed: %w", model.InternalErr)
	}

	defer res.Body.Close()

	var payRes PayResponse
	if err := json.NewDecoder(res.Body).Decode(&payRes); err != nil {
		return "", fmt.Errorf("request to the payment plugin resulted with incorrect formats: %w", model.InternalErr)
 	}

	orderCreationErr := b.usersRepository.CreateUserOrder(userId, payRes.Id, amount)
	if orderCreationErr != nil {
		return "", orderCreationErr
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/api/status/%s", paymentPlugin, payRes.Id))

	var cResp usersServices.OrderStatus
	if err == nil {
		if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
			return "", err
		}
	}

	defer resp.Body.Close()

	return cResp.Data.Operation[0].Link, nil
}
