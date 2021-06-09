package api

import (
	"fmt"
	"log"
	"net/http"
)

func ServeTestAPI(port string) {
	svr := http.NewServeMux()
	svr.HandleFunc("/natsEvent", NatsEventHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), svr))
}
