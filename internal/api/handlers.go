package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TasSM/EDA/internal/model"
	"github.com/gorilla/schema"
	"github.com/nats-io/nats.go"
)

var decoder = schema.NewDecoder()

type NATSEventRequest struct {
	Count   uint16
	DelayMs uint16
}

type NATSRequest struct {
	Subject  string
	Message  string
	Checksum string
}

func NatsPubSubHandler(nc *nats.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var natsRequest NATSEventRequest
		err := decoder.Decode(&natsRequest, r.URL.Query())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request: %s", err.Error())))
			return
		}
		go model.GenerateNATSEventsPubSub(nc, natsRequest.Count, natsRequest.DelayMs)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(natsRequest)
	}
}

func NatsRequestHandler(nc *nats.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var natsRequest NATSRequest
		err := decoder.Decode(&natsRequest, r.URL.Query())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request: %s", err.Error())))
			return
		}
		res, err := model.MakeNatsRequest(nc, natsRequest.Subject, []byte(natsRequest.Message))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("NATS Request Failed: %s", err.Error())))
			return
		}
		natsRequest.Checksum = res
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(natsRequest)
	}
}
