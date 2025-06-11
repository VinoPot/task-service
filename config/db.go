package config

import (
	"fmt"
	"os"
	"strconv" // ← отсутствовал

	"github.com/joho/godotenv" // ← тоже отсутствовал
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Ошибка загрузки .env")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	port, err := strconv.Atoi(dbPort)
	if err != nil {
		panic("DB_PORT должен быть числом")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, port, dbSSLMode, dbTimezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	return db
}
