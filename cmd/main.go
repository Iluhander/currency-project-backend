package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/Iluhander/currency-project-backend/internal/config"
	pluginsControllers "github.com/Iluhander/currency-project-backend/internal/controllers/plugins"
	usersControllers "github.com/Iluhander/currency-project-backend/internal/controllers/users"
	"github.com/Iluhander/currency-project-backend/internal/repository/pipelines"
	"github.com/Iluhander/currency-project-backend/internal/repository/users"
	pluginsService "github.com/Iluhander/currency-project-backend/internal/services/plugins"
	usersService "github.com/Iluhander/currency-project-backend/internal/services/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dev := flag.Bool("dev", false, "Is in the dev mode")
	flag.Parse()

	prod := !*dev

	cfg, err := config.Init(prod)
	if err != nil {
		panic(err)
	}
	
	dbRepo, closeCallback, err := users.Init(cfg)
	if err != nil {
		panic(err)
	}

	defer closeCallback()


	pipeRepo, err := pipelines.Init("pipeline.json")
	if err != nil {
		panic(err)
	}

	userService := usersService.Init(dbRepo)
	executionService := pluginsService.Init(pipeRepo)

	if prod {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setting up the router.
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
	}))

	usersControllers.Route(r.Group("users"), userService)
	pluginsControllers.Route(r.Group("plugins"), executionService)

	r.Run(fmt.Sprint(":", strconv.Itoa(int(cfg.ServePort))))
}


