package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

type Store struct {
	*gorm.DB
}

type SpeedMeasurement struct {
	ID        int
	Download  float64   `json:"download"`
	Upload    float64   `json:"upload"`
	Ping      float64   `json:"ping"`
	Timestamp time.Time `json:"timestamp"`
}

func CreateDB(filename string) *Store {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	migrationErr := db.AutoMigrate(&SpeedMeasurement{})
	if migrationErr != nil {
		log.Fatalf("DB Migration failed: %v", migrationErr)
	}

	return &Store{
		db,
	}
}

func (store *Store) SaveMeasurement(measurement SpeedMeasurement) {
	store.Create(&measurement)
}
