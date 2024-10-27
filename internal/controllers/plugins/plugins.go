package controllers

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/Iluhander/currency-project-backend/internal/model"
	pluginsModel "github.com/Iluhander/currency-project-backend/internal/model/plugins"
	"github.com/Iluhander/currency-project-backend/internal/services/plugins"
	"github.com/gin-gonic/gin"
)

type PluginsController struct {
	execService *plugins.ExecutionService 
}

func Route(r *gin.RouterGroup, execService *plugins.ExecutionService) (controller *PluginsController) {
	c := PluginsController{execService}

	r.GET("/", c.findPlugins)
	r.POST("/", c.add)
	r.PUT("/", c.update)
	r.DELETE("/:id", c.delete)

	return &c
}

func (c *PluginsController) findPlugins(ctx *gin.Context) {
	pluginId := ctx.Query("type")
	plugins := c.execService.GetPipeline(pluginId)

	ctx.JSON(200, gin.H{
		"plugins": plugins,
	})
}

func (c *PluginsController) add(ctx *gin.Context) {
	jsonData, jsonErr := io.ReadAll(ctx.Request.Body)
	if jsonErr != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Incorrect request body"})

		return
	}

	plugin := pluginsModel.Plugin{}
	unmarshallErr := json.Unmarshal(jsonData, &plugin)
	if unmarshallErr != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Incorrect request body"})

		return
	}

	addedPlugin, addErr := c.execService.AddPlugin(&plugin)
	if addErr != nil {
		if errors.Is(addErr, model.InvalidDataErr) {
			ctx.JSON(400, gin.H{
				"err": addErr,
			})
		} else {
			ctx.JSON(500, gin.H{
				"err": addErr,
			})
		}

		return;
	}

	ctx.JSON(201, gin.H{
		"plugin": addedPlugin,
	})
}

func (c *PluginsController) update(ctx *gin.Context) {
	jsonData, jsonErr := io.ReadAll(ctx.Request.Body)
	if jsonErr != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Incorrect request body"})

		return
	}

	plugin := pluginsModel.Plugin{}
	unmarshallErr := json.Unmarshal(jsonData, &plugin)
	if unmarshallErr != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Incorrect request body"})

		return
	}

	addedPlugin, addErr := c.execService.UpdatePlugin(&plugin)
	if addErr != nil {
		if errors.Is(addErr, model.InvalidDataErr) {
			ctx.JSON(400, gin.H{
				"err": addErr,
			})
		} else if errors.Is(addErr, model.NotFoundErr) {
			ctx.JSON(404, gin.H{
				"err": addErr,
			})
		} else {
			ctx.JSON(500, gin.H{
				"err": addErr,
			})
		}

		return;
	}

	ctx.JSON(200, gin.H{
		"plugin": addedPlugin,
	})
}

func (c *PluginsController) delete(ctx *gin.Context) {
	pluginId := ctx.Param("id")

	err := c.execService.DeletePlugin(pluginId)
	if err != nil {
		if errors.Is(err, model.InvalidDataErr) {
			ctx.JSON(400, gin.H{
				"err": err,
			})
		} else if errors.Is(err, model.NotFoundErr) {
			ctx.JSON(404, gin.H{
				"err": err,
			})
		} else {
			ctx.JSON(500, gin.H{
				"err": err,
			})
		}

		return;
	}

	ctx.JSON(200, gin.H{})
}
