package shared

import (
	"log"

	"github.com/glebarez/sqlite"
	"github.com/openlyfree/gopher-ctf/internal/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := "ctf.db?_pragma=journal_mode(WAL)&_pragma=synchronous(normal)&_pragma=busy_timeout(5000)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Challenge{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.ChallengeCompletion{}); err != nil {
		log.Fatal(err)
	}

	return db
}
