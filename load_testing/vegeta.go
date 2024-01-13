package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mahalichev/WB-L0/api/inits"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func init() {
	inits.LoadEnvironment()
}

func main() {
	attacker := vegeta.NewAttacker()
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    fmt.Sprintf("http://%s:%s/orders?id=%s", os.Getenv("SERVICE_ADDRESS"), os.Getenv("SERVICE_PORT"), os.Getenv("VEGETA_ORDERUID")),
	})

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, vegeta.Rate{Freq: 1000, Per: time.Second}, 3*time.Second, "GET order test") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Println("Latency stats")
	fmt.Printf(" - Mean: %s\n", metrics.Latencies.Mean)
	fmt.Printf(" - 99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf(" - Max: %s\n", metrics.Latencies.Max)
	fmt.Println("\nOverall")
	fmt.Printf(" - Total success: %.5f\n", metrics.Success)
	fmt.Printf(" - Errors: %v\n", metrics.Errors)
}
