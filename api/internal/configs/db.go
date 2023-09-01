package configs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int64
	Database string
}

func formatError(err error) (error){
	return fmt.Errorf("loading db config: %w", err)
}

func NewDBConfig(cp *ConfigParser)(*DBConfig, error) {
	

	username, err := cp.GetStringEnv("DB_USERNAME")
	if err != nil {
		return nil, formatError(err)
	}

	password, err := cp.GetStringEnv("DB_PASSWORD")
	if err != nil {
		return nil, formatError(err)
	}

	host, err := cp.GetStringEnv("DB_HOST")
	if err != nil {
		return nil , formatError(err)
	}

	port, err := cp.GetIntEnv("DB_PORT")
	if err != nil {
		return nil, formatError(err)
	}


	database, err := cp.GetStringEnv("DB_NAME")
	if err != nil {
		return nil , formatError(err)
	}

	return &DBConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}, nil
}

func ConnectDB(config DBConfig) (*gorm.DB, error) {
	dns := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local`,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	dbLogger := logger.Default.LogMode(logger.Silent)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return db, nil
}
