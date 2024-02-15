package db

import (
	"fmt"
	"github.com/amirh-khali/go-playground/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	config := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		"localhost",
		"5432",
		"postgres",
		"go-playground",
		"1234",
	)

	database, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = database.AutoMigrate(&models.Recipe{})
	if err != nil {
		return
	}

	DB = database
}
