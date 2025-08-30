package db

import "gorm.io/gorm"

type DB interface {
	Open(migrate bool) (*gorm.DB, error)
}
