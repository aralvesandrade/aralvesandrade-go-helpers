package main

import (
	"fmt"
	"os"
	"strconv"

	database "github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/db"
	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/health"
	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/health/db"
	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/health/url"
	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logger.NewLogger()
	mysqlService := database.NewMySQL(logger)
	sqliteService := database.NewSQLite(logger)

	config := database.Config{
		Port:       os.Getenv("MYSQL_PORT"),
		Host:       os.Getenv("MYSQL_HOST"),
		Database:   os.Getenv("MYSQL_DATABASE"),
		Username:   os.Getenv("MYSQL_USER"),
		Password:   os.Getenv("MYSQL_PASSWORD"),
		DataSource: os.Getenv("SQLITE_DATASOURCE"),
	}

	mysqlDB, err := mysqlService.Connect(config)
	if err != nil {
		logMsg := fmt.Sprintf("Error on connecting to mysql: %v", err)
		logger.LogIt("ERROR", logMsg)
		os.Exit(1)
	}
	defer mysqlDB.Close()

	mysqlDB1, err := mysqlService.ConnectByDSN(os.Getenv("MYSQL_URL"))
	if err != nil {
		logMsg := fmt.Sprintf("Error on connecting to mysql: %v", err)
		logger.LogIt("ERROR", logMsg)
		os.Exit(1)
	}
	defer mysqlDB1.Close()

	mysqlDB2, err := mysqlService.ConnectGorm(config)
	if err != nil {
		logMsg := fmt.Sprintf("Error on connecting to mysql: %v", err)
		logger.LogIt("ERROR", logMsg)
		os.Exit(1)
	}
	mysqlDB2Aux, _ := mysqlDB2.DB()
	defer mysqlDB2Aux.Close()

	mysqlDB3, err := mysqlService.ConnectByDSNGorm(os.Getenv("MYSQL_URL"))
	if err != nil {
		logMsg := fmt.Sprintf("Error on connecting to mysql: %v", err)
		logger.LogIt("ERROR", logMsg)
		os.Exit(1)
	}
	mysqlDB3Aux, _ := mysqlDB3.DB()
	defer mysqlDB3Aux.Close()

	sqliteDB, err := sqliteService.Connect(config)
	if err != nil {
		logMsg := fmt.Sprintf("Error on connecting to sqlite: %v", err)
		logger.LogIt("ERROR", logMsg)
		os.Exit(1)
	}
	defer sqliteDB.Close()

	sqliteDB1, err := sqliteService.ConnectGorm(config)
	if err != nil {
		logMsg := fmt.Sprintf("Error on connecting to sqlite: %v", err)
		logger.LogIt("ERROR", logMsg)
		os.Exit(1)
	}
	sqliteDB1Aux, _ := sqliteDB1.DB()
	defer sqliteDB1Aux.Close()

	handler := health.NewHandler()
	mysql := db.NewMySQLChecker(mysqlDB)
	sqlite := db.NewSqlite3Checker(sqliteDB)

	handler.AddChecker("Go", url.NewChecker("https://golang.org/"))
	handler.AddChecker("MySQL", mysql)
	handler.AddChecker("SQLite", sqlite)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	logger.LogIt("INFO", fmt.Sprintf("Starting on port: %d", port))

	router := gin.Default()
	router.GET("/health", handler.Health)

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
