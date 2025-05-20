package initializers

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	dbHandle, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		log.Fatal("Connection with database couldn't be established.")
	}
	dbHandle.Exec("PRAGMA foreign_keys = ON")
	DB = dbHandle
}
