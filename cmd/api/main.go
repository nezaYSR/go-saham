package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gitlab.com/nezaysr/go-saham.git/config"
	"gitlab.com/nezaysr/go-saham.git/storage"
)

type Config struct {
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portRaw := os.Getenv("PORT")
	port, err := strconv.ParseInt(portRaw, 0, 0)
	if err != nil {
		log.Fatal("Error converting port number")
	}

	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	database, err := config.NewDatabase(redisPort, redisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	e := echo.New()
	storage.NewDB()

	Routes(e, database)
	e.Start(fmt.Sprintf(":%d", port))

}

// logger := logrus.New()
// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
// 	Format: "time=${time_rfc3339} method=${method}, uri=${uri}, status=${status}\n",
// 	Output: logger.Out,
// }))

// e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))

// file, err := os.OpenFile("../../server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// if err != nil {
// 	logger.Fatal(err)
// }
// defer file.Close()

// logger.SetOutput(file)
