package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/logger"
)

type SQLite struct {
	logger logger.ILogger
}

type ISQLite interface {
	Connect(config Config) (*sql.DB, error)
	ConnectGorm(config Config) (*gorm.DB, error)
}

func NewSQLite(logger logger.ILogger) ISQLite {
	return &SQLite{logger}
}

func (c *SQLite) Connect(config Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.DataSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	c.logger.LogIt("DEBUG", "Database sqlite connected")

	return db, nil
}

func (c *SQLite) ConnectGorm(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DataSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	c.logger.LogIt("DEBUG", "Database sqlite connected")

	return db, nil
}
