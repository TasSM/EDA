package clients

import (
	"log"
	"strconv"

	"github.com/nats-io/nats.go"
)

func CreateNATSWorkerPool(workers int16) chan *nats.Msg {
	subch := make(chan *nats.Msg, 128)
	log.Printf("Starting %d workers", workers)
	for i := int16(0); i < workers; i++ {
		go func(id string) {
			log.Printf("Worker %s starting up", id)
			for {
				m, ok := <-subch
				if !ok {
					log.Printf("Worker %s channel was closed", id)
					return
				}
				log.Printf("Worker %s received message: %s", id, string(m.Data))
			}
		}(strconv.Itoa(int(i)))
	}
	return subch
}
