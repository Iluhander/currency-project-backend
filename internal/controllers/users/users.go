package controllers

import (
	"errors"
	"strconv"

	"github.com/Iluhander/currency-project-backend/internal/model"
	"github.com/Iluhander/currency-project-backend/internal/services/users"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	service *users.UsersService
}

func Route(r *gin.RouterGroup, usersService *users.UsersService) (controller *UsersController) {
	c := UsersController{usersService}

	r.GET("/", c.findUsers)

	return &c
}

func (c *UsersController) findUsers(ctx *gin.Context) {
	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	if pageSize <= 0 {
		pageSize = 10
	}

	sortField := ctx.Query("sortField")
	sort := ctx.Query("sort")
	if sort == "" {
		sort = model.TSortDesc
	}
	
	users, err := c.service.GetUsers(pageNumber * pageSize, pageSize, sortField, sort)
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
		"users": users,
	})
}
