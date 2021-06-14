package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TasSM/EDA/internal/api"
	"github.com/TasSM/EDA/internal/client"
	"github.com/TasSM/EDA/internal/model"
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
	NATS_QUEUEGROUP  = "test-qg01"
	API_PORT         = "8088"
)

func main() {

	log.Printf("Connecting to NATS server on %s", nats.DefaultURL)
	rc := model.CreateNATSConn(nats.DefaultURL)
	wc := model.CreateNATSConn(nats.DefaultURL)
	defer rc.Close()
	defer wc.Close()
	log.Printf("Starting Web API on %s", API_PORT)
	go api.ServeTestAPI(API_PORT, wc)

	// register clients / subscribers

	// pubsub client
	psc := client.CreateNATSWorkerPool(2)
	pssub, _ := rc.ChanSubscribe(PUBSUB_EVENT, psc)
	defer close(psc)
	defer pssub.Unsubscribe()

	// request/response client
	rrsub, _ := rc.QueueSubscribe(QUEUEGROUP_EVENT, NATS_QUEUEGROUP, func(msg *nats.Msg) {
		reply := md5.Sum(msg.Data)
		msg.Respond([]byte(hex.EncodeToString(reply[:])))
	})
	defer rrsub.Drain()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
