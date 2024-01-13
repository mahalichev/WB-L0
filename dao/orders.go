package dao

import (
	"encoding/json"
	"fmt"

	"github.com/mahalichev/WB-L0/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DAO struct {
	db *gorm.DB
}

func New(db *gorm.DB) *DAO {
	return &DAO{db}
}

func (dao *DAO) InsertOrder(orderUID string, order datatypes.JSON) error {
	entry := models.DBOrder{OrderUID: orderUID, OrderData: order}

	result := dao.db.Create(&entry)
	// if result.Error != nil {
	// 23502 order_id is empty
	// check if items is empty
	// 	if err, ok := result.Error.(*pgconn.PgError); ok && err.Code == "23505" {
	// 		fmt.Printf("OrderUID %s already in use\n", order)
	// 	} else {
	// 		fmt.Println(result.Error)
	// 		return
	// 	}
	// }
	return result.Error
}

func (dao *DAO) SelectAll() ([]models.Order, error) {
	var dbOrders []models.DBOrder
	result := dao.db.Find(&dbOrders)
	if result.Error != nil {
		return nil, result.Error
	}
	orders := make([]models.Order, 0, len(dbOrders))
	for _, dbOrder := range dbOrders {
		var order models.Order
		jsonData, err := dbOrder.OrderData.MarshalJSON()
		if err != nil {
			return orders, fmt.Errorf("OrderData of %s is corrupted: %s. Skipping", dbOrder.OrderUID, err.Error())
		}
		if err := json.Unmarshal(jsonData, &order); err != nil {
			return orders, fmt.Errorf("OrderData of %s is not valid: %s. Skipping", dbOrder.OrderUID, err.Error())
		}
		orders = append(orders, order)
	}
	return orders, nil
}
