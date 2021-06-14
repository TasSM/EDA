package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TasSM/EDA/internal/api"
	"github.com/TasSM/EDA/internal/clients"
	"github.com/TasSM/EDA/internal/models"
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

	log.Printf("Connecting to NATS server on %s", nats.DefaultURL)
	rc := models.CreateNATSConn(nats.DefaultURL)
	wc := models.CreateNATSConn(nats.DefaultURL)
	defer rc.Close()
	defer wc.Close()
	log.Printf("Starting Web API on %s", API_PORT)
	go api.ServeTestAPI(API_PORT, wc)

	// register clients / subscribers
	psc := clients.CreateNATSWorkerPool(2)
	pssub, _ := rc.ChanSubscribe(PUBSUB_EVENT, psc)
	defer close(psc)
	defer pssub.Unsubscribe()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
