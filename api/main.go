package main

import (

	"log"


	"github.com/putto11262002/expense-tracker/api/configs"
	"github.com/putto11262002/expense-tracker/api/server"

)

func main() {

	if configs.GetGoEnv() != "production" {
		err := configs.LoadEnv(".env")
		if err != nil {
			log.Fatal(err)
		}
	}

	port, err := configs.GetIntEnv("PORT")
	if err != nil {
		port = 3001
	}

	username, err := configs.GetStringEnv("DB_USERNAME")
	if err != nil {
		log.Fatalf("loading database config: %v", err)
	}

	password, err := configs.GetStringEnv("DB_PASSWORD")
	if err != nil {
		log.Fatalf("loading database config: %v", err)
	}

	host, err := configs.GetStringEnv("DB_HOST")
	if err != nil {
		log.Fatalf("loading database config: %v", err)
	}

	dbPort, err := configs.GetIntEnv("DB_PORT")
	if err != nil {
		log.Fatalf("loading database config: %v", err)
	}

	database, err := configs.GetStringEnv("DB_NAME")
	if err != nil {
		log.Fatalf("loading database config: %v", err)
	}

	jwtSecret, err := configs.GetStringEnv("JWT_SECRET")
	if err != nil {
		jwtSecret = "secret"
	}

	config := configs.AppConfig{
		Port: int(port),
		DBConfig: configs.DBConfig{
			Username: username,
			Password: password,
			Host: host,
			Port: dbPort,
			Database: database,
		},
		JwtSecret: jwtSecret,
		Env: configs.GetGoEnv(),
	}

	app := server.NewApp(config)
	app.Run()
}
