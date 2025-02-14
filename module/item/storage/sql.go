package storage

import (
	"gorm.io/gorm"
)

type sqlStorage struct {
	db *gorm.DB
}

func newSqlStorage(db *gorm.DB) *sqlStorage {
	return &sqlStorage{db: db}
}
