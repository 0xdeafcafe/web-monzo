package helpers

import (
	"github.com/0xdeafcafe/web-monzo/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// NewDatabaseConnection creates a new database connection
func NewDatabaseConnection(connectionString string) *gorm.DB {
	db, err := gorm.Open("mysql", connectionString)
	db.LogMode(true)

	if err != nil {
		panic(err)
	}

	// Run those dank migrations
	db.AutoMigrate(&models.Cookie{})

	return db
}
