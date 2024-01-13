package models

import "gorm.io/datatypes"

type DBOrder struct {
	OrderUID  string `gorm:"primarykey;unique;not null;type:varchar(100);default:null"`
	OrderData datatypes.JSON
}
