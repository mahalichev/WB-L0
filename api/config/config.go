package config

import (
	"log"
	"os"

	"github.com/mahalichev/WB-L0/api/dao"
	"github.com/mahalichev/WB-L0/api/inits"
	"github.com/mahalichev/WB-L0/api/models"
)

type AppConfig struct {
	InfoLog *log.Logger
	ErrLog  *log.Logger
	Dao     *dao.DAO
	Cache   map[string]models.Order
}

func (app *AppConfig) AddToCache(order models.Order) {
	app.Cache[order.OrderUID] = order
}

func New() *AppConfig {
	app := &AppConfig{
		InfoLog: log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime),
		ErrLog:  log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile),
		Cache:   make(map[string]models.Order),
	}
	db, err := inits.GetDatabase()
	if err != nil {
		app.ErrLog.Fatalf("can't connect to db: %s", err.Error())
	}
	app.Dao = dao.New(db)

	orders, err := app.Dao.SelectAll()
	if err != nil {
		app.ErrLog.Fatalf("can't recover cache: %s", err.Error())
	}

	for _, order := range orders {
		app.AddToCache(order)
	}
	app.InfoLog.Printf("recovered %d orders", len(app.Cache))
	return app
}
