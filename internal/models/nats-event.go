package models

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func CreateNATSConn(addr string) *nats.Conn {
	nc, err := nats.Connect(addr)
	if err != nil {
		log.Fatalf("Failed to connect to NATS server at %s", addr)
	}
	return nc
}

func GenerateNATSEventsPubSub(nc *nats.Conn, count uint16, delay uint16) {
	var i uint16
	log.Printf("Received NATS event generation command: count %d, delay(ms): %d", count, delay)
	for i = 0; i < count; i++ {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		nc.Publish("pubsub", []byte(fmt.Sprintf("this is a test event at time %d", time.Now().UTC().Unix())))
	}
	log.Printf("NATS event creation completed.")
}
