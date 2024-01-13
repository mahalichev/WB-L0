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

	"github.com/mahalichev/WB-L0/config"
	"github.com/mahalichev/WB-L0/handlers"
	"github.com/mahalichev/WB-L0/models"
	"github.com/nats-io/stan.go"
)

func newMux(app *config.App) *http.ServeMux {
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

	sub, err := sc.Subscribe(os.Getenv("STAN_SUBJECT"), func(msg *stan.Msg) {
		// TODO (add data to sql and to cache)
		log.Printf("Received a message: %s\n", string(msg.Data))
		order := models.Order{}
		_ = json.Unmarshal(msg.Data, &order)
		app.Cache[order.OrderUID] = order
	})
	if err != nil {
		app.ErrLog.Print(err)
		return
	}
	defer sub.Unsubscribe()

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("SERVICE_ADDRESS"), os.Getenv("SERVICE_PORT")),
		ErrorLog: app.ErrLog,
		Handler:  newMux(app),
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			app.ErrLog.Printf("HTTP server error: %s", err.Error())
			return
		}
	}()
	app.InfoLog.Print("Service started successfully")
	UntilInterrupt()
	app.InfoLog.Printf("Exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		app.ErrLog.Printf("An error occurred while shutting down the server: %s", err.Error())
	}
	app.InfoLog.Print("Service shutted down")
}

func UntilInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
