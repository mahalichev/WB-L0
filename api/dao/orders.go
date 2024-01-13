package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mahalichev/WB-L0/api/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DAO struct {
	db       *gorm.DB
	validate *validator.Validate
}

func New(db *gorm.DB) *DAO {
	return &DAO{db, validator.New()}
}

func (dao *DAO) InsertOrder(order models.Order) error {
	if err := dao.validate.Struct(order); err != nil {
		return fmt.Errorf("object must contain all values: %s", err.Error())
	}

	if len(order.Items) == 0 {
		return errors.New("field items is empty")
	}

	if _, err := time.Parse("2006-01-02T15:04:05Z", order.DateCreated); err != nil {
		return errors.New("field date_created has wrong time format")
	}

	raw, err := json.Marshal(order)
	if err != nil {
		return err
	}

	dbOrder := models.DBOrder{OrderUID: order.OrderUID, OrderData: datatypes.JSON(raw)}
	if result := dao.db.Create(&dbOrder); result.Error != nil {
		if err, ok := result.Error.(*pgconn.PgError); ok && err.Code == "23505" {
			return fmt.Errorf("order_uid %s already in use", order.OrderUID)
		} else {
			return result.Error
		}
	}
	return nil
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
			return orders, fmt.Errorf("field \"OrderData\" of %s is corrupted: %s", dbOrder.OrderUID, err.Error())
		}
		if err := json.Unmarshal(jsonData, &order); err != nil {
			return orders, fmt.Errorf("can't unmarshal \"OrderData\" field of %s: %s", dbOrder.OrderUID, err.Error())
		}
		if err := dao.validate.Struct(order); err != nil {
			return orders, fmt.Errorf("object with order_uid %s must contain all values: %s", dbOrder.OrderUID, err.Error())
		}
		orders = append(orders, order)
	}
	return orders, nil
}
