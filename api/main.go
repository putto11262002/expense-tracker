package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/internal/configs"
	"github.com/putto11262002/expense-tracker/api/internal/routes"
)

func main() {
	cp := configs.NewConfigParser(".env")
	err := cp.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	dbConfig, err := configs.NewDBConfig(cp)
	if err != nil {
		log.Fatal(err)
	}

	db, err := configs.ConnectDB(*dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	// loading PORT from environment
	port, err := cp.GetIntEnv("PORT")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	routes.NewUserRoutes(db, r)

	r.Run(fmt.Sprintf(":%d", port))
}
