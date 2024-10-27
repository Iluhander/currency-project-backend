package balance

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/Iluhander/currency-project-backend/internal/model"
	"github.com/Iluhander/currency-project-backend/internal/services/balance"
	"github.com/gin-gonic/gin"
)

type BalanceController struct {
	s *balance.BalanceService
}

type CurrencyChange struct {
	UserId model.TId `json:"userId"`
	Amount float64 `json:"amount"`
}

type CurrencyPricedChange struct {
	CurrencyChange
	Price float64 `json:"price"`
}

func Route(r *gin.RouterGroup, s *balance.BalanceService) (controller *BalanceController) {
	c := BalanceController{s}

	r.POST("/subtract", c.subtractCurrency)
	r.POST("/add", c.addCurrency)

	return &c
}

func (c *BalanceController) subtractCurrency(ctx *gin.Context) {
	jsonData, jsonErr := io.ReadAll(ctx.Request.Body)

	if jsonErr != nil {
		ctx.JSON(400, gin.H{
			"err": jsonErr,
		})

		return
	}

	var change CurrencyChange
	marshalErr := json.Unmarshal(jsonData, &change)

	if marshalErr != nil {
		ctx.JSON(400, gin.H{
			"err": marshalErr,
		})

		return
	}

	res, err := c.s.SubtractCurrency(change.UserId, change.Amount)
	if err != nil {
		if errors.Is(err, model.InvalidDataErr) {
			ctx.JSON(400, gin.H{
				"err": err,
			})
		} else {
			ctx.JSON(500, gin.H{
				"err": err,
			})
		}

		return;
	}

	ctx.JSON(200, gin.H{
		"balance": res,
	})
}

func (c *BalanceController) addCurrency(ctx *gin.Context) {
	jsonData, jsonErr := io.ReadAll(ctx.Request.Body)

	if jsonErr != nil {
		ctx.JSON(400, gin.H{
			"err": jsonErr,
		})

		return
	}

	var change CurrencyPricedChange
	marshalErr := json.Unmarshal(jsonData, &change)

	if marshalErr != nil {
		ctx.JSON(400, gin.H{
			"err": marshalErr,
		})

		return
	}

	link, err := c.s.AddCurrency(change.UserId, change.Amount, change.Price)
	if err != nil {
		if errors.Is(err, model.InvalidDataErr) {
			ctx.JSON(400, gin.H{
				"err": err,
			})
		} else {
			ctx.JSON(500, gin.H{
				"err": err,
			})
		}

		return;
	}

	ctx.JSON(200, gin.H{
		"link": link,
	})
}
