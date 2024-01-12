package models

import "gorm.io/datatypes"

type DBEntry struct {
	OrderUID string `gorm:"primarykey;unique;not null"`
	Value    datatypes.JSON
}
