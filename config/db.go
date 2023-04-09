package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	DBType = "postgres"
)

func GetDBType() string {
	return DBType
}

func GetPostgresConnString() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBUser := os.Getenv("DBUser")
	DBPassword := os.Getenv("DBPassword")
	DBName := os.Getenv("DBName")
	DBHost := os.Getenv("DBHost")
	DBPort := os.Getenv("DBPort")

	dataBase := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		DBHost,
		DBPort,
		DBUser,
		DBName,
		DBPassword,
	)
	return dataBase
}
