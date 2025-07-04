package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"Go-Exercise/pkg/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	db.AutoMigrate(&model.User{}, &model.RefreshToken{})
	DB = db
}

func InitDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	log.Println("Connecting to database...")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Updating existing records...")
	now := time.Now()

	if db.Migrator().HasTable(&model.User{}) {
		db.Model(&model.User{}).Where("created_at IS NULL").Updates(map[string]interface{}{
			"created_at": now,
			"updated_at": now,
		})
	}

	if db.Migrator().HasTable(&model.RefreshToken{}) {
		db.Model(&model.RefreshToken{}).Where("created_at IS NULL").Updates(map[string]interface{}{
			"created_at": now,
			"updated_at": now,
		})
	}

	log.Println("Running migrations...")
	err = db.AutoMigrate(&model.User{}, &model.RefreshToken{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database initialization completed successfully")
	return db
}
