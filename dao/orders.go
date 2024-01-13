package dao

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
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

func (dao *DAO) InsertOrder(order models.Order) error {

	if len(order.OrderUID) == 0 {
		return errors.New("field order_id is empty")
	}

	if len(order.Items) == 0 {
		return errors.New("field items is empty")
	}

	raw, err := json.Marshal(order)
	if err != nil {
		return err
	}

	dbOrder := models.DBOrder{OrderUID: order.OrderUID, OrderData: datatypes.JSON(raw)}
	result := dao.db.Create(&dbOrder)
	if result.Error != nil {
		if err, ok := result.Error.(*pgconn.PgError); ok && err.Code == "23505" {
			return fmt.Errorf("order_uid %s already in use", order.OrderUID)
		} else {
			return result.Error
		}
	}
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
			return orders, fmt.Errorf("field \"OrderData\" of %s is corrupted: %s. Skipping", dbOrder.OrderUID, err.Error())
		}
		if err := json.Unmarshal(jsonData, &order); err != nil {
			return orders, fmt.Errorf("field \"OrderData\" of %s is not valid: %s. Skipping", dbOrder.OrderUID, err.Error())
		}
		orders = append(orders, order)
	}
	return orders, nil
}
