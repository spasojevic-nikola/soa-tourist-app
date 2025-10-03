package database

import (
	"log"
	"tour-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=tour-db user=postgres password=password dbname=tours port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("!!! FAILED TO CONNECT TO DATABASE:", err)
	}

	log.Println("Database connection successful. Running migration...")

	// --- KLJUČNA IZMENA JE OVDE ---
	// Eksplicitno proveravamo da li je migracija uspela.
	// Ako ne uspe, aplikacija će se srušiti i ispisati tačnu grešku.
	err = db.AutoMigrate(&models.Tour{}, &models.KeyPoint{}, &models.TourDuration{})
	if err != nil {
		log.Fatal("!!! FAILED TO MIGRATE DATABASE:", err)
	}

	log.Println("Database schema migrated successfully.")
	return db
}

/*package database

import (
	"log"
	"tour-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=tour-db user=postgres password=password dbname=tours port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.AutoMigrate(&models.Tour{})
	log.Println("Database connection successful and schema migrated.")
	return db
}*/
