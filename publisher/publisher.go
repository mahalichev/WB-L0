package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/mahalichev/WB-L0/models"
	"github.com/nats-io/stan.go"
)

const subject = "service"
const clusterID = "test-cluster"
const clientID = "publisher"

func pseudoRandomString() string {
	hasher := md5.New()
	hasher.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	timer := time.NewTimer(10)
	<-timer.C
	return hashString
}

func generateItem() models.Item {
	return models.Item{
		ChrtID:      rand.Intn(1<<31 - 1),
		TrackNumber: pseudoRandomString(),
		Price:       rand.Intn(1<<31 - 1),
		RID:         pseudoRandomString(),
		Name:        pseudoRandomString(),
		Sale:        rand.Intn(1<<31 - 1),
		Size:        pseudoRandomString(),
		TotalPrice:  rand.Intn(1<<31 - 1),
		NMID:        rand.Intn(1<<31 - 1),
		Brand:       pseudoRandomString(),
		Status:      rand.Intn(1<<31 - 1),
	}
}

func generatePayment() models.Payment {
	return models.Payment{
		Transaction:  pseudoRandomString(),
		RequestID:    pseudoRandomString(),
		Currency:     pseudoRandomString(),
		Provider:     pseudoRandomString(),
		Amount:       rand.Intn(1<<31 - 1),
		PaymentDT:    rand.Intn(1<<31 - 1),
		Bank:         pseudoRandomString(),
		DeliveryCost: rand.Intn(1<<31 - 1),
		GoodsTotal:   rand.Intn(1<<31 - 1),
		CustomFEE:    rand.Intn(1<<31 - 1),
	}
}

func generateDelivery() models.Delivery {
	return models.Delivery{
		Name:    pseudoRandomString(),
		Phone:   pseudoRandomString(),
		Zip:     pseudoRandomString(),
		City:    pseudoRandomString(),
		Address: pseudoRandomString(),
		Region:  pseudoRandomString(),
		Email:   pseudoRandomString(),
	}
}

func generateOrder() models.Order {
	order := models.Order{
		OrderUID:          pseudoRandomString(),
		TrackNumber:       pseudoRandomString(),
		Entry:             pseudoRandomString(),
		Locale:            pseudoRandomString(),
		InternalSignature: pseudoRandomString(),
		CustomerID:        pseudoRandomString(),
		DeliveryService:   pseudoRandomString(),
		Shardkey:          pseudoRandomString(),
		SMID:              rand.Intn(1<<31 - 1),
		DateCreated:       time.Now().Format("2006-01-02T15:04:05Z"),
		OOFShard:          pseudoRandomString(),
	}
	itemsLen := rand.Intn(3) + 1
	order.Items = make([]models.Item, itemsLen)
	for i := 0; i < itemsLen; i++ {
		order.Items[i] = generateItem()
	}
	order.Payment = generatePayment()
	order.Delivery = generateDelivery()
	return order
}

func main() {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	sendMessages(sc)
}

func sendMessages(sc stan.Conn) {
	ticker := time.NewTicker(3 * time.Second)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		order := generateOrder()
		message, err := json.Marshal(order)
		if err != nil {
			log.New(os.Stdin, "ERROR", 1).Println("something went wrong while creating order")
			return
		}
		select {
		case <-ticker.C:
			err := sc.Publish(subject, message)
			if err != nil {
				log.Println("Error publishing message:", err)
			} else {
				log.Printf("Sent order with UID: %s\n", order.OrderUID)
			}
		case <-c:
			fmt.Println("Exit")
			return
		}
	}
}
