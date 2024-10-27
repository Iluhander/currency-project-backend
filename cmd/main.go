package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/Iluhander/currency-project-backend/internal/config"
	balanceControllers "github.com/Iluhander/currency-project-backend/internal/controllers/balance"
	pluginsControllers "github.com/Iluhander/currency-project-backend/internal/controllers/plugins"
	usersControllers "github.com/Iluhander/currency-project-backend/internal/controllers/users"
	"github.com/Iluhander/currency-project-backend/internal/migrations"
	usersMigrations "github.com/Iluhander/currency-project-backend/internal/migrations/users"
	"github.com/Iluhander/currency-project-backend/internal/repository/pipelines"
	"github.com/Iluhander/currency-project-backend/internal/repository/users"
	"github.com/Iluhander/currency-project-backend/internal/services/balance"
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

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	} else {
		log.Println("Opened a db connection")
	}

	defer func(conn *sql.DB) {
		log.Println("Closing the db connection")
		conn.Close()
	}(conn)

	migrationErr := migrations.Init(conn).Run()
	if migrationErr != nil {
		panic(migrationErr)
	} else {
		log.Println("Basic migrations executed successfully")
	}

	usersMigrationsErr := usersMigrations.Init(conn).Run()
	if usersMigrationsErr != nil {
		panic(usersMigrationsErr)
	} else {
		log.Println("Users migrations executed successfully")
	}
	
	dbRepo := users.Init(conn)
	if err != nil {
		panic(err)
	} else {
		log.Printf("Connected to db %s:%d\n", cfg.DBHost, cfg.DBPort)
	}

	pipeRepo, err := pipelines.Init("pipeline.json")
	if err != nil {
		panic(err)
	}

	pluginsService := pluginsService.Init(pipeRepo)
	userService := usersService.Init(dbRepo, conn, pluginsService)
	balanceService := balance.Init(userService, dbRepo, pluginsService)

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
	pluginsControllers.Route(r.Group("plugins"), pluginsService)
	balanceControllers.Route(r.Group("balance"), balanceService)

	listenTo := fmt.Sprint(cfg.ServeEnpoint, ":", strconv.Itoa(int(cfg.ServePort)))
	log.Println("Listening to", listenTo)
	r.Run(listenTo)
}


