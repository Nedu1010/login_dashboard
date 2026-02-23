package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Client *gorm.DB
}

func NewDB(database_url string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(database_url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{Client: db}, nil
}

func (db *DB) Close() {
	sqlDB, err := db.Client.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}

func (db *DB) HealthCheck() error {
	sqlDB, err := db.Client.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
