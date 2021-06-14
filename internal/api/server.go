package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

func ServeTestAPI(port string, nc *nats.Conn) {
	svr := http.NewServeMux()
	svr.HandleFunc("/natsEvent", NatsPubSubHandler(nc))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), svr))
}
