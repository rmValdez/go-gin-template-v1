package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database instance (similar to how PrismaClient works)
var DB *gorm.DB

func Connect(dsn string) {
	if dsn == "" {
		log.Println("DATABASE_URL is not set, skipping database connection...")
		return
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection successfully established")
}
