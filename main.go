package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/mahalichev/WB-L0/dao"
	"github.com/mahalichev/WB-L0/inits"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const clusterID = "test-cluster"
const clientID = "listener"
const subject = "service"

func init() {
	inits.LoadEnvironment()
}

func main() {
	db, err := inits.ConnectToDatabase()
	if err != nil {
		log.Fatalf("can't connect to database: %s\n", err)
	}
	_ = dao.New(db)

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(nats.DefaultURL))
	if err != nil {
		log.Fatal("Can't connect to NATS-Streaming subsystem")
	}
	defer sc.Close()
	sub, err := sc.Subscribe(subject, func(msg *stan.Msg) {
		log.Printf("Received a message: %s\n", string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Exiting...")
	os.Exit(0)
}
