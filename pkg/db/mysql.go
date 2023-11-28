package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/logger"
)

const defaultMySQLPort = 3306

type Config struct {
	Port       string
	Host       string
	Database   string
	Username   string
	Password   string
	TimeoutSec int
	DataSource string
}

type MySQL struct {
	logger logger.ILogger
}

type IMySQL interface {
	Connect(config Config) (*sql.DB, error)
	ConnectByDSN(dsn string) (*sql.DB, error)
	ConnectGorm(config Config) (*gorm.DB, error)
	ConnectByDSNGorm(dsn string) (*gorm.DB, error)
}

func NewMySQL(logger logger.ILogger) IMySQL {
	return &MySQL{logger}
}

func (c *MySQL) Connect(config Config) (*sql.DB, error) {
	dsn := config.toMySQLConnectionString()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if config.TimeoutSec == 0 {
		config.TimeoutSec = 10
	}

	db.SetConnMaxLifetime(time.Duration(config.TimeoutSec) * time.Second)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	c.logger.LogIt("DEBUG", "Database mysql connected")

	return db, nil
}

func (c *MySQL) ConnectByDSN(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	c.logger.LogIt("DEBUG", "Database mysql connected")

	return db, nil
}

func (c *MySQL) ConnectGorm(config Config) (*gorm.DB, error) {
	dsn := config.toMySQLConnectionString()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.TimeoutSec == 0 {
		config.TimeoutSec = 10
	}

	sqlDB.SetConnMaxLifetime(time.Duration(config.TimeoutSec) * time.Second)

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	c.logger.LogIt("DEBUG", "Database mysql connected")

	return db, nil
}

func (c *MySQL) ConnectByDSNGorm(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

	c.logger.LogIt("DEBUG", "Database mysql connected")

	return db, nil
}

func (c *Config) toMySQLConnectionString() string {
	port, _ := strconv.Atoi(c.Port)

	if port == 0 {
		port = defaultMySQLPort
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		port,
		c.Database,
	)
}
