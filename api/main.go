package main

import (
	"fmt"
	"github.com/putto11262002/expense-tracker/api/configs"
	"github.com/putto11262002/expense-tracker/api/middlewares"
	"github.com/putto11262002/expense-tracker/api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	err := configs.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := configs.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := configs.AutoMigrate(db); err != nil {
		log.Fatal(err)
	}

	// loading PORT from environment
	port, err := configs.GetIntEnv("PORT")
	if err != nil {
		port = 3001
	}

	r := gin.Default()

	r.Use(middlewares.GlobalErrorHandler())

	routes.NewUserRoutes(db, r)
	routes.NewAuthRoutes(db, r)
	routes.NewGroupRoutes(db, r)

	routes.SetUpRoutes(r, db)

	err = r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(fmt.Errorf("starting server: %w", err))
	}
}
