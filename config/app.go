package config

import (
	"log"
	"os"

	"github.com/mahalichev/WB-L0/dao"
	"github.com/mahalichev/WB-L0/inits"
	"github.com/mahalichev/WB-L0/models"
)

type App struct {
	InfoLog *log.Logger
	ErrLog  *log.Logger
	Dao     *dao.DAO
	Cache   map[string]models.Order
}

func (app *App) AddToCache(order models.Order) {
	app.Cache[order.OrderUID] = order
}

func New() *App {
	app := &App{
		InfoLog: log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime),
		ErrLog:  log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile),
		Cache:   make(map[string]models.Order),
	}
	db, err := inits.GetDatabase()
	if err != nil {
		app.ErrLog.Fatalf("Can't connect to db: %s", err.Error())
	}
	app.Dao = dao.New(db)

	orders, err := app.Dao.SelectAll()
	if err != nil {
		app.ErrLog.Fatalf("Can't recover cache: %s", err.Error())
	}

	for _, order := range orders {
		app.AddToCache(order)
	}
	return app
}
