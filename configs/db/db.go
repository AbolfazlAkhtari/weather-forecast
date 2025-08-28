package db

import "gorm.io/gorm"

type DB interface {
	Open() gorm.Dialector
}
