package controllers

import (
	"github.com/Iluhander/currency-project-backend/internal/services"
	"github.com/Iluhander/currency-project-backend/internal/services/users"
	"github.com/gin-gonic/gin"
)

type S3Controller struct {
	usersService *users.UsersService
	execService *services.ExecutionService 
}

func Route(r *gin.RouterGroup, usersService *users.UsersService, execService *services.ExecutionService) (controller *S3Controller) {
	c := S3Controller{usersService, execService}

	pluginsGroup := r.Group("/plugins")
	pluginsGroup.GET("/", c.getPlugins)

	return &c
}

func (c *S3Controller) getPlugins(ctx *gin.Context) {
	pipeline := c.execService.GetPipeline(-1)

	ctx.JSON(200, gin.H{
		"url": pipeline,
	})
}
