package config

import (
	"expensetracker/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBUrl     string
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		DBUrl:     os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}

func ConnectDatabase(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&models.Expense{})
	if err != nil {
		log.Fatalf("Failed to run AutoMigrate: %v", err)
	}
	return db
}
