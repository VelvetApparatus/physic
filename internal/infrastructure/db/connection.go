package db

import (
	"context"
	"database/sql"
	"fmt"
	"physk/internal/config"
	"time"
)

func NewConnection(ctx context.Context) (*sql.DB, error) {
	conf := config.C().Db
	database, err := sql.Open("postgres", stringConfig(conf))
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	database.SetMaxOpenConns(200)
	database.SetConnMaxLifetime(120 * time.Second)
	database.SetMaxIdleConns(30)
	database.SetConnMaxIdleTime(60 * time.Second)

	if err = database.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return database, nil
}

func stringConfig(c config.DBSettings) string {

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s application_name=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode, c.AppName)
}
