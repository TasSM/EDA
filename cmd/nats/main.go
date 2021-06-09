package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TasSM/EDA/internal/api"
	"github.com/nats-io/nats.go"
)

/*
* Run the NATS server:
* docker run -d -p 4222:4222 -p 8222:8222 nats
 */

const (
	WORKER_COUNT     = 2
	PUBSUB_EVENT     = "pubsub"
	QUEUEGROUP_EVENT = "queuegroup"
	API_PORT         = "8088"
)

func main() {
	log.Printf("Starting Web API on %s", API_PORT)
	go api.ServeTestAPI(API_PORT)
	log.Printf("Connecting to NATS server on %s", nats.DefaultURL)
	nc, _ := nats.Connect(nats.DefaultURL)

	// pubsub subscriber channel
	subch := make(chan *nats.Msg, 128)
	sub, _ := nc.ChanSubscribe(PUBSUB_EVENT, subch)
	defer sub.Unsubscribe()

	for i := 0; i < WORKER_COUNT; i++ {
		go newWorker(fmt.Sprint(i), subch)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func newWorker(id string, ch chan *nats.Msg) {
	log.Printf("Worker %s starting up", id)
	for {
		m := <-ch
		log.Printf("Worker %s received message: %s", id, string(m.Data))
	}
}
