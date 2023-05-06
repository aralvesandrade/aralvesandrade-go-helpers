package main

import (
	"fmt"
	"os"
	"strconv"

	dbMySQL "github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/db"
	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/health"
	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/health/db"
	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/health/url"
	"github.com/aralvesandrade/aralvesandrade-go-helpers/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logger.NewLogger()
	mysqlService := dbMySQL.NewMySQL(logger)

	mysqlConfig := dbMySQL.Config{
		Port:     os.Getenv("MYSQL_PORT"),
		Host:     os.Getenv("MYSQL_HOST"),
		Database: os.Getenv("MYSQL_DATABASE"),
		Username: os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	}

	//database, err := mysqlService.ConnectByDSN(os.Getenv("MYSQL_URL"))
	database, err := mysqlService.Connect(mysqlConfig)

	if err != nil {
		logMsg := fmt.Sprintf("Error on connecting to mysql database: %v", err)
		logger.LogIt("ERROR", logMsg)
		os.Exit(1)
	}
	defer database.Close()

	handler := health.NewHandler()
	mysql := db.NewMySQLChecker(database)

	handler.AddChecker("Go", url.NewChecker("https://golang.org/"))
	handler.AddChecker("MySQL", mysql)

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
