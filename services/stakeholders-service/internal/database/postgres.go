package database

import (
	"log"
	"stakeholders-service/internal/models" // Obavezno zamenite "stakeholders-service" sa imenom vašeg modula iz go.mod

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB uspostavlja konekciju sa PostgreSQL bazom i pokreće automatsku migraciju.
func InitDB() *gorm.DB {
	// Connection string za povezivanje na bazu
	dsn := "host=postgres user=postgres password=password dbname=stakeholders port=5432 sslmode=disable"
	
	// Otvaranje konekcije
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Ako konekcija ne uspe, program se gasi uz poruku o grešci
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate automatski kreira ili ažurira tabelu "stakeholders_users"
	// na osnovu `User` modela iz `models` paketa.
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connection successful and schema migrated.")
	
	// Vraća objekat konekcije koji će se koristiti u celoj aplikaciji
	return db
}