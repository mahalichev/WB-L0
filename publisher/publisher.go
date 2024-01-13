package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/mahalichev/WB-L0/app"
	"github.com/mahalichev/WB-L0/inits"
	"github.com/nats-io/stan.go"
)

func init() {
	inits.LoadEnvironment()
}

func main() {
	sc, err := stan.Connect(os.Getenv("STAN_CLUSTERID"), os.Getenv("STAN_PUBLISHERID"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	go sendMessages(sc, os.Getenv("STAN_SUBJECT"), ticker)
	app.UntilInterrupt()
}

func sendMessages(sc stan.Conn, subject string, ticker *time.Ticker) {
	defer ticker.Stop()
	for {
		order := app.GenerateOrder()
		message, err := json.Marshal(order)
		if err != nil {
			log.Printf("something went wrong while creating order: %s", err.Error())
			return
		}
		<-ticker.C
		err = sc.Publish(subject, message)
		if err != nil {
			log.Printf("error publishing message: %s", err.Error())
		} else {
			log.Printf("sent order with UID: %s", order.OrderUID)
		}
	}
}
