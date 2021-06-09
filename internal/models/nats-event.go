package models

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func GenerateNATSEvents(count uint16, delay uint16) {
	var i uint16
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()
	for i = 0; i < count; i++ {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		log.Println("Publishing test NATS event to subject: pubsub")
		nc.Publish("pubsub", []byte(fmt.Sprintf("this is a test event at time %d", time.Now().UTC().Unix())))
	}
}
