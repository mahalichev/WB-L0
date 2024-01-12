package dao

import (
	"github.com/mahalichev/WB-L0/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DAO struct {
	db *gorm.DB
}

func New(db *gorm.DB) DAO {
	return DAO{db}
}

func (dao *DAO) InsertOrder(orderUID string, order datatypes.JSON) error {
	entry := models.DBEntry{OrderUID: orderUID, Value: order}

	result := dao.db.Create(&entry)
	// if result.Error != nil {
	// 	if err, ok := result.Error.(*pgconn.PgError); ok && err.Code == "23505" {
	// 		fmt.Printf("OrderUID %s already in use\n", order)
	// 	} else {
	// 		fmt.Println(result.Error)
	// 		return
	// 	}
	// }
	return result.Error
}

func (dao *DAO) SelectAll() ([]models.DBEntry, error) {
	var orders []models.DBEntry
	result := dao.db.Find(&orders)
	return orders, result.Error
}
