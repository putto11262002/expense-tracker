package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/configs"
	"github.com/putto11262002/expense-tracker/api/middlewares"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/routes"
	"github.com/putto11262002/expense-tracker/api/services"
	"gorm.io/gorm"
)

type App struct {
	httpServer        *http.Server
	userRepository    repositories.IUserRepository
	userService       services.IUserService
	groupRepository   repositories.IGroupRepository
	groupService      services.IGroupService
	expenseRepository repositories.IExpenseRepository
	expenseService    services.IExpenseService
	db                *gorm.DB
	config            configs.AppConfig
}

func NewApp(config configs.AppConfig) *App {

	db, error := configs.ConnectDB(config.DBConfig)
	if error != nil {
		log.Fatalf("error connecting to database: %v", error)
	}

	if err := configs.AutoMigrate(db); err != nil {
		log.Fatalf("error migrating database: %v", err)
	}

	userRepository := repositories.NewUserRepository(db)
	groupRepository := repositories.NewGroupRepository(db)
	expenseRepository := repositories.NewExpenseRepository(db)

	return &App{
		db:                db,
		userRepository:    userRepository,
		userService:       services.NewUserService(userRepository, config.JwtSecret),
		groupRepository:   groupRepository,
		groupService:      services.NewGroupService(groupRepository),
		expenseRepository: expenseRepository,
		expenseService:    services.NewExpenseService(expenseRepository),
		config:            config,
	}
}

func (app *App) Run() {
	router := gin.Default()

	if app.config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(middlewares.GlobalErrorHandler(), middlewares.CORSMiddleware())

	api := router.Group("/api")

	routes.NewExpenseRoutes(api, app.expenseService, app.userService)
	routes.NewUserRoutes(api, app.userService)
	routes.NewAuthRoutes(api, app.userService)
	routes.NewGroupRoutes(api, app.groupService, app.userService)

	// Check database connection
	api.GET("/health-check", func(ctx *gin.Context) {
		db, err := app.db.DB()
		if err != nil {
			ctx.Status(500)
			log.Printf("error retrieving database: %v", err)
			return
		}
		if err := db.Ping(); err != nil {
			log.Printf("error pinging database: %v", err)
			ctx.Status(500)
			return
		}

		ctx.Status(http.StatusOK)
	})

	app.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", app.config.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("Server running on port %d", app.config.Port)
		if err := app.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		db, err := app.db.DB()
		if err != nil {
			log.Fatalf("error retrieving database object: %v", err)

		}
		if err := db.Close(); err != nil {
			log.Fatalf("error closing database connection: %v", err)
		}
		log.Println("database connection closed")
		cancel()
	}()

	if err := app.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server:%+v", err)
	}
	log.Print("server shutdown gracefully")
}
