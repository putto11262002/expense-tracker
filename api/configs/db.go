package configs

import (
	"fmt"
	"github.com/putto11262002/expense-tracker/api/domains"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int64
	Database string
}

func formatError(err error) error {
	return fmt.Errorf("loading db config: %w", err)
}

func loadDBConfig() (*DBConfig, error) {

	username, err := GetStringEnv("DB_USERNAME")
	if err != nil {
		return nil, formatError(err)
	}

	password, err := GetStringEnv("DB_PASSWORD")
	if err != nil {
		return nil, formatError(err)
	}

	host, err := GetStringEnv("DB_HOST")
	if err != nil {
		return nil, formatError(err)
	}

	port, err := GetIntEnv("DB_PORT")
	if err != nil {
		return nil, formatError(err)
	}

	database, err := GetStringEnv("DB_NAME")
	if err != nil {
		return nil, formatError(err)
	}

	return &DBConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}, nil
}

func ConnectDB() (*gorm.DB, error) {
	config, err := loadDBConfig()
	if err != nil {
		return nil, err
	}
	dns := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC`,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	// dbLogger := logger.Default.LogMode(logger.Silent)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		// Logger: dbLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domains.User{}, &domains.Group{}, &domains.Expense{}, &domains.Split{}); err != nil {
		return fmt.Errorf("auto migration: %w", err)
	}
	return nil
}
