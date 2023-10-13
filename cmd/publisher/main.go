package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/racoon-proger/wb-l0/internal/domain"
)

type configuration struct {
	NatsAddr string `envconfig:"NATS_ADDR" default:"nats://127.0.0.1:4222"`
}

func main() {
	var cfg configuration
	nc, err := nats.Connect(cfg.NatsAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile("./data/order.json")
	if err != nil {
		log.Fatal(err)
	}
	var orders []domain.Order

	err = json.Unmarshal(data, &orders)
	if err != nil {
		log.Fatal(err)
	}

	for i := range orders {
		var data []byte
		data, err = json.Marshal(&orders[i])
		if err != nil {
			log.Fatal(err)
		}
		_, err = js.PublishAsync("foo", data)
		if err != nil {
			log.Fatal(err)
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
