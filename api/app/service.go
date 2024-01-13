package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mahalichev/WB-L0/api/config"
	"github.com/mahalichev/WB-L0/api/handlers"
	"github.com/mahalichev/WB-L0/api/models"
	"github.com/nats-io/stan.go"
)

func newMux(app *config.AppConfig) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetRoot(w, r)
		}
	})
	mux.HandleFunc("/static/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetOrderHTML(app)(w, r)
		}
	})
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetOrderJSON(app)(w, r)
		}
	})
	return mux
}

func RunService() {
	app := config.New()
	sc, err := stan.Connect(os.Getenv("STAN_CLUSTERID"), os.Getenv("STAN_SERVICEID"))
	if err != nil {
		log.Fatal("can't connect to NATS-Streaming subsystem")
	}
	defer sc.Close()

	_, err = sc.Subscribe(os.Getenv("STAN_SUBJECT"), func(msg *stan.Msg) {
		var order models.Order
		app.InfoLog.Print("recieved message")
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			app.ErrLog.Printf("can't unmarshal json data: %s", err.Error())
			return
		}
		if err := app.Dao.InsertOrder(order); err != nil {
			app.ErrLog.Printf("can't insert order into db: %s", err.Error())
			return
		}
		app.AddToCache(order)
		app.InfoLog.Printf("order with order_uid %s cached successfully", order.OrderUID)
	}, stan.DurableName(os.Getenv("STAN_DURABLENAME")), stan.DeliverAllAvailable())

	if err != nil {
		app.ErrLog.Print(err)
		return
	}

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("SERVICE_ADDRESS"), os.Getenv("SERVICE_PORT")),
		ErrorLog: app.ErrLog,
		Handler:  newMux(app),
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			app.ErrLog.Printf("HTTP server error: %s", err.Error())
		}
	}()
	app.InfoLog.Print("service started successfully")
	UntilInterrupt()
	app.InfoLog.Printf("exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		app.ErrLog.Printf("an error occurred while shutting down the server: %s", err.Error())
	}
	app.InfoLog.Print("service shutted down")
}

func UntilInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
